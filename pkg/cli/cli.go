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

func execCmd(yamlBlock map[interface{}]interface{}, _envvars []string, _level string, _output string) {
	var (
		cmds []string
		desc string
	)
	wg := new(sync.WaitGroup)

	cmds = functions.ExtractAndExpand(yamlBlock, "values")

	wg.Add(1)
	go exec.Command(cmds, desc, wg, _envvars, _level, _output)
	wg.Wait()
}

func shellCmd(yamlBlock map[interface{}]interface{}, _envvars []string, _level string, _output string) {
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
		go exec.CommandShell(shcmd, desc, wg, ind, _envvars, _level, _output)
		wg.Wait()
	}
}

func dockerCmd(yamlBlock map[interface{}]interface{}, _envvars []string, _level string, _output string) {
	var (
		cmds []string
		desc string
	)
	wg := new(sync.WaitGroup)
	cmds = functions.ExtractAndExpand(yamlBlock, "values")

	wg.Add(1)
	go exec.CommandDockerRun(yamlBlock["command"].(string), yamlBlock["container"].(string), cmds, desc, _envvars, wg, _level, _output)
	wg.Wait()
}

func dockerComposeCmd(yamlBlock map[interface{}]interface{}, _envvars []string, _level string, _output string) {
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
	go exec.CommandDockerComposeExec(yamlBlock["command"].(string), service, cmdoptions, dcoptions, cmds, desc, _envvars, wg, _level, _output)
	wg.Wait()
}

func sshCmd(yamlBlock map[interface{}]interface{}, _envvars []string, _level string, _output string) {
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
	go exec.CommandSSH(user, yamlBlock["port"].(int), host, options, cmds, desc, _envvars, wg, _level, _output)
	wg.Wait()
}

func conf(yamlBlock map[interface{}]interface{}, _level string, _output string) {
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
	confdata = "# " + desc + "\n" + confdata

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
	functions.PrintSwitch(color.FgGreen, _level, _output, "# create ", confdest)
}

func Runfromyaml(yamlFile []byte, debug bool) {
	var (
		_output                   string
		_level                    string
		yamlDocument              map[interface{}][]interface{}
		yamlBlock                 map[interface{}]interface{}
		environmentVariablesShell []string
		ok                        bool
	)

	if err := yaml.Unmarshal(yamlFile, &yamlDocument); err != nil {
		functions.PrintSwitch(color.FgHiWhite, "info", "stdout", "could not unmarshal YAML data ("+err.Error()+")")
	}

	for key := range yamlDocument["logging"] {
		setting := yamlDocument["logging"][key].(map[interface{}]interface{})
		if reflect.ValueOf(setting["output"]).IsValid() {
			_output = setting["output"].(string)
		}
		if reflect.ValueOf(setting["level"]).IsValid() {
			_level = setting["level"].(string)
		}
	}

	if _output == "file" {
		functions.PrintSwitch(color.FgHiWhite, "info", "stdout", "logfile temp file: "+os.TempDir()+"logrus-"+time.Now().Format("20060102")+".log")
	}

	getEnvironmentVariables := func(data []string, getkeyval func(item string) (key, val string)) map[string]string {
		items := make(map[string]string)
		for _, item := range data {
			key, val := getkeyval(item)
			items[key] = val
		}
		return items
	}
	tempEnvironmentVariables := getEnvironmentVariables(os.Environ(), func(item string) (key, val string) {
		splits := strings.Split(item, "=")
		key = splits[0]
		val = splits[1]
		return
	})

	for key := range yamlDocument["env"] {
		_env, ok := yamlDocument["env"][key].(map[interface{}]interface{})
		if !ok {
			functions.PrintSwitch(color.FgRed, _level, _output, "it was not successfull to read yamlDocument['env'][key]")
		}
		envkey, ok := _env["key"].(string)
		if !ok {
			functions.PrintSwitch(color.FgRed, _level, _output, "it was not successfull to read _env['key'].(string)")
		}
		envvalue, ok := _env["value"].(string)
		if !ok {
			functions.PrintSwitch(color.FgRed, _level, _output, "it was not successfull to read _env['value'].(string)")
		}
		os.Setenv(envkey, envvalue)
		environmentVariablesShell = append(environmentVariablesShell, envkey+"="+envvalue)
		tempEnvironmentVariables[envkey] = envvalue
	}
	EnvironmentVariables = tempEnvironmentVariables

	for key := range yamlDocument["cmd"] {

		if !reflect.ValueOf(yamlDocument["cmd"][key].(map[interface{}]interface{})).IsNil() {
			yamlBlock, ok = yamlDocument["cmd"][key].(map[interface{}]interface{})
			if !ok {
				functions.PrintSwitch(color.FgRed, _level, _output, "it was not successfull to read yamlDocument['cmd'][key].(map[interface{}]interface{}))")
			}

			if reflect.ValueOf(yamlBlock["desc"]).IsValid() {
				var aidesc string
				if openai.IsAiEnabled {
					for {
						response, err := openai.OpenAI(openai.Key, openai.Model, functions.EvaluateDescription(yamlBlock), openai.ShellType)
						if err == nil {
							aidesc = "\n" + "# example: " + openai.PrintAiResponse(response)
							break
						}
					}
				}
				functions.PrintSwitch(color.FgGreen, _level, _output, "\n"+"# "+functions.EvaluateDescription(yamlBlock)+aidesc)
			}

			if yamlBlock["type"] == "exec" {
				execCmd(yamlBlock, environmentVariablesShell, _level, _output)
			}
			if yamlBlock["type"] == "shell" {
				shellCmd(yamlBlock, environmentVariablesShell, _level, _output)
			}
			if yamlBlock["type"] == "docker" {
				dockerCmd(yamlBlock, environmentVariablesShell, _level, _output)
			}
			if yamlBlock["type"] == "docker-compose" {
				dockerComposeCmd(yamlBlock, environmentVariablesShell, _level, _output)
			}
			if yamlBlock["type"] == "ssh" {
				sshCmd(yamlBlock, environmentVariablesShell, _level, _output)
			}
			if yamlBlock["type"] == "conf" {
				if debug {
					functions.PrintSwitch(color.FgHiBlue, "debug", _output, "Destination: %+v\n", yamlBlock["confdest"])
					functions.PrintSwitch(color.FgHiBlue, "debug", _output, "Permissions: %+v\n", yamlBlock["confperm"])
					functions.PrintSwitch(color.FgHiBlue, "debug", _output, "Data:\n---\n%+v\n---\n", yamlBlock["confdata"])
				}
				conf(yamlBlock, _level, _output)
			}
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
