package cli

import (
	"bufio"
	"fmt"
	"os"
	execute "os/exec"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"gopkg.in/yaml.v2"

	"github.com/lanixx/runfromyaml/pkg/exec"
	functions "github.com/lanixx/runfromyaml/pkg/functions"
	"github.com/lanixx/runfromyaml/pkg/openai"
)

var (
	EnvironmentVariables map[string]string
)

func execCmd(yamlBlock map[interface{}]interface{}, environmentVariablesShell []string, outputLevel string, outputType string) {
	var (
		cmds []string
		desc string
	)
	wg := new(sync.WaitGroup)

	cmds = functions.ExtractAndExpand(yamlBlock, "values")

	wg.Add(1)
	go exec.Command(cmds, desc, wg, environmentVariablesShell, outputLevel, outputType)
	wg.Wait()
}

func shellCmd(yamlBlock map[interface{}]interface{}, environmentVariablesShell []string, outputLevel string, outputType string) {
	var (
		cmds []string
		desc string
	)
	cmds = functions.ExtractAndExpand(yamlBlock, "values")

	temp_cmds := strings.Join(cmds, " ")
	cmds = strings.Split(temp_cmds, ";")
	wg := new(sync.WaitGroup)
	for ind, shcmds := range cmds {
		shcmd := strings.Split(shcmds, " ")
		wg.Add(1)
		go exec.CommandShell(shcmd, desc, wg, ind, environmentVariablesShell, outputLevel, outputType)
		wg.Wait()
	}
}

func dockerCmd(yamlBlock map[interface{}]interface{}, environmentVariablesShell []string, outputLevel string, outputType string) {
	var (
		cmds []string
		desc string
	)
	wg := new(sync.WaitGroup)
	cmds = functions.ExtractAndExpand(yamlBlock, "values")

	wg.Add(1)
	go exec.CommandDockerRun(yamlBlock["command"].(string), yamlBlock["container"].(string), cmds, desc, environmentVariablesShell, wg, outputLevel, outputType)
	wg.Wait()
}

func dockerComposeCmd(yamlBlock map[interface{}]interface{}, environmentVariablesShell []string, outputLevel string, outputType string) {
	var (
		cmds       []string
		dcoptions  []string
		cmdoptions []string
		desc       string
		service    string
	)
	wg := new(sync.WaitGroup)
	cmds = functions.ExtractAndExpand(yamlBlock, "values")
	dcoptions = functions.ExtractAndExpand(yamlBlock, "dcoptions")
	cmdoptions = functions.ExtractAndExpand(yamlBlock, "cmdoptions")

	service = "no-service"
	if reflect.ValueOf(yamlBlock["service"]).IsValid() {
		service = yamlBlock["service"].(string)
	}
	wg.Add(1)
	go exec.CommandDockerComposeExec(yamlBlock["command"].(string), service, cmdoptions, dcoptions, cmds, desc, environmentVariablesShell, wg, outputLevel, outputType)
	wg.Wait()
}

func sshCmd(yamlBlock map[interface{}]interface{}, environmentVariablesShell []string, outputLevel string, outputType string) {
	var (
		cmds    []string
		options []string
		desc    string
	)
	var user string
	var host string
	wg := new(sync.WaitGroup)
	cmds = functions.ExtractAndExpand(yamlBlock, "values")
	options = functions.ExtractAndExpand(yamlBlock, "options")

	user = yamlBlock["user"].(string)
	host = yamlBlock["host"].(string)
	if reflect.ValueOf(yamlBlock["expandenv"]).IsValid() {
		if yamlBlock["expandenv"].(bool) {
			user = os.ExpandEnv(user)
			host = os.ExpandEnv(host)
		}
	}

	wg.Add(1)
	go exec.CommandSSH(user, yamlBlock["port"].(int), host, options, cmds, desc, environmentVariablesShell, wg, outputLevel, outputType)
	wg.Wait()
}

