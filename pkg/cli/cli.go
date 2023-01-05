package cli

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/fatih/color"
	"gopkg.in/yaml.v2"

	"github.com/lanixx/runfromyaml/pkg/exec"
	functions "github.com/lanixx/runfromyaml/pkg/functions"
)

var (
	ydoc       map[interface{}][]interface{}
	yblock     map[interface{}]interface{}
	values     string
	envs       []string
	cmds       []string
	dcoptions  []string
	cmdoptions []string
	options    []string
	desc       string
	confdata   string
	confdest   string
	confperm   os.FileMode
	debug      bool
)

func execCmd(yblock map[interface{}]interface{}, _envs []string) {
	wg := new(sync.WaitGroup)

	if !reflect.ValueOf(yblock["values"]).IsNil() {
		values = strings.Trim(fmt.Sprint(yblock["values"]), "[]")
		if reflect.ValueOf(yblock["expandenv"]).Bool() && yblock["expandenv"].(bool) {
			values = os.ExpandEnv(values)
			if debug {
				color.New(color.FgHiBlack).Println("# environment variables are expanded")
			}
		}
		cmds = strings.Fields(values)
	}
	if string(yblock["desc"].(string)) != "" {
		desc = fmt.Sprintf("%v", yblock["desc"])
	}
	wg.Add(1)
	go exec.Command(cmds, desc, wg, _envs)
	wg.Wait()
}

func shellCmd(yblock map[interface{}]interface{}, _envs []string) {
	if !reflect.ValueOf(yblock["values"]).IsNil() {
		values = strings.Trim(fmt.Sprint(yblock["values"]), "[]")
		if reflect.ValueOf(yblock["expandenv"]).Bool() && yblock["expandenv"].(bool) {
			values = os.ExpandEnv(values)
			if debug {
				color.New(color.FgHiBlack).Println("# environment variables are expanded")
			}
		}
		cmds = strings.Fields(values)
	}
	if string(yblock["desc"].(string)) != "" {
		desc = fmt.Sprintf("%v", yblock["desc"])
	}
	temp_cmds := strings.Join(cmds, " ")
	cmds := strings.Split(temp_cmds, ";")
	wg := new(sync.WaitGroup)
	for ind, shcmds := range cmds {
		shcmd := strings.Split(shcmds, " ")
		wg.Add(1)
		go exec.CommandShell(shcmd, desc, wg, ind, _envs)
		wg.Wait()
	}
}

func dockerCmd(yblock map[interface{}]interface{}, _envs []string) {
	wg := new(sync.WaitGroup)
	if !reflect.ValueOf(yblock["values"]).IsNil() {
		values = strings.Trim(fmt.Sprint(yblock["values"]), "[]")
		if reflect.ValueOf(yblock["expandenv"]).Bool() && yblock["expandenv"].(bool) {
			values = os.ExpandEnv(values)
			if debug {
				color.New(color.FgHiBlack).Println("# environment variables are expanded")
			}
		}
		cmds = strings.Fields(values)
	}
	if string(yblock["desc"].(string)) != "" {
		desc = fmt.Sprintf("%v", yblock["desc"])
	}
	wg.Add(1)
	go exec.CommandDockerRun(yblock["command"].(string), yblock["container"].(string), cmds, desc, _envs, wg)
	wg.Wait()
}

func dockerComposeCmd(yblock map[interface{}]interface{}, _envs []string) {
	var service string
	wg := new(sync.WaitGroup)
	if !reflect.ValueOf(yblock["values"]).IsNil() {
		values = strings.Trim(fmt.Sprint(yblock["values"]), "[]")
		if reflect.ValueOf(yblock["expandenv"]).Bool() && yblock["expandenv"].(bool) {
			values = os.ExpandEnv(values)
			if debug {
				color.New(color.FgHiBlack).Println("# environment variables are expanded")
			}
		}
		cmds = strings.Fields(values)
	}
	if !reflect.ValueOf(yblock["dcoptions"]).IsNil() {
		values = strings.Trim(fmt.Sprint(yblock["dcoptions"]), "[]")
		if reflect.ValueOf(yblock["expandenv"]).Bool() && yblock["expandenv"].(bool) {
			values = os.ExpandEnv(values)
		}
		dcoptions = strings.Fields(values)
	}
	if !reflect.ValueOf(yblock["cmdoptions"]).IsNil() {
		values = strings.Trim(fmt.Sprint(yblock["cmdoptions"]), "[]")
		if reflect.ValueOf(yblock["expandenv"]).Bool() && yblock["expandenv"].(bool) {
			values = os.ExpandEnv(values)
		}
		cmdoptions = strings.Fields(values)
	}
	if string(yblock["desc"].(string)) != "" {
		desc = fmt.Sprintf("%v", yblock["desc"])
	}
	if string(yblock["service"].(string)) != "" {
		service = fmt.Sprintf("%v", yblock["service"])
	}
	wg.Add(1)
	go exec.CommandDockerComposeExec(yblock["command"].(string), service, cmdoptions, dcoptions, cmds, desc, _envs, wg)
	wg.Wait()
}

