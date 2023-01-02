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
	results    map[interface{}][]interface{}
	types      map[interface{}]interface{}
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

func execCmd(types map[interface{}]interface{}, _envs []string) {
	wg := new(sync.WaitGroup)

	if !reflect.ValueOf(types["values"]).IsNil() {
		values = strings.Trim(fmt.Sprint(types["values"]), "[]")
		if reflect.ValueOf(types["expandenv"]).Bool() && types["expandenv"].(bool) {
			values = os.ExpandEnv(values)

		}
		cmds = strings.Fields(values)
	}
	if string(types["desc"].(string)) != "" {
		desc = fmt.Sprintf("%v", types["desc"])
	}
	wg.Add(1)
	go exec.Command(cmds, desc, wg, _envs)
	wg.Wait()
}

func shellCmd(types map[interface{}]interface{}, _envs []string) {
	if !reflect.ValueOf(types["values"]).IsNil() {
		values = strings.Trim(fmt.Sprint(types["values"]), "[]")
		if reflect.ValueOf(types["expandenv"]).Bool() && types["expandenv"].(bool) {
			values = os.ExpandEnv(values)

		}
		cmds = strings.Fields(values)
	}
	if string(types["desc"].(string)) != "" {
		desc = fmt.Sprintf("%v", types["desc"])
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

func dockerCmd(types map[interface{}]interface{}, _envs []string) {
	wg := new(sync.WaitGroup)
	if !reflect.ValueOf(types["values"]).IsNil() {
		values = strings.Trim(fmt.Sprint(types["values"]), "[]")
		if reflect.ValueOf(types["expandenv"]).Bool() && types["expandenv"].(bool) {
			values = os.ExpandEnv(values)

		}
		cmds = strings.Fields(values)
	}
	if string(types["desc"].(string)) != "" {
		desc = fmt.Sprintf("%v", types["desc"])
	}
	wg.Add(1)
	go exec.CommandDockerRun(types["command"].(string), types["container"].(string), cmds, desc, _envs, wg)
	wg.Wait()
}

func dockerComposeCmd(types map[interface{}]interface{}, _envs []string) {
	var service string
	wg := new(sync.WaitGroup)
	if !reflect.ValueOf(types["values"]).IsNil() {
		values = strings.Trim(fmt.Sprint(types["values"]), "[]")
		if reflect.ValueOf(types["expandenv"]).Bool() && types["expandenv"].(bool) {
			values = os.ExpandEnv(values)

		}
		cmds = strings.Fields(values)
	}
	if !reflect.ValueOf(types["dcoptions"]).IsNil() {
		values = strings.Trim(fmt.Sprint(types["dcoptions"]), "[]")
		if reflect.ValueOf(types["expandenv"]).Bool() && types["expandenv"].(bool) {
			values = os.ExpandEnv(values)

		}
		dcoptions = strings.Fields(values)
	}
	if !reflect.ValueOf(types["cmdoptions"]).IsNil() {
		values = strings.Trim(fmt.Sprint(types["cmdoptions"]), "[]")
		if reflect.ValueOf(types["expandenv"]).Bool() && types["expandenv"].(bool) {
			values = os.ExpandEnv(values)

		}
		cmdoptions = strings.Fields(values)
	}
	if string(types["desc"].(string)) != "" {
		desc = fmt.Sprintf("%v", types["desc"])
	}
	if string(types["service"].(string)) != "" {
		service = fmt.Sprintf("%v", types["service"])
	}
	wg.Add(1)
	go exec.CommandDockerComposeExec(types["command"].(string), service, cmdoptions, dcoptions, cmds, desc, _envs, wg)
	wg.Wait()
}

func sshCmd(types map[interface{}]interface{}, _envs []string) {
	wg := new(sync.WaitGroup)
	if !reflect.ValueOf(types["values"]).IsNil() {
		values = strings.Trim(fmt.Sprint(types["values"]), "[]")
		cmds = strings.Fields(values)
	}
	if !reflect.ValueOf(types["options"]).IsNil() {
		values = strings.Trim(fmt.Sprint(types["options"]), "[]")
		options = strings.Fields(values)
	}
	if string(types["desc"].(string)) != "" {
		desc = fmt.Sprintf("%v", types["desc"])
	}
	wg.Add(1)
	go exec.CommandSSH(types["user"].(string), types["port"].(int), types["host"].(string), options, cmds, desc, _envs, wg)
	wg.Wait()
}

func conf(types map[interface{}]interface{}) {
	if reflect.ValueOf(types["confdata"].(string)).String() != "" {
		confdata = types["confdata"].(string)
	}
	if reflect.ValueOf(types["confdest"].(string)).String() != "" {
		confdest = types["confdest"].(string)
	}
	if reflect.ValueOf(types["confperm"].(int)).Int() != 0 {
		confperm = os.FileMode(int(types["confperm"].(int)))
	}
	if confdata != "" && confdest != "" && string(rune(confperm)) != "" {
		functions.WriteFile(confdata, confdest, confperm)
		//readFile(string(confdest))
	}
	color.New(color.FgGreen).Println("# create " + confdest)
}

func Runfromyaml(yamlFile []byte) {

	yaml.Unmarshal(yamlFile, &results)

	for key := range results["env"] {
		if debug {
			functions.PrintColor(color.BgBlue, results["env"])
			functions.PrintColor(color.BgBlue, envs)
			functions.PrintColor(color.BgBlue, key)
		}
		envs = append(envs, results["env"][key].(string))
	}

	for key := range results["cmd"] {
		if !reflect.ValueOf(results["cmd"][key].(map[interface{}]interface{})).IsNil() {
			types = results["cmd"][key].(map[interface{}]interface{})
			if debug {
				functions.PrintColor(color.FgHiCyan, "\n\n%+v\n\n", types)
				functions.PrintColor(color.FgBlue, "Name: %+v\n", types["name"])
				functions.PrintColor(color.FgBlue, "Beschreibung: %+v\n", types["desc"])
				functions.PrintColor(color.FgBlue, "Key: %+v\n", key)
				functions.PrintColor(color.FgBlue, "Command: %+v\n", values)
				functions.PrintColor(color.FgBlue, "Data:\n---\n%+v\n---\n", types["confdata"])
				functions.PrintColor(color.FgBlue, "Destination: %+v\n", types["confdest"])
				functions.PrintColor(color.FgBlue, "Permissions: %+v\n", types["confperm"])
				functions.PrintColor(color.FgHiWhite, "Key: %+v\n", key)
				functions.PrintColor(color.FgGreen, "Name: %+v\n", types["name"])
				functions.PrintColor(color.FgGreen, "Beschreibung: %+v\n", types["desc"])
				functions.PrintColor(color.FgYellow, "Command: %+v\n", values)
				fmt.Printf("\n")
				fmt.Printf("\n")
			}
			if types["type"] == "exec" {
				execCmd(types, envs)
			}
			if types["type"] == "shell" {
				shellCmd(types, envs)
			}
			if types["type"] == "docker" {
				dockerCmd(types, envs)
			}
			if types["type"] == "docker-compose" {
				dockerComposeCmd(types, envs)
			}
			if types["type"] == "ssh" {
				sshCmd(types, envs)
			}
			if types["type"] == "conf" {
				if debug {
					functions.PrintColor(color.FgYellow, "Config: %+v\n", types["confdata"])
					functions.PrintColor(color.FgYellow, "Config: %+v\n", types["confdest"])
					functions.PrintColor(color.FgYellow, "Config: %+v\n", types["confperm"])
					fmt.Printf("\n")
				}
				conf(types)
			}
		}
	}
}