func conf(yamlBlock map[interface{}]interface{}, outputLevel string, outputType string) {
	var (
		desc     string
		confdata string
		confdest string
		confperm os.FileMode
	)
	if reflect.ValueOf(yamlBlock["confdata"].(string)).String() != "" {
		confdata = yamlBlock["confdata"].(string)
		if reflect.ValueOf(yamlBlock["expandenv"]).IsValid() {
			if yamlBlock["expandenv"].(bool) {
				confdata = functions.GoTemplate(EnvironmentVariables, confdata)
			}
		}
	}

	desc = functions.EvaluateDescription(yamlBlock)
	confdata = desc + confdata

	if reflect.ValueOf(yamlBlock["confdest"].(string)).String() != "" {
		confdest = yamlBlock["confdest"].(string)
	}
	if reflect.ValueOf(yamlBlock["confperm"].(int)).Int() != 0 {
		confperm = os.FileMode(int(yamlBlock["confperm"].(int)))
	}
	if confdata != "" && confdest != "" && string(rune(confperm)) != "" {
		functions.WriteFile(confdata, confdest, confperm)
	}
	if reflect.ValueOf(yamlBlock["expandenv"]).IsValid() {
		if yamlBlock["expandenv"].(bool) {
			confdest = os.ExpandEnv(confdest)
		}
	}
	functions.PrintSwitch(color.FgGreen, outputLevel, outputType, "# create ", confdest)
}

func execFunctionsMap(map[interface{}]interface{}, []string, string, string) map[string]func(map[interface{}]interface{}, []string, string, string) {
	execFunctions := map[string]func(map[interface{}]interface{}, []string, string, string){
		"exec": func(yamlBlock map[interface{}]interface{}, environmentVariablesShell []string, outputLevel string, outputType string) {
			execCmd(yamlBlock, environmentVariablesShell, outputLevel, outputType)
		},
		"shell": func(yamlBlock map[interface{}]interface{}, environmentVariablesShell []string, outputLevel string, outputType string) {
			shellCmd(yamlBlock, environmentVariablesShell, outputLevel, outputType)
		},
		"docker": func(yamlBlock map[interface{}]interface{}, environmentVariablesShell []string, outputLevel string, outputType string) {
			dockerCmd(yamlBlock, environmentVariablesShell, outputLevel, outputType)
		},
		"docker-compose": func(yamlBlock map[interface{}]interface{}, environmentVariablesShell []string, outputLevel string, outputType string) {
			dockerComposeCmd(yamlBlock, environmentVariablesShell, outputLevel, outputType)
		},
		"ssh": func(yamlBlock map[interface{}]interface{}, environmentVariablesShell []string, outputLevel string, outputType string) {
			sshCmd(yamlBlock, environmentVariablesShell, outputLevel, outputType)
		},
		"conf": func(yamlBlock map[interface{}]interface{}, environmentVariablesShell []string, outputLevel string, outputType string) {
			conf(yamlBlock, outputLevel, outputType)
		},
	}
	return execFunctions
}

func defineEnvironmentVariables(yamlDocument map[interface{}][]interface{}) (map[string]string, []string) {
	var (
		environmentVariablesShell []string
		tempEnvironmentVariables  map[string]string
	)
	// define temp function to convert variable to key value map
	getEnvironmentVariables := func(data []string, getkeyval func(item string) (key, val string)) map[string]string {
		items := make(map[string]string)
		for _, item := range data {
			key, val := getkeyval(item)
			items[key] = val
		}
		return items
	}
	// parse os environment variables and covert it to key value map
	tempEnvironmentVariables = getEnvironmentVariables(os.Environ(), func(item string) (key, val string) {
		splits := strings.Split(item, "=")
		key = splits[0]
		val = splits[1]
		return
	})

	// read all defined environment variables as key value map and store it in string slice as variable definiton and additionaly store parsed variables to temporary variables map
	for key := range yamlDocument["env"] {
		_env, ok := yamlDocument["env"][key].(map[interface{}]interface{})
		if !ok {
			functions.PrintSwitch(color.FgRed, "error", "file", "it was not successfull to read yamlDocument['env'][key]")
		}
		envkey, ok := _env["key"].(string)
		if !ok {
			functions.PrintSwitch(color.FgRed, "error", "file", "it was not successfull to read _env['key'].(string)")
		}
		envvalue, ok := _env["value"].(string)
		if !ok {
			functions.PrintSwitch(color.FgRed, "error", "file", "it was not successfull to read _env['value'].(string)")
		}
		os.Setenv(envkey, envvalue)
		environmentVariablesShell = append(environmentVariablesShell, envkey+"="+envvalue)
		tempEnvironmentVariables[envkey] = envvalue
	}
	return tempEnvironmentVariables, environmentVariablesShell
}

