package functions

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"text/template"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

var (
	RestOut http.ResponseWriter
	ReqOut  *http.Request
)

// WriteFile write a file
func WriteFile(file string, path string, perm os.FileMode) {
	bytefile := []byte(file)
	err := os.WriteFile(os.ExpandEnv(path), bytefile, perm)
	if err != nil {
		panic(err)
	}
}

// ReadFile read file
func ReadFile(file string) {
	content, err := os.ReadFile(os.ExpandEnv(file))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("File contents: %s", content)
}

// Remove element from slice
func Remove(slice []string, i int) []string {
	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
}

func Config() {

}

func PrintSwitch(ctype color.Attribute, _level string, _output string, cstring ...interface{}) {
	switch _output {
	case "rest":
		PrintRest(ctype, _level, cstring...)
	case "file":
		PrintFile(_level, cstring...)
	case "stdout":
		PrintColor(ctype, _level, cstring...)
	}
}

func PrintFile(_level string, cstring ...interface{}) {
	log := logrus.New()
	// file
	log.Formatter = new(logrus.JSONFormatter)                      //default
	log.Formatter.(*logrus.JSONFormatter).PrettyPrint = true       // pretty print
	log.Formatter.(*logrus.JSONFormatter).DisableTimestamp = false // remove timestamp from test output
	dir := os.TempDir()
	file, err := os.OpenFile(dir+"logrus-"+time.Now().Format("20060102")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	switch _level {
	case "info":
		log.Info(cstring...)
	case "warn":
		log.Warn(cstring...)
	case "error":
		log.Error(cstring...)
	case "debug":
		log.Debug(cstring...)
	case "trace":
		log.Trace(cstring...)
	case "fatal":
		log.Fatal(cstring...)
	case "panic":
		log.Panic(cstring...)
	}
}

func PrintRestHeader() {
	log := logrus.New()
	//file
	log.Formatter = new(logrus.JSONFormatter)                      //default
	log.Formatter.(*logrus.JSONFormatter).PrettyPrint = false      // pretty print
	log.Formatter.(*logrus.JSONFormatter).DisableTimestamp = true  // remove timestamp from test output
	log.Formatter.(*logrus.JSONFormatter).DisableHTMLEscape = true // remove timestamp from test output

	log.Out = RestOut
	log.WithFields(logrus.Fields{
		"uri":        ReqOut.RequestURI,
		"method":     ReqOut.Method,
		"host":       ReqOut.Host,
		"remoteaddr": ReqOut.RemoteAddr,
		"header":     ReqOut.Header,
	}).Info("Header Information")
}

func PrintRest(ctype color.Attribute, _level string, cstring ...interface{}) {
	mystring := color.New(ctype)
	_, _ = mystring.Fprintln(RestOut, cstring...)
}

func PrintColor(ctype color.Attribute, _level string, cstring ...interface{}) {
	mystring := color.New(ctype)
	_, _ = mystring.Println(cstring...)
}

func GoTemplate(mymap map[string]string, mytemplate string) string {
	var writer bytes.Buffer
	t, err := template.New("todos").Parse(mytemplate)

	if err != nil {
		panic(err)
	}
	err = t.Execute(&writer, mymap)
	if err != nil {
		panic(err)
	}
	return writer.String()
}

// The function first checks if the key exists in the yblock map and if it does, it trims the value associated with the key and converts it to a string.
// If the expandenv key exists in the yblock map and is set to true, the function then expands any environment variables in the string.
// Finally, the function splits the string into an array of strings and returns it. If the key does not exist, then the function returns nil.

// ExtractAndExpand extracts values from a YAML block and optionally expands environment variables.
// It now handles empty values gracefully, returning an empty slice instead of nil for empty values.
func ExtractAndExpand(yblock map[interface{}]interface{}, key string) []string {
	if !reflect.ValueOf(yblock[key]).IsValid() {
		// Return empty slice instead of nil to allow empty values blocks
		return []string{}
	}

	var result []string

	// Handle different types of values
	switch v := yblock[key].(type) {
	case []interface{}:
		// Handle YAML array
		for _, item := range v {
			if str, ok := item.(string); ok {
				result = append(result, str)
			} else {
				result = append(result, fmt.Sprint(item))
			}
		}
	case string:
		// Handle single string value
		if v != "" {
			result = []string{v}
		}
	case nil:
		// Handle explicit nil values
		return []string{}
	default:
		// Handle other types by converting to string
		str := fmt.Sprint(v)
		if str != "" && str != "[]" && str != "<nil>" {
			result = []string{str}
		}
	}

	// Apply environment variable expansion if requested
	if reflect.ValueOf(yblock["expandenv"]).IsValid() && yblock["expandenv"].(bool) {
		for i, val := range result {
			result[i] = os.ExpandEnv(val)
		}
	}

	return result
}

func PrintShellCommandsAsYaml(commands []string, envs map[string]string) map[string]interface{} {

	mymap := map[string]interface{}{
		"logging": []map[string]interface{}{
			{
				"level": "info",
			},
			{
				"output": "stdout",
			},
		},
		// "env": []map[string]interface{}{
		// 	{
		// 		"key":   "TEST",
		// 		"value": "foo",
		// 	},
		// 	{
		// 		"key":   "BLA",
		// 		"value": "TEST",
		// 	},
		// },
		// "cmd": []map[string]interface{}{
		// 	{
		// 		"type": "shell",
		// 		"values": []string{
		// 			"ls",
		// 		},
		// 	},
		// },
	}

	if len(commands) > 0 {

		mymap["cmd"] = []map[string]interface{}{}
		for i, command := range commands {
			mymap["cmd"] = append(mymap["cmd"].([]map[string]interface{}), map[string]interface{}{
				"type":      "shell",
				"name":      fmt.Sprintf("command-%d", i+1),
				"desc":      fmt.Sprintf("Interactive command: %s", command),
				"expandenv": true,
				"values": []string{
					command,
				},
			})
		}
	}

	// Filter environment variables - only include relevant/custom ones
	filteredEnvs := filterRelevantEnvVars(envs)
	if len(filteredEnvs) > 0 {
		mymap["env"] = []map[string]interface{}{}

		for key, value := range filteredEnvs {
			mymap["env"] = append(mymap["env"].([]map[string]interface{}), map[string]interface{}{
				"key":   key,
				"value": value,
			})
		}
	}

	return mymap
}

// filterRelevantEnvVars filters out system/session-specific environment variables
// and keeps only relevant/custom ones that should be documented
func filterRelevantEnvVars(envs map[string]string) map[string]string {
	// System/session variables to exclude (common across Unix/Linux/macOS/Windows)
	systemVars := map[string]bool{
		// System paths and directories
		"HOME":                true,
		"TMPDIR":              true,
		"TMP":                 true,
		"TEMP":                true,
		"PATH":                true,
		"LD_LIBRARY_PATH":     true,
		"DYLD_LIBRARY_PATH":   true,
		"PWD":                 true,
		"OLDPWD":              true,
		
		// User and session info
		"USER":                true,
		"USERNAME":            true,
		"LOGNAME":             true,
		"SHELL":               true,
		"SHLVL":               true,
		"TTY":                 true,
		"SSH_AUTH_SOCK":       true,
		"SSH_SOCKET_DIR":      true,
		
		// Terminal and display
		"TERM":                true,
		"TERM_PROGRAM":        true,
		"TERM_PROGRAM_VERSION": true,
		"COLORTERM":           true,
		"DISPLAY":             true,
		"WARP_HONOR_PS1":      true,
		"WARP_IS_LOCAL_SHELL_SESSION": true,
		"WARP_USE_SSH_WRAPPER": true,
		
		// System internals
		"XPC_SERVICE_NAME":    true,
		"XPC_FLAGS":           true,
		"COMMAND_MODE":        true,
		"LC_CTYPE":            true,
		"LC_ALL":              true,
		"LANG":                true,
		"__CF_USER_TEXT_ENCODING": true,
		"__CFBundleIdentifier": true,
		"_":                   true,
		
		// Package managers (Homebrew, Conda, etc.)
		"HOMEBREW_PREFIX":     true,
		"HOMEBREW_CELLAR":     true,
		"HOMEBREW_REPOSITORY": true,
		"CONDA_CHANGEPS1":     true,
		"CONDA_DEFAULT_ENV":   true,
		"CONDA_PREFIX":        true,
		"INFOPATH":            true,
		
		// Process/shell specific
		"SHELL_PID":           true,
		"PPID":                true,
		"PID":                 true,
		
		// Q CLI specific (our own tool)
		"Q_SET_PARENT_CHECK":  true,
	}
	
	// Prefixes to exclude (variables starting with these)
	excludePrefixes := []string{
		"BASH_",
		"ZSH_",
		"FISH_",
		"PS1",
		"PS2",
		"PS3",
		"PS4",
		"PROMPT_",
		"LESS",
		"PAGER",
		"EDITOR",
		"VISUAL",
		"MANPATH",
		"HISTFILE",
		"HISTSIZE",
		"HISTCONTROL",
		"XDG_",
		"DBUS_",
		"GNOME_",
		"KDE_",
		"QT_",
		"GTK_",
	}
	
	filtered := make(map[string]string)
	
	for key, value := range envs {
		// Skip if it's a known system variable
		if systemVars[key] {
			continue
		}
		
		// Skip if it starts with excluded prefixes
		shouldSkip := false
		for _, prefix := range excludePrefixes {
			if strings.HasPrefix(key, prefix) {
				shouldSkip = true
				break
			}
		}
		if shouldSkip {
			continue
		}
		
		// Include variables that are likely custom/relevant:
		// - AWS, DOCKER, KUBERNETES related
		// - Custom application variables
		// - Development environment variables
		if isRelevantEnvVar(key) {
			filtered[key] = value
		}
	}
	
	return filtered
}

// isRelevantEnvVar determines if an environment variable is relevant for documentation
func isRelevantEnvVar(key string) bool {
	relevantPrefixes := []string{
		"AWS_",
		"DOCKER_",
		"KUBE",
		"K8S_",
		"HELM_",
		"TERRAFORM_",
		"ANSIBLE_",
		"JENKINS_",
		"CI_",
		"BUILD_",
		"DEPLOY_",
		"ENV_",
		"APP_",
		"API_",
		"DB_",
		"DATABASE_",
		"REDIS_",
		"MONGO_",
		"MYSQL_",
		"POSTGRES_",
		"NODE_",
		"PYTHON_",
		"JAVA_",
		"GO_",
		"RUST_",
		"PHP_",
		"RUBY_",
		"GIT_",
		"GITHUB_",
		"GITLAB_",
		"BITBUCKET_",
	}
	
	for _, prefix := range relevantPrefixes {
		if strings.HasPrefix(key, prefix) {
			return true
		}
	}
	
	// Also include some specific relevant variables
	relevantVars := map[string]bool{
		"PORT":        true,
		"HOST":        true,
		"DEBUG":       true,
		"ENVIRONMENT": true,
		"STAGE":       true,
		"VERSION":     true,
		"REGION":      true,
		"ZONE":        true,
		"NAMESPACE":   true,
		"SERVICE":     true,
		"CONFIG":      true,
		"SECRET":      true,
		"TOKEN":       true,
		"KEY":         true,
		"URL":         true,
		"ENDPOINT":    true,
	}
	
	return relevantVars[key]
}

func EvaluateDescription(yamlBlock map[interface{}]interface{}, defaultDescription ...string) string {
	var desc string
	if len(defaultDescription) > 0 {
		desc = defaultDescription[0]
	}
	if reflect.ValueOf(yamlBlock["desc"]).IsValid() {
		desc = fmt.Sprintf("# %v", yamlBlock["desc"])
	}
	return desc
}
