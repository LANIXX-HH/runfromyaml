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
)

var (
	envstruct map[string]string
)

func execCmd(yblock map[interface{}]interface{}, _envs []string, _level string, _output string) {
	var (
		cmds []string
		desc string
	)
	wg := new(sync.WaitGroup)

	cmds = functions.ExtractAndExpand(yblock, "values")

	if reflect.ValueOf(yblock["desc"]).IsValid() {
		desc = fmt.Sprintf("%v", yblock["desc"])
	}
	wg.Add(1)
	go exec.Command(cmds, desc, wg, _envs, _level, _output)
	wg.Wait()
}

func shellCmd(yblock map[interface{}]interface{}, _envs []string, _level string, _output string) {
	var (
		cmds []string
		desc string
	)
	cmds = functions.ExtractAndExpand(yblock, "values")

	desc = "<no description>"
	if reflect.ValueOf(yblock["desc"]).IsValid() {
		desc = fmt.Sprintf("%v", yblock["desc"])
	}
	temp_cmds := strings.Join(cmds, " ")
	cmds = strings.Split(temp_cmds, ";")
	wg := new(sync.WaitGroup)
	for ind, shcmds := range cmds {
		shcmd := strings.Split(shcmds, " ")
		wg.Add(1)
		go exec.CommandShell(shcmd, desc, wg, ind, _envs, _level, _output)
		wg.Wait()
	}
}

func dockerCmd(yblock map[interface{}]interface{}, _envs []string, _level string, _output string) {
	var (
		cmds []string
		desc string
	)
	wg := new(sync.WaitGroup)
	cmds = functions.ExtractAndExpand(yblock, "values")

	desc = "<no description>"
	if reflect.ValueOf(yblock["desc"]).IsValid() {
		desc = fmt.Sprintf("%v", yblock["desc"])
	}
	wg.Add(1)
	go exec.CommandDockerRun(yblock["command"].(string), yblock["container"].(string), cmds, desc, _envs, wg, _level, _output)
	wg.Wait()
}

func dockerComposeCmd(yblock map[interface{}]interface{}, _envs []string, _level string, _output string) {
	var (
		cmds       []string
		dcoptions  []string
		cmdoptions []string
		desc       string
		service    string
	)
	wg := new(sync.WaitGroup)
	cmds = functions.ExtractAndExpand(yblock, "values")
	dcoptions = functions.ExtractAndExpand(yblock, "dcoptions")
	cmdoptions = functions.ExtractAndExpand(yblock, "cmdoptions")

	desc = "<no description>"
	if reflect.ValueOf(yblock["desc"]).IsValid() {
		desc = yblock["desc"].(string)
	}
	service = "no-service"
	if reflect.ValueOf(yblock["service"]).IsValid() {
		service = yblock["service"].(string)
	}
	wg.Add(1)
	go exec.CommandDockerComposeExecNew(yblock["command"].(string), service, cmdoptions, dcoptions, cmds, desc, _envs, wg, _level, _output)
	wg.Wait()
}

func sshCmd(yblock map[interface{}]interface{}, _envs []string, _level string, _output string) {
	var (
		cmds    []string
		options []string
		desc    string
	)
	var user string
	var host string
	wg := new(sync.WaitGroup)
	cmds = functions.ExtractAndExpand(yblock, "values")
	options = functions.ExtractAndExpand(yblock, "options")

	desc = "<no description>"
	if reflect.ValueOf(yblock["desc"]).IsValid() {
		desc = fmt.Sprintf("%v", yblock["desc"])
	}

	user = yblock["user"].(string)
	if reflect.ValueOf(yblock["expandenv"]).IsValid() {
		if yblock["expandenv"].(bool) {
			user = os.ExpandEnv(user)
		}
	}

	host = yblock["host"].(string)
	if reflect.ValueOf(yblock["expandenv"]).IsValid() {
		if yblock["expandenv"].(bool) {
			host = os.ExpandEnv(host)
		}
	}

	wg.Add(1)
	go exec.CommandSSH(user, yblock["port"].(int), host, options, cmds, desc, _envs, wg, _level, _output)
	wg.Wait()
}

