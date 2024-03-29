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
	"github.com/lanixx/runfromyaml/pkg/openai"
	"github.com/lanixx/runfromyaml/pkg/restapi"
	"gopkg.in/yaml.v2"
)

func init() {
	functions.Config()
}

func main() {
	var (
		ydoc  map[interface{}][]interface{}
		err   error
		ydata []byte
		yfile string
	)

	programm := os.Args
	flags := make(map[string]interface{})

	// parse flags
	flags["debug"] = flag.Bool("debug", false, "debug - activate debug mode to print more informations")
	flags["rest"] = flag.Bool("rest", false, "restapi - start this instance in background mode in rest api mode")
	flags["no-auth"] = flag.Bool("no-auth", false, "no-auth - disable rest auth")
	flags["restout"] = flag.Bool("restout", false, "rest output - activate output to http response")
	flags["no-file"] = flag.Bool("no-file", false, "no-file - file option should be disabled")
	flags["ai"] = flag.Bool("ai", false, "ai - interact with OpenAI")
	flags["shell"] = flag.Bool("shell", false, "shell - interactive shell ")

	flags["file"] = flag.String("file", "commands.yaml", "file - file with all defined commands, descriptions and configuration blocks in yaml fromat")
	flags["host"] = flag.String("host", "localhost", "host - set host for rest api mode (default host is localhost)")
	flags["user"] = flag.String("user", "rest", "user - set username for rest api authentication (default username is rest) ")
	flags["ai-in"] = flag.String("ai-in", "", "ai - interact with OpenAI")
	flags["ai-key"] = flag.String("ai-key", "", "ai - OpenAI API Key")
	flags["ai-model"] = flag.String("ai-model", "text-davinci-003", "ai-model - OpenAI Model for answer generation")
	flags["ai-cmdtype"] = flag.String("ai-cmdtype", "shell", "ai-cmdtype - For which type of code should be examples generated")
	flags["shell-type"] = flag.String("shell-type", "bash", "shell-type - which shell type should be used for recording all the commands to generate yaml structure")

	flags["port"] = flag.Int("port", 8080, "port - set http port for rest api mode (default http port is 8080)")

	flag.Parse()
	yfile = *flags["file"].(*string)
	ydata, err = os.ReadFile(yfile)
	if err != nil {
		fmt.Println("\nfile option was set, but it was not possible to read this file:\n\t", yfile)
	}

	err = yaml.Unmarshal(ydata, &ydoc)
	if err != nil {
		fmt.Println("It was not possibile to read yaml structure from this file", yfile, "with following error message:\n", err)
	}

	for key := range ydoc["options"] {
		options := ydoc["options"][key].(map[interface{}]interface{})
		if options["key"] == "file" || options["key"] == "host" || options["key"] == "user" || options["key"] == "ai-key" || options["key"] == "ai-model" || options["key"] == "ai-cmdtype" {
			*flags[options["key"].(string)].(*string) = options["value"].(string)
		}
		if options["key"] == "debug" || options["key"] == "rest" || options["key"] == "no-auth" || options["key"] == "restout" || options["key"] == "no-file" || options["key"] == "ai" {
			*flags[options["key"].(string)].(*bool) = options["value"].(bool)
		}
		if options["key"] == "port" {
			*flags[options["key"].(string)].(*int) = options["value"].(int)
		}

	}

	if *flags["debug"].(*bool) {
		functions.PrintColor(color.FgRed, "debug", "\n", programm)
	}
	if *flags["ai"].(*bool) {

		var response map[string][]interface{}

		if len(*flags["ai-key"].(*string)) > 0 {
			openai.Key = *flags["ai-key"].(*string)
			openai.IsAiEnabled = true
		} else {
			openai.IsAiEnabled = false
		}

		openai.Model = *flags["ai-model"].(*string)
		openai.ShellType = *flags["ai-cmdtype"].(*string)

		if openai.IsAiEnabled {
			if len(*flags["ai-in"].(*string)) > 0 {
				for {
					response, err = openai.OpenAI(openai.Key, openai.Model, *flags["ai-in"].(*string), openai.ShellType)
					if err == nil {
						fmt.Println(openai.PrintAiResponse(response))
						break
					}
				}
			}
		} else {
			fmt.Println("OpenAI is not enabled. Probably OpenAI-Key is empty.")
		}
	}

	_, filerr := os.Stat(yfile)
	if filerr != nil {
		fmt.Println("\t", filerr)
	}
	if reflect.ValueOf(yfile).IsValid() && filerr == nil && !*flags["no-file"].(*bool) {
		cli.Runfromyaml(ydata, *flags["debug"].(*bool))
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
	if *flags["shell"].(*bool) {
		fmt.Println("your input commands will be written to create a YAML structure")
		fmt.Println("enter 'exit' + '\\n' to stop interactive recording")
		commands := cli.InteractiveShell(*flags["shell-type"].(*string))
		tempmap := functions.PrintShellCommandsAsYaml(commands, cli.EnvironmentVariables)
		tempyaml, err := yaml.Marshal(tempmap)
		if err != nil {
			fmt.Println("error by marshaling temporary map to yaml")
		}
		fmt.Println(string(tempyaml))
	}
}
