package main

import (
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
	var conf string
	var confdest string
	var perm os.FileMode

	argsWithoutProg := os.Args[1:]
	argsWithoutProgAsString := strings.Join(argsWithoutProg, ",")

	yamlFile, err := ioutil.ReadFile("commands.yaml")
	err = yaml.Unmarshal(yamlFile, &results)
	for key := range results["cmd"] {
		if !reflect.ValueOf(results["cmd"][key].(map[interface{}]interface{})).IsNil() {
			types = results["cmd"][key].(map[interface{}]interface{})
			if strings.Contains(argsWithoutProgAsString, "--debug") {
				fmt.Printf("\n%+v\n\n", types)
				fmt.Printf("Config: %+v\n", key)
				fmt.Printf("Name: %+v\n", types["name"])
				fmt.Printf("Beschreibung: %+v\n", types["desc"])
				fmt.Printf("Command: %+v\n", values)
				fmt.Printf("Config: %+v\n", types["conf"])
				fmt.Printf("Config: %+v\n", types["confdest"])
				fmt.Printf("Config: %+v\n", types["confperm"])
				fmt.Printf("\n")
			}
			if conf != "" && confdest != "" && string(perm) != "" {
				writeFile(conf, confdest, perm)
			}
			if types["type"] == "shell" {
				if strings.Contains(argsWithoutProgAsString, "--debug") {
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
			if types["conf"] == "conf" {
				if strings.Contains(argsWithoutProgAsString, "--debug") {
					fmt.Printf("\n%+v\n\n", types)
					fmt.Printf("Config: %+v\n", types["conf"])
					fmt.Printf("Config: %+v\n", types["confdest"])
					fmt.Printf("Config: %+v\n", types["confperm"])
					fmt.Printf("\n")
				}
				if reflect.ValueOf(types["conf"].(string)).String() != "" {
					conf = types["confdata"].(string)
				}
				if reflect.ValueOf(types["confdest"].(string)).String() != "" {
					confdest = types["confdest"].(string)
				}
				if reflect.ValueOf(types["confperm"].(int)).Int() != 0 {
					perm = os.FileMode(int(types["confperm"].(int)))
				}
			}
			fmt.Printf("\n")
		}
	}
	check(err)
}
