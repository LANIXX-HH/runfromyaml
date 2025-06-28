package functions

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
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
		for _, command := range commands {
			mymap["cmd"] = append(mymap["cmd"].([]map[string]interface{}), map[string]interface{}{
				"type": "shell",
				"values": []string{
					command,
				},
			})
		}
	}

	if len(envs) > 0 {
		mymap["env"] = []map[string]interface{}{}

		for key, value := range envs {
			mymap["env"] = append(mymap["env"].([]map[string]interface{}), map[string]interface{}{
				"key":   key,
				"value": value,
			})
		}
	}

	return mymap
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
