package cli

import (
	"fmt"
	"os"
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
	debug bool
)

func execCmd(yblock map[interface{}]interface{}, _envs []string, _level string, _output string) {
	var (
		values string
		cmds   []string
		desc   string
	)
	wg := new(sync.WaitGroup)

	if reflect.ValueOf(yblock["values"]).IsValid() {
		values = strings.Trim(fmt.Sprint(yblock["values"]), "[]")
		if reflect.ValueOf(yblock["expandenv"]).IsValid() {
			if yblock["expandenv"].(bool) {
				values = os.ExpandEnv(values)
				if debug {
					functions.PrintColor(color.FgHiBlack, "debug", _output, "# environment variables are expanded")
				}
			}
		}
		cmds = strings.Fields(values)
	}
	if reflect.ValueOf(yblock["desc"]).IsValid() {
		desc = fmt.Sprintf("%v", yblock["desc"])
	}
	wg.Add(1)
	go exec.Command(cmds, desc, wg, _envs, _level, _output)
	wg.Wait()
}

func shellCmd(yblock map[interface{}]interface{}, _envs []string, _level string, _output string) {
	var (
		values string
		cmds   []string
		desc   string
	)
	if reflect.ValueOf(yblock["values"]).IsValid() {
		values = strings.Trim(fmt.Sprint(yblock["values"]), "[]")
		if reflect.ValueOf(yblock["expandenv"]).IsValid() {
			if yblock["expandenv"].(bool) {
				values = os.ExpandEnv(values)
				if debug {
					functions.PrintColor(color.FgHiBlack, "debug", _output, "# environment variables are expanded")
				}
			}
		}
		cmds = strings.Fields(values)
	}
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
		values string
		cmds   []string
		desc   string
	)
	wg := new(sync.WaitGroup)
	if reflect.ValueOf(yblock["values"]).IsValid() {
		values = strings.Trim(fmt.Sprint(yblock["values"]), "[]")
		if reflect.ValueOf(yblock["expandenv"]).IsValid() {
			if yblock["expandenv"].(bool) {
				values = os.ExpandEnv(values)
				if debug {
					functions.PrintColor(color.FgHiBlack, "debug", _output, "# environment variables are expanded")
				}
			}
		}
		cmds = strings.Fields(values)
	}
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
		values     string
		cmds       []string
		dcoptions  []string
		cmdoptions []string
		desc       string
		service    string
	)
	wg := new(sync.WaitGroup)
	if reflect.ValueOf(yblock["values"]).IsValid() {
		values = strings.Trim(fmt.Sprint(yblock["values"]), "[]")
		if reflect.ValueOf(yblock["expandenv"]).IsValid() {
			if yblock["expandenv"].(bool) {
				values = os.ExpandEnv(values)
				if debug {
					functions.PrintColor(color.FgHiBlack, "debug", _output, "# environment variables are expanded")
				}
			}
		}
		cmds = strings.Fields(values)
	}
	if reflect.ValueOf(yblock["dcoptions"]).IsValid() {
		values = strings.Trim(fmt.Sprint(yblock["dcoptions"]), "[]")
		if reflect.ValueOf(yblock["expandenv"]).IsValid() {
			if yblock["expandenv"].(bool) {
				values = os.ExpandEnv(values)
			}
		}
		dcoptions = strings.Fields(values)
	}
	if reflect.ValueOf(yblock["cmdoptions"]).IsValid() {
		values = strings.Trim(fmt.Sprint(yblock["cmdoptions"]), "[]")
		if reflect.ValueOf(yblock["expandenv"]).IsValid() {
			if yblock["expandenv"].(bool) {
				values = os.ExpandEnv(values)
			}
		}
		cmdoptions = strings.Fields(values)
	}
	desc = "<no description>"
	if reflect.ValueOf(yblock["desc"]).IsValid() {
		desc = fmt.Sprintf("%v", yblock["desc"])
	}
	service = "no-service"
	if reflect.ValueOf(yblock["service"]).IsValid() {
		service = fmt.Sprintf("%v", yblock["service"])
	}
	wg.Add(1)
	go exec.CommandDockerComposeExec(yblock["command"].(string), service, cmdoptions, dcoptions, cmds, desc, _envs, wg, _level, _output)
	wg.Wait()
}

