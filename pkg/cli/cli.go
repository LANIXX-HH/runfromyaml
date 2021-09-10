package cli

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/fatih/color"

	"github.com/lanixx/runfromyaml/pkg/exec"
	functions "github.com/lanixx/runfromyaml/pkg/functions"

	"gopkg.in/yaml.v2"
)

var (
	results  map[interface{}][]interface{}
	types    map[interface{}]interface{}
	values   string
	cmds     []string
	options  []string
	desc     string
	confdata string
	confdest string
	confperm os.FileMode
	file     string
	help     bool
	debug    bool
)

func printColor(ctype color.Attribute, cstring ...interface{}) {
	mystring := color.New(ctype)
	mystring.Println(cstring)
}

func execCmd(types map[interface{}]interface{}) {
	wg := new(sync.WaitGroup)
	if !reflect.ValueOf(types["values"]).IsNil() {
		values = fmt.Sprintf("%v", types["values"])
		values = strings.TrimPrefix(values, "[")
		values = strings.TrimSuffix(values, "]")
		cmds = strings.Fields(values)
	}
	if string(types["desc"].(string)) != "" {
		desc = fmt.Sprintf("%v", types["desc"])
	}
	wg.Add(1)
	go exec.Command(cmds, desc, wg)
	wg.Wait()
}

func shellCmd(types map[interface{}]interface{}) {
	if !reflect.ValueOf(types["values"]).IsNil() {
		values = fmt.Sprintf("%v", types["values"])
		values = strings.TrimPrefix(values, "[")
		values = strings.TrimSuffix(values, "]")
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
		go exec.CommandShell(shcmd, desc, wg, ind)
		wg.Wait()
	}
}

func dockerCmd(types map[interface{}]interface{}) {
	wg := new(sync.WaitGroup)
	if !reflect.ValueOf(types["values"]).IsNil() {
		values = fmt.Sprintf("%v", types["values"])
		values = strings.TrimPrefix(values, "[")
		values = strings.TrimSuffix(values, "]")
		cmds = strings.Fields(values)
	}
	if string(types["desc"].(string)) != "" {
		desc = fmt.Sprintf("%v", types["desc"])
	}
	wg.Add(1)
	go exec.CommandDockerRun(types["command"].(string), types["container"].(string), cmds, desc, wg)
	wg.Wait()
}

func dockerComposeCmd(types map[interface{}]interface{}) {
	wg := new(sync.WaitGroup)
	if !reflect.ValueOf(types["values"]).IsNil() {
		values = fmt.Sprintf("%v", types["values"])
		values = strings.TrimPrefix(values, "[")
		values = strings.TrimSuffix(values, "]")
		cmds = strings.Fields(values)
	}
	if !reflect.ValueOf(types["options"]).IsNil() {
		values = fmt.Sprintf("%v", types["options"])
		values = strings.TrimPrefix(values, "[")
		values = strings.TrimSuffix(values, "]")
		options = strings.Fields(values)
	}
	if string(types["desc"].(string)) != "" {
		desc = fmt.Sprintf("%v", types["desc"])
	}
	wg.Add(1)
	go exec.CommandDockerComposeExec(options, cmds, desc, wg)
	wg.Wait()
}

func sshCmd(types map[interface{}]interface{}) {
	wg := new(sync.WaitGroup)
	if !reflect.ValueOf(types["values"]).IsNil() {
		values = fmt.Sprintf("%v", types["values"])
		values = strings.TrimPrefix(values, "[")
		values = strings.TrimSuffix(values, "]")
		cmds = strings.Fields(values)
	}
	if !reflect.ValueOf(types["options"]).IsNil() {
		values = fmt.Sprintf("%v", types["options"])
		values = strings.TrimPrefix(values, "[")
		values = strings.TrimSuffix(values, "]")
		options = strings.Fields(values)
	}
	if string(types["desc"].(string)) != "" {
		desc = fmt.Sprintf("%v", types["desc"])
	}
	wg.Add(1)
	go exec.CommandSSH(types["user"].(string), types["port"].(int), types["host"].(string), options, cmds, desc, wg)
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
	fmt.Printf("\n")
}

func Runfromyaml() {

	programm := os.Args

	// parse flags
	flag.StringVar(&file, "file", "commands.yaml", "input config filename")
	flag.BoolVar(&help, "help", false, "Display this help")
	flag.BoolVar(&debug, "debug", false, "Debug Mode")
	flag.Parse()

	if debug {
		printColor(color.FgRed, "\n", programm)
	}

	yamlFile, err := ioutil.ReadFile(file)
	yaml.Unmarshal(yamlFile, &results)
	for key := range results["cmd"] {
		if !reflect.ValueOf(results["cmd"][key].(map[interface{}]interface{})).IsNil() {
			types = results["cmd"][key].(map[interface{}]interface{})
			if debug {
				printColor(color.FgHiCyan, "\n\n%+v\n\n", types)
				printColor(color.FgBlue, "Name: %+v\n", types["name"])
				printColor(color.FgBlue, "Beschreibung: %+v\n", types["desc"])
				printColor(color.FgBlue, "Key: %+v\n", key)
				printColor(color.FgBlue, "Command: %+v\n", values)
				printColor(color.FgBlue, "Data:\n---\n%+v\n---\n", types["confdata"])
				printColor(color.FgBlue, "Destination: %+v\n", types["confdest"])
				printColor(color.FgBlue, "Permissions: %+v\n", types["confperm"])
				printColor(color.FgHiWhite, "Key: %+v\n", key)
				printColor(color.FgGreen, "Name: %+v\n", types["name"])
				printColor(color.FgGreen, "Beschreibung: %+v\n", types["desc"])
				printColor(color.FgYellow, "Command: %+v\n", values)
				fmt.Printf("\n")
				fmt.Printf("\n")
			}
			if types["type"] == "exec" {
				execCmd(types)
			}
			if types["type"] == "shell" {
				shellCmd(types)
			}
			if types["type"] == "docker" {
				dockerCmd(types)
			}
			if types["type"] == "docker-compose" {
				dockerComposeCmd(types)
			}
			if types["type"] == "ssh" {
				sshCmd(types)
			}
			if types["type"] == "conf" {
				if debug {
					printColor(color.FgYellow, "Config: %+v\n", types["confdata"])
					printColor(color.FgYellow, "Config: %+v\n", types["confdest"])
					printColor(color.FgYellow, "Config: %+v\n", types["confperm"])
					fmt.Printf("\n")
				}
				conf(types)
			}
		}
	}
	functions.Check(err)
}
