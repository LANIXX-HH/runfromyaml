package main

import (
	"fmt"
	"io/ioutil"
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

func exeCmd(cmd string, wg *sync.WaitGroup) {
	fmt.Println(cmd)
	parts := strings.Fields(cmd)
	out, err := exec.Command(parts[0], parts[1]).Output()
	if err != nil {
		fmt.Println("error occured")
		fmt.Printf("%s", err)
	}
	fmt.Printf("%s", out)
	wg.Done()
}

func writeFile(file string, path string, perm os.FileMode) {
	bytefile := []byte(file)
	err := ioutil.WriteFile(path, bytefile, perm)
	check(err)
}

func main() {
	var config Configs
	wg := new(sync.WaitGroup)

	yamlFile, err := ioutil.ReadFile("commands.yaml")
	check(err)

	err = yaml.UnmarshalStrict(yamlFile, &config)
	//err = yaml.Unmarshal(yamlFile, &config)
	check(err)

	fmt.Printf("\n%+v\n\n", config)
	for i := 0; i < len(config.Cfgs); i++ {
		if config.Cfgs[i].Type == "shell" {

			fmt.Printf("Name: %+v\n", config.Cfgs[i].Name)
			fmt.Printf("Beschreibung: %+v\n", config.Cfgs[i].Desc)
			fmt.Printf("Command: %+v\n", config.Cfgs[i].Values[0])
			fmt.Printf("Config: %+v\n", config.Cfgs[i].Conf)
			fmt.Printf("\n")

			runcommand := config.Cfgs[i].Values[0]
			for j := 1; j < len(config.Cfgs); j++ {
				runcommand = runcommand + " " + config.Cfgs[i].Values[j]
			}
			fmt.Printf(runcommand + "\n")

			wg.Add(1)
			writeFile(config.Cfgs[i].Conf, config.Cfgs[i].Confdest, config.Cfgs[i].Confperm)
			go exeCmd(runcommand, wg)
			wg.Wait()

			fmt.Printf("\n")
		}
	}
}