func sshCmd(yblock map[interface{}]interface{}, _envs []string) {
	var user string
	var host string
	wg := new(sync.WaitGroup)
	if !reflect.ValueOf(yblock["values"]).IsNil() {
		values = strings.Trim(fmt.Sprint(yblock["values"]), "[]")
		if reflect.ValueOf(yblock["expandenv"]).Bool() && yblock["expandenv"].(bool) {
			values = os.ExpandEnv(values)
			if debug {
				color.New(color.FgHiBlack).Println("# environment variables are expanded")
			}
		}
		cmds = strings.Fields(values)
	}
	if !reflect.ValueOf(yblock["options"]).IsNil() {
		values = strings.Trim(fmt.Sprint(yblock["options"]), "[]")
		if reflect.ValueOf(yblock["expandenv"]).Bool() && yblock["expandenv"].(bool) {
			values = os.ExpandEnv(values)
		}
		options = strings.Fields(values)
	}
	if string(yblock["desc"].(string)) != "" {
		desc = fmt.Sprintf("%v", yblock["desc"])
	}
	if reflect.ValueOf(yblock["expandenv"]).Bool() && yblock["expandenv"].(bool) {
		user = os.ExpandEnv(yblock["user"].(string))
	} else {
		user = yblock["user"].(string)
	}
	if reflect.ValueOf(yblock["expandenv"]).Bool() && yblock["expandenv"].(bool) {
		host = os.ExpandEnv(yblock["host"].(string))
	} else {
		host = yblock["host"].(string)
	}
	wg.Add(1)
	go exec.CommandSSH(user, yblock["port"].(int), host, options, cmds, desc, _envs, wg)
	wg.Wait()
}

func conf(yblock map[interface{}]interface{}) {
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
		//readFile(string(confdest))
	}
	color.New(color.FgGreen).Println("# create " + confdest)
}

func Runfromyaml(yamlFile []byte) {

	yaml.Unmarshal(yamlFile, &ydoc)

	if debug {
		functions.PrintColor(color.BgCyan, ydoc["env"])
	}

	for key := range ydoc["env"] {
		_env := ydoc["env"][key].(map[interface{}]interface{})
		os.Setenv(_env["key"].(string), _env["value"].(string))
		envs = append(envs, _env["key"].(string)+"="+_env["value"].(string))
		if debug {
			functions.PrintColor(color.FgCyan, _env["key"].(string)+"="+_env["value"].(string))
		}
	}

	for key := range ydoc["cmd"] {

		if debug {
			color.New(color.FgBlue).Println("\n" + "# " + fmt.Sprint(key+1))
		}

		if !reflect.ValueOf(ydoc["cmd"][key].(map[interface{}]interface{})).IsNil() {
			yblock = ydoc["cmd"][key].(map[interface{}]interface{})
			if debug {
				functions.PrintColor(color.FgHiCyan, "\n\n%+v\n\n", yblock)
				functions.PrintColor(color.FgBlue, "Name: %+v\n", yblock["name"])
				functions.PrintColor(color.FgBlue, "Beschreibung: %+v\n", yblock["desc"])
				functions.PrintColor(color.FgBlue, "Key: %+v\n", key)
				functions.PrintColor(color.FgBlue, "Command: %+v\n", values)
				functions.PrintColor(color.FgBlue, "Data:\n---\n%+v\n---\n", yblock["confdata"])
				functions.PrintColor(color.FgBlue, "Destination: %+v\n", yblock["confdest"])
				functions.PrintColor(color.FgBlue, "Permissions: %+v\n", yblock["confperm"])
				functions.PrintColor(color.FgHiWhite, "Key: %+v\n", key)
				functions.PrintColor(color.FgGreen, "Name: %+v\n", yblock["name"])
				functions.PrintColor(color.FgGreen, "Beschreibung: %+v\n", yblock["desc"])
				functions.PrintColor(color.FgYellow, "Command: %+v\n", values)
				fmt.Printf("\n")
				fmt.Printf("\n")
			}

			if reflect.ValueOf(yblock["desc"]).IsValid() {
				color.New(color.FgGreen).Println("\n" + "# " + yblock["desc"].(string))
			}

			if yblock["type"] == "exec" {
				execCmd(yblock, envs)
			}
			if yblock["type"] == "shell" {
				shellCmd(yblock, envs)
			}
			if yblock["type"] == "docker" {
				dockerCmd(yblock, envs)
			}
			if yblock["type"] == "docker-compose" {
				dockerComposeCmd(yblock, envs)
			}
			if yblock["type"] == "ssh" {
				sshCmd(yblock, envs)
			}
			if yblock["type"] == "conf" {
				if debug {
					functions.PrintColor(color.FgYellow, "Config: %+v\n", yblock["confdata"])
					functions.PrintColor(color.FgYellow, "Config: %+v\n", yblock["confdest"])
					functions.PrintColor(color.FgYellow, "Config: %+v\n", yblock["confperm"])
					fmt.Printf("\n")
				}
				conf(yblock)
			}
		}
	}
}