func sshCmd(yblock map[interface{}]interface{}, _envs []string, _level string, _output string) {
	var (
		values  string
		cmds    []string
		options []string
		desc    string
	)
	var user string
	var host string
	wg := new(sync.WaitGroup)
	if reflect.ValueOf(yblock["values"]).IsValid() {
		values = strings.Trim(fmt.Sprint(yblock["values"]), "[]")
		if reflect.ValueOf(yblock["expandenv"]).IsValid() {
			if yblock["expandenv"].(bool) {
				values = os.ExpandEnv(values)
				if debug {
					functions.PrintColor(color.FgHiBlack, "debug", _output, "# environment variables are expanded")
				}
			}
		}
		cmds = strings.Fields(values)
	}
	if reflect.ValueOf(yblock["options"]).IsValid() {
		values = strings.Trim(fmt.Sprint(yblock["options"]), "[]")
		if reflect.ValueOf(yblock["expandenv"]).IsValid() {
			if yblock["expandenv"].(bool) {
				values = os.ExpandEnv(values)
			}
		}
		options = strings.Fields(values)
	}
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
	functions.PrintColor(color.FgGreen, _level, _output, "# create ", confdest)
}

func Runfromyaml(yamlFile []byte, debug bool) {
	var (
		_output string
		_level  string
		ydoc    map[interface{}][]interface{}
		yblock  map[interface{}]interface{}
		envs    []string
	)

	yaml.Unmarshal(yamlFile, &ydoc)

	for key := range ydoc["logging"] {
		setting := ydoc["logging"][key].(map[interface{}]interface{})
		if reflect.ValueOf(setting["output"]).IsValid() {
			_output = fmt.Sprint(setting["output"])
		}
		if reflect.ValueOf(setting["level"]).IsValid() {
			_level = fmt.Sprint(setting["level"])
		}
	}
	if _output == "file" {
		fmt.Println("logfile temp file: " + os.TempDir() + "logrus-" + time.Now().Format("20060102") + ".log")
	}
	if debug {
		functions.PrintColor(color.FgRed, "debug", _output, ydoc["env"])
	}

	for key := range ydoc["env"] {
		_env := ydoc["env"][key].(map[interface{}]interface{})
		os.Setenv(_env["key"].(string), _env["value"].(string))
		envs = append(envs, _env["key"].(string)+"="+_env["value"].(string))
		if debug {
			functions.PrintColor(color.FgRed, "debug", _output, _env["key"].(string)+"="+_env["value"].(string))
		}
	}

	for key := range ydoc["cmd"] {

		if debug {
			functions.PrintColor(color.FgHiBlue, "debug", _output, "\n"+"# "+fmt.Sprint(key+1))
		}

		if !reflect.ValueOf(ydoc["cmd"][key].(map[interface{}]interface{})).IsNil() {
			yblock = ydoc["cmd"][key].(map[interface{}]interface{})
			if debug {
				functions.PrintColor(color.FgHiBlue, "debug", _output, "YAML Block: \n---\n", yblock, "\n---\n")
			}

			if reflect.ValueOf(yblock["desc"]).IsValid() {
				functions.PrintColor(color.FgGreen, _level, _output, "\n"+"# "+yblock["desc"].(string))
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
					functions.PrintColor(color.FgHiBlue, "debug", _output, "Destination: %+v\n", yblock["confdest"])
					functions.PrintColor(color.FgHiBlue, "debug", _output, "Permissions: %+v\n", yblock["confperm"])
					functions.PrintColor(color.FgHiBlue, "debug", _output, "Data:\n---\n%+v\n---\n", yblock["confdata"])
				}
				conf(yblock, _level, _output)
			}
		}
	}
}
