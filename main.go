package main

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/fatih/color"
	cli "github.com/lanixx/runfromyaml/pkg/cli"
	functions "github.com/lanixx/runfromyaml/pkg/functions"
)

var (
	file  string
	help  bool
	debug bool
)

func init() {
	functions.Config()
}

func main() {
	programm := os.Args

	// parse flags
	flag.StringVar(&file, "file", "commands.yaml", "input config filename")
	flag.BoolVar(&help, "help", false, "Display this help")
	flag.BoolVar(&debug, "debug", false, "Debug Mode")
	flag.Parse()

	if debug {
		functions.PrintColor(color.FgRed, "debug", "stdout", "\n", programm)
	}

	yamlFile, err := ioutil.ReadFile(file)
	cli.Runfromyaml(yamlFile)
	functions.Check(err)

}
