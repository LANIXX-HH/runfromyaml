package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/fatih/color"

	exec "github.com/lanixx/runfromyaml/pkg/exec"
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
	command  string
	file     string
	help     bool
	debug    bool
	user     string
	port     int
	host     string
)

func main() {

	programm := os.Args

	// parse flags
	flag.StringVar(&file, "file", "commands.yaml", "input config filename")
	flag.BoolVar(&help, "help", false, "Display this help")
	flag.BoolVar(&debug, "debug", false, "Debug Mode")
	flag.Parse()

	red := color.New(color.FgRed).PrintfFunc()
	yellow := color.New(color.FgYellow).PrintfFunc()
	green := color.New(color.FgGreen).PrintfFunc()
	blue := color.New(color.FgBlue).PrintfFunc()
	white := color.New(color.FgHiWhite).PrintfFunc()
	cyan := color.New(color.FgHiCyan).PrintfFunc()

	if debug {
		red("\n", programm)
	}

	yamlFile, err := ioutil.ReadFile(file)
	err = yaml.Unmarshal(yamlFile, &results)
	for key := range results["cmd"] {
		if !reflect.ValueOf(results["cmd"][key].(map[interface{}]interface{})).IsNil() {
			types = results["cmd"][key].(map[interface{}]interface{})
			if debug {
				cyan("\n\n%+v\n\n", types)
				blue("Name: %+v\n", types["name"])
				blue("Beschreibung: %+v\n", types["desc"])
				blue("Key: %+v\n", key)
				blue("Command: %+v\n", values)
				blue("Data:\n---\n%+v\n---\n", types["confdata"])
				blue("Destination: %+v\n", types["confdest"])
				blue("Permissions: %+v\n", types["confperm"])
				fmt.Printf("\n")
			}
			if types["type"] == "exec" {
				if debug {
					white("Key: %+v\n", key)
					green("Name: %+v\n", types["name"])
					green("Beschreibung: %+v\n", types["desc"])
					yellow("Command: %+v\n", values)
					fmt.Printf("\n")
				}
				wg := new(sync.WaitGroup)
				if !reflect.ValueOf(types["values"].(interface{})).IsNil() {
					values = fmt.Sprintf("%v", types["values"].(interface{}))
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
			if types["type"] == "shell" {
				if debug {
					white("Key: %+v\n", key)
					green("Name: %+v\n", types["name"])
					green("Beschreibung: %+v\n", types["desc"])
					yellow("Command: %+v\n", values)
					fmt.Printf("\n")
				}
				wg := new(sync.WaitGroup)
				if !reflect.ValueOf(types["values"].(interface{})).IsNil() {
					values = fmt.Sprintf("%v", types["values"].(interface{}))
					values = strings.TrimPrefix(values, "[")
					values = strings.TrimSuffix(values, "]")
					cmds = strings.Fields(values)
				}
				if string(types["desc"].(string)) != "" {
					desc = fmt.Sprintf("%v", types["desc"])
				}
				wg.Add(1)
				go exec.CommandShell(cmds, desc, wg)
				wg.Wait()
			}
			if types["type"] == "docker" {
				if debug {
					white("Key: %+v\n", key)
					green("Name: %+v\n", types["name"])
					green("Beschreibung: %+v\n", types["desc"])
					yellow("Command: %+v\n", values)
					fmt.Printf("\n")
				}
				wg := new(sync.WaitGroup)
				if !reflect.ValueOf(types["values"].(interface{})).IsNil() {
					values = fmt.Sprintf("%v", types["values"].(interface{}))
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
			if types["type"] == "docker-compose" {
				if debug {
					white("Key: %+v\n", key)
					green("Name: %+v\n", types["name"])
					green("Beschreibung: %+v\n", types["desc"])
					yellow("Command: %+v\n", values)
					fmt.Printf("\n")
				}
				wg := new(sync.WaitGroup)
				if !reflect.ValueOf(types["values"].(interface{})).IsNil() {
					values = fmt.Sprintf("%v", types["values"].(interface{}))
					values = strings.TrimPrefix(values, "[")
					values = strings.TrimSuffix(values, "]")
					cmds = strings.Fields(values)
				}
				if !reflect.ValueOf(types["options"].(interface{})).IsNil() {
					values = fmt.Sprintf("%v", types["options"].(interface{}))
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
			if types["type"] == "ssh" {
				if debug {
					white("Key: %+v\n", key)
					green("Name: %+v\n", types["name"])
					green("Beschreibung: %+v\n", types["desc"])
					yellow("Command: %+v\n", values)
					fmt.Printf("\n")
				}
				wg := new(sync.WaitGroup)
				if !reflect.ValueOf(types["values"].(interface{})).IsNil() {
					values = fmt.Sprintf("%v", types["values"].(interface{}))
					values = strings.TrimPrefix(values, "[")
					values = strings.TrimSuffix(values, "]")
					cmds = strings.Fields(values)
				}
				if !reflect.ValueOf(types["options"].(interface{})).IsNil() {
					values = fmt.Sprintf("%v", types["options"].(interface{}))
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
			if types["type"] == "conf" {
				if debug {
					yellow("Config: %+v\n", types["confdata"])
					yellow("Config: %+v\n", types["confdest"])
					yellow("Config: %+v\n", types["confperm"])
					fmt.Printf("\n")
				}
				if reflect.ValueOf(types["confdata"].(string)).String() != "" {
					confdata = types["confdata"].(string)
				}
				if reflect.ValueOf(types["confdest"].(string)).String() != "" {
					confdest = types["confdest"].(string)
				}
				if reflect.ValueOf(types["confperm"].(int)).Int() != 0 {
					confperm = os.FileMode(int(types["confperm"].(int)))
				}
				if confdata != "" && confdest != "" && string(confperm) != "" {
					functions.WriteFile(confdata, confdest, confperm)
					//readFile(string(confdest))
				}
			}
			fmt.Printf("\n")
		}
	}
	functions.Check(err)
}