func conf(yblock map[interface{}]interface{}, _level string, _output string) {
	var (
		confdata string
		confdest string
		confperm os.FileMode
	)
	if reflect.ValueOf(yblock["confdata"].(string)).String() != "" {
		confdata = yblock["confdata"].(string)
		if reflect.ValueOf(yblock["expandenv"]).IsValid() {
			if yblock["expandenv"].(bool) {
				confdata = functions.GoTemplate(envstruct, confdata)
			}
		}
	}
	if reflect.ValueOf(yblock["confdest"].(string)).String() != "" {
		confdest = yblock["confdest"].(string)
	}
	if reflect.ValueOf(yblock["confperm"].(int)).Int() != 0 {
		confperm = os.FileMode(int(yblock["confperm"].(int)))
	}
	if confdata != "" && confdest != "" && string(rune(confperm)) != "" {
		functions.WriteFile(confdata, confdest, confperm)
	}
	if reflect.ValueOf(yblock["expandenv"]).IsValid() {
		if yblock["expandenv"].(bool) {
			confdest = os.ExpandEnv(confdest)
		}
	}
	functions.PrintSwitch(color.FgGreen, _level, _output, "# create ", confdest)
}

func Runfromyaml(yamlFile []byte, debug bool) {
	var (
		_output string
		_level  string
		ydoc    map[interface{}][]interface{}
		yblock  map[interface{}]interface{}
		envs    []string
		ok      bool
	)

	if err := yaml.Unmarshal(yamlFile, &ydoc); err != nil {
		functions.PrintSwitch(color.FgHiWhite, "info", "stdout", "could not unmarshal YAML data ("+err.Error()+")")
	}

	for key := range ydoc["logging"] {
		setting := ydoc["logging"][key].(map[interface{}]interface{})
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

	getenvironment := func(data []string, getkeyval func(item string) (key, val string)) map[string]string {
		items := make(map[string]string)
		for _, item := range data {
			key, val := getkeyval(item)
			items[key] = val
		}
		return items
	}
	environment := getenvironment(os.Environ(), func(item string) (key, val string) {
		splits := strings.Split(item, "=")
		key = splits[0]
		val = splits[1]
		return
	})

	for key := range ydoc["env"] {
		_env, ok := ydoc["env"][key].(map[interface{}]interface{})
		if !ok {
			functions.PrintSwitch(color.FgRed, _level, _output, "it was not successfull to read ydoc['env'][key]")
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
		envs = append(envs, envkey+"="+envvalue)
		environment[envkey] = envvalue
	}
	envstruct = environment

	for key := range ydoc["cmd"] {

		if !reflect.ValueOf(ydoc["cmd"][key].(map[interface{}]interface{})).IsNil() {
			yblock, ok = ydoc["cmd"][key].(map[interface{}]interface{})
			if !ok {
				functions.PrintSwitch(color.FgRed, _level, _output, "it was not successfull to read ydoc['cmd'][key].(map[interface{}]interface{}))")
			}

			if reflect.ValueOf(yblock["desc"]).IsValid() {
				functions.PrintSwitch(color.FgGreen, _level, _output, "\n"+"# "+yblock["desc"].(string))
			}

			if yblock["type"] == "exec" {
				execCmd(yblock, envs, _level, _output)
			}
			if yblock["type"] == "shell" {
				shellCmd(yblock, envs, _level, _output)
			}
			if yblock["type"] == "docker" {
				dockerCmd(yblock, envs, _level, _output)
			}
			if yblock["type"] == "docker-compose" {
				dockerComposeCmd(yblock, envs, _level, _output)
			}
			if yblock["type"] == "ssh" {
				sshCmd(yblock, envs, _level, _output)
			}
			if yblock["type"] == "conf" {
				if debug {
					functions.PrintSwitch(color.FgHiBlue, "debug", _output, "Destination: %+v\n", yblock["confdest"])
					functions.PrintSwitch(color.FgHiBlue, "debug", _output, "Permissions: %+v\n", yblock["confperm"])
					functions.PrintSwitch(color.FgHiBlue, "debug", _output, "Data:\n---\n%+v\n---\n", yblock["confdata"])
				}
				conf(yblock, _level, _output)
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
