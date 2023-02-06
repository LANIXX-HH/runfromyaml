package main

import (
	"flag"
	"fmt"
	"os"
	"reflect"

	"github.com/dchest/uniuri"
	"github.com/fatih/color"
	"github.com/lanixx/runfromyaml/pkg/cli"
	"github.com/lanixx/runfromyaml/pkg/functions"
	"github.com/lanixx/runfromyaml/pkg/restapi"
	"gopkg.in/yaml.v2"
)

func init() {
	functions.Config()
}

func main() {
	var (
		ydoc map[interface{}][]interface{}
	)

	programm := os.Args
	flags := make(map[string]interface{})

	// parse flags
	flags["debug"] = flag.Bool("debug", false, "debug - activate debug mode to print more informations")
	flags["rest"] = flag.Bool("rest", false, "restapi - start this instance in background mode in rest api mode")
	flags["no-auth"] = flag.Bool("no-auth", false, "no-auth - disable rest auth")
	flags["restout"] = flag.Bool("restout", false, "rest output - activate output to http response")
	flags["no-file"] = flag.Bool("no-file", false, "no-file - file option should be disabled")

	flags["file"] = flag.String("file", "commands.yaml", "file - file with all defined commands, descriptions and configuration blocks in yaml fromat")
	flags["host"] = flag.String("host", "localhost", "host - set host for rest api mode (default host is localhost)")
	flags["user"] = flag.String("user", "rest", "user - set username for rest api authentication (default username is rest) ")

	flags["port"] = flag.Int("port", 8080, "port - set http port for rest api mode (default http port is 8080)")

	flag.Parse()
	yamlFile, err := os.ReadFile(*flags["file"].(*string))
	if err != nil {
		fmt.Println("\n file option was set, but it was not possible to read input yaml file.")
	}

	yaml.Unmarshal(yamlFile, &ydoc)

	for key := range ydoc["options"] {
		options := ydoc["options"][key].(map[interface{}]interface{})
		if options["key"] == "file" || options["key"] == "host" || options["key"] == "user" {
			*flags[options["key"].(string)].(*string) = options["value"].(string)
		}
		if options["key"] == "debug" || options["key"] == "rest" || options["key"] == "no-auth" || options["key"] == "restout" || options["key"] == "no-file" {
			*flags[options["key"].(string)].(*bool) = options["value"].(bool)
		}
		if options["key"] == "port" {
			*flags[options["key"].(string)].(*int) = options["value"].(int)
		}

	}

	if *flags["debug"].(*bool) {
		functions.PrintColor(color.FgRed, "debug", "\n", programm)
	}

	_, filerr := os.Stat("commands.yaml")
	if reflect.ValueOf(*flags["file"].(*string)).IsValid() && filerr == nil && !*flags["no-file"].(*bool) {
		cli.Runfromyaml(yamlFile, *flags["debug"].(*bool))
	}

	if *flags["rest"].(*bool) {
		fmt.Println("start command in rest api mode on", *flags["host"].(*string), "host", *flags["port"].(*int), "port")

		if *flags["restout"].(*bool) {
			restapi.RestOut = *flags["restout"].(*bool)
			fmt.Println("output should be redirected to rest http response")
		}

		if *flags["no-auth"].(*bool) {
			restapi.RestAuth = false
		} else {
			restapi.RestAuth = true
			restapi.TempPass = uniuri.New()
			restapi.TempUser = *flags["user"].(*string)
			fmt.Println("temporary password for rest api connection with user", restapi.TempUser, "is", restapi.TempPass)
		}
		restapi.RestApi(*flags["port"].(*int), *flags["host"].(*string))
	}
}
