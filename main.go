package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"
)

// Config 2
type Config struct {
	Type     string
	Name     string
	Desc     string
	Values   []string
	Conf     string
	Confdest string
	Confperm os.FileMode
}

// Configs cool
type Configs struct {
	Cfgs []Config `yaml:"cmd"`
}

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
	var config Configs
	wg := new(sync.WaitGroup)
	argsWithoutProg := os.Args[1:]
	argsWithoutProgAsString := strings.Join(argsWithoutProg, ",")

	yamlFile, err := ioutil.ReadFile("commands.yaml")
	check(err)

	err = yaml.UnmarshalStrict(yamlFile, &config)
	//err = yaml.Unmarshal(yamlFile, &config)
	check(err)

	for i := 0; i < len(config.Cfgs); i++ {
		if config.Cfgs[i].Type == "shell" {
			//if config.Cfgs[i].Debug {
			if strings.Contains(argsWithoutProgAsString, "--debug") {
				fmt.Printf("\n%+v\n\n", config)
				fmt.Printf("Name: %+v\n", config.Cfgs[i].Name)
				fmt.Printf("Beschreibung: %+v\n", config.Cfgs[i].Desc)
				fmt.Printf("Command: %+v\n", config.Cfgs[i].Values[0])
				fmt.Printf("Config: %+v\n", config.Cfgs[i].Conf)
				fmt.Printf("\n")
			}

			if len(config.Cfgs[i].Conf) > 0 && len(config.Cfgs[i].Confdest) > 0 {
				if config.Cfgs[i].Confperm == os.FileMode(0000) {
					config.Cfgs[i].Confperm = 0644
				}
				writeFile(config.Cfgs[i].Conf, config.Cfgs[i].Confdest, config.Cfgs[i].Confperm)
			}

			wg.Add(1)
			go exeCommand(config.Cfgs[i].Values, config.Cfgs[i].Desc, wg)
			wg.Wait()
			fmt.Printf("\n")
		}
	}
}
