package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func exeCommand(cmd []string, desc string, wg *sync.WaitGroup) {
	fmt.Println("==> command: " + desc)
	command := exec.Command(cmd[0], cmd[1:]...)
	fmt.Println(command)
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		log.Fatalf("Start: %v", err)
	}
	wg.Done()
}

func writeFile(file string, path string, perm os.FileMode) {
	bytefile := []byte(file)
	err := ioutil.WriteFile(path, bytefile, perm)
	check(err)
}

func readFile(file string) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("File contents: %s", content)
}

func remove(slice []string, i int) []string {
	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
}

func main() {
	var results map[interface{}][]interface{}
	var types map[interface{}]interface{}
	var values string
	var cmds []string
	var desc string
	var confdata string
	var confdest string
	var confperm os.FileMode
	var file string
	var help bool
	var debug bool

	programm := os.Args

	// parse flags
	flag.StringVar(&file, "file", "commands.yaml", "input config filename")
	flag.BoolVar(&help, "help", false, "Display this help")
	flag.BoolVar(&debug, "debug", false, "Debug Mode")
	flag.Parse()

	yamlFile, err := ioutil.ReadFile(file)
	err = yaml.Unmarshal(yamlFile, &results)
	for key := range results["cmd"] {
		if !reflect.ValueOf(results["cmd"][key].(map[interface{}]interface{})).IsNil() {
			types = results["cmd"][key].(map[interface{}]interface{})
			if debug {
				fmt.Printf("\n", programm)
				fmt.Printf("\n\n%+v\n\n", types)
				fmt.Printf("Config: %+v\n", key)
				fmt.Printf("Name: %+v\n", types["name"])
				fmt.Printf("Beschreibung: %+v\n", types["desc"])
				fmt.Printf("Command: %+v\n", values)
				fmt.Printf("Config: %+v\n", types["confdata"])
				fmt.Printf("Config: %+v\n", types["confdest"])
				fmt.Printf("Config: %+v\n", types["confperm"])
				fmt.Printf("\n")
			}
			if types["type"] == "shell" {
				if debug {
					fmt.Printf("\n%+v\n\n", types)
					fmt.Printf("Config: %+v\n", key)
					fmt.Printf("Name: %+v\n", types["name"])
					fmt.Printf("Beschreibung: %+v\n", types["desc"])
					fmt.Printf("Command: %+v\n", values)
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
				go exeCommand(cmds, desc, wg)
				wg.Wait()
			}
			if types["type"] == "conf" {
				if debug {
					fmt.Printf("\n%+v\n\n", types)
					fmt.Printf("Config: %+v\n", types["confdata"])
					fmt.Printf("Config: %+v\n", types["confdest"])
					fmt.Printf("Config: %+v\n", types["confperm"])
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
					writeFile(confdata, confdest, confperm)
					//readFile(string(confdest))
				}
			}
			fmt.Printf("\n")
		}
	}
	check(err)
}
