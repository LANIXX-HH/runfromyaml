package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dchest/uniuri"
	"github.com/fatih/color"
	"github.com/lanixx/runfromyaml/pkg/cli"
	"github.com/lanixx/runfromyaml/pkg/functions"
	"github.com/lanixx/runfromyaml/pkg/restapi"
)

func init() {
	functions.Config()
}

func main() {
	var (
		file     string
		debug    bool
		rest     bool
		port     int
		host     string
		user     string
		restauth bool
	)

	programm := os.Args

	// parse flags
	flag.StringVar(&file, "f", "commands.yaml", "file - file with all defined commands, descriptions and configuration blocks in yaml fromat")
	flag.BoolVar(&debug, "d", false, "debug - activate debug mode to print more informations")
	flag.BoolVar(&rest, "r", false, "restapi - start this instance in background mode in rest api mode")
	flag.IntVar(&port, "p", 8080, "port - set http port for rest api mode (default http port is 8080)")
	flag.StringVar(&host, "h", "localhost", "host - set host for rest api mode (default host is localhost)")
	flag.StringVar(&user, "u", "rest", "user - set username for rest api authentication (default username is rest) ")
	flag.BoolVar(&restauth, "n", false, "no-auth - disable rest auth")

	flag.Parse()

	if debug {
		functions.PrintColor(color.FgRed, "debug", "stdout", "\n", programm)
	}

	yamlFile, err := os.ReadFile(file)

	if rest {
		fmt.Println("start command in rest api mode on", host, "host", port, "port")

		if !restauth {
			restapi.RestAuth = true
			restapi.TempPass = uniuri.New()
			restapi.TempUser = user
			fmt.Println("temporary password for rest api connection with user", restapi.TempUser, "is", restapi.TempPass)
		} else {
			restapi.RestAuth = false
		}
		restapi.RestApi(port, host)
	} else {
		if err != nil {
			fmt.Println("\ninput file not found. please use -f option to set input file or create default commands.yaml file \n\n valid options are:\n")
			flag.PrintDefaults()
			fmt.Println("\n")
		}
		cli.Runfromyaml(yamlFile, debug)
	}
}