// parse all environment variables: defined by runfromyaml logging block or from current shell environment
func loggingSettings(yamlDocument map[interface{}][]interface{}) (string, string) {
	var (
		outputType  string
		outputLevel string
	)

	// parse logging block
	for key := range yamlDocument["logging"] {

		// store current element in setting
		setting := yamlDocument["logging"][key].(map[interface{}]interface{})

		// set outputType
		if reflect.ValueOf(setting["output"]).IsValid() {
			outputType = setting["output"].(string)
		}

		// set outputLevel
		if reflect.ValueOf(setting["level"]).IsValid() {
			outputLevel = setting["level"].(string)
		}
	}
	return outputType, outputLevel
}

func Runfromyaml(yamlFile []byte, debug bool) {
	// define all the variables
	var (
		outputType                string
		outputLevel               string
		yamlDocument              map[interface{}][]interface{}
		yamlBlock                 map[interface{}]interface{}
		environmentVariablesShell []string
		ok                        bool
	)

	// define functions map
	execFunctions := execFunctionsMap(yamlBlock, environmentVariablesShell, outputLevel, outputType)

	// parse YAML structure
	if err := yaml.Unmarshal(yamlFile, &yamlDocument); err != nil {
		functions.PrintSwitch(color.FgHiWhite, "error", "file", "could not unmarshal YAML data ("+err.Error()+")")
	}

	// set output logging settings
	outputType, outputLevel = loggingSettings(yamlDocument)

	// print out filename if output type is as file defined
	if outputType == "file" {
		if debug {
			functions.PrintSwitch(color.FgHiWhite, "info", "stdout", "logfile temp file: "+os.TempDir()+"logrus-"+time.Now().Format("20060102")+".log")
		}
	}

	// parse environment variables and store all variables to EnvironmentVariables and only defined variables store to environmentVariablesShell
	EnvironmentVariables, environmentVariablesShell = defineEnvironmentVariables(yamlDocument)

	// parse command yaml blocks
	for key := range yamlDocument["cmd"] {

		// if the main command block is not empty, then apply all the steps
		if !reflect.ValueOf(yamlDocument["cmd"][key].(map[interface{}]interface{})).IsNil() {
			// read the structure and store it to yamlBlock variable otherwise print the error
			yamlBlock, ok = yamlDocument["cmd"][key].(map[interface{}]interface{})
			if !ok {
				functions.PrintSwitch(color.FgRed, "error", "file", "it was not successfull to read yamlDocument['cmd'][key].(map[interface{}]interface{}))")
			}

			// if desciption block is not empty try to create from description AI generated command and print everything as comment
			if reflect.ValueOf(yamlBlock["desc"]).IsValid() {
				var aidesc string
				if openai.IsAiEnabled {
					for {
						response, err := openai.OpenAI(openai.Key, openai.Model, functions.EvaluateDescription(yamlBlock), openai.ShellType)
						if err == nil {
							aidesc = "# example: " + openai.PrintAiResponse(response) + "\n"
							break
						}
					}
					functions.PrintSwitch(color.FgHiCyan, outputLevel, outputType, aidesc)
				}
				functions.PrintSwitch(color.FgGreen, outputLevel, outputType, "\n"+functions.EvaluateDescription(yamlBlock))
			}

			// forward all the collected options for specific command and execute specific command defined by function map
			execFunctions[yamlBlock["type"].(string)](yamlBlock, environmentVariablesShell, outputLevel, outputType)
		}
	}
}

func InteractiveShell(shell string) []string {
	bash := []string{shell, "--login"}
	cmd := execute.Command(bash[0], bash[1:]...)
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()

	var commands []string

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print("your session will be recorded > ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		if input != "" && input != "exit" {
			commands = append(commands, input)
		}
		_, err = stdin.Write([]byte(input + "\n"))
		if err != nil {
			break
		}
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Println(err)
	}

	return commands

}
