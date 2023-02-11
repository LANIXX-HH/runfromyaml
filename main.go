package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dchest/uniuri"
	"github.com/fatih/color"
	"github.com/lanixx/runfromyaml/pkg/cli"
	"github.com/lanixx/runfromyaml/pkg/config"
	"github.com/lanixx/runfromyaml/pkg/functions"
	"github.com/lanixx/runfromyaml/pkg/openai"
	"github.com/lanixx/runfromyaml/pkg/restapi"
	"gopkg.in/yaml.v2"
)

func init() {
	functions.Config()
}

func main() {
	cfg := config.New()
	cfg.ParseFlags()

	if cfg.Debug {
		functions.PrintColor(color.FgRed, "debug", "\n", os.Args)
	}

	// Load YAML configuration if file exists
	if !cfg.NoFile {
		ydata, err := os.ReadFile(cfg.File)
		if err == nil {
			if err := cfg.LoadFromYAML(ydata); err != nil {
				fmt.Printf("Warning: Failed to load YAML configuration: %v\n", err)
			}
		}
	}

	// Handle AI mode
	if cfg.AI {
		handleAIMode(cfg)
	}

	// Handle file-based execution
	if !cfg.NoFile {
		handleFileExecution(cfg)
	}

	// Handle REST API mode
	if cfg.Rest {
		handleRestMode(cfg)
	}

	// Handle interactive shell mode
	if cfg.Shell {
		handleShellMode(cfg)
	}
}

func handleAIMode(cfg *config.Config) {
	if len(cfg.AIKey) > 0 {
		openai.Key = cfg.AIKey
		openai.IsAiEnabled = true
	} else {
		openai.IsAiEnabled = false
	}

	openai.Model = cfg.AIModel
	openai.ShellType = cfg.AICmdType

	if openai.IsAiEnabled {
		if len(cfg.AIInput) > 0 {
			for {
				response, err := openai.OpenAI(openai.Key, openai.Model, cfg.AIInput, openai.ShellType)
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

func handleFileExecution(cfg *config.Config) {
	ydata, err := os.ReadFile(cfg.File)
	if err != nil {
		fmt.Printf("\nfile option was set, but it was not possible to read this file:\n\t%s\n", cfg.File)
		return
	}

	var ydoc map[interface{}][]interface{}
	if err := yaml.Unmarshal(ydata, &ydoc); err != nil {
		fmt.Printf("It was not possible to read yaml structure from this file %s with following error message:\n%s\n", cfg.File, err)
		return
	}

	cli.Runfromyaml(ydata, cfg.Debug)
}

func handleRestMode(cfg *config.Config) {
	fmt.Printf("start command in rest api mode on %s host %d port\n", cfg.Host, cfg.Port)

	if cfg.RestOut {
		restapi.RestOut = cfg.RestOut
		fmt.Println("output should be redirected to rest http response")
	}

	if cfg.NoAuth {
		restapi.RestAuth = false
	} else {
		restapi.RestAuth = true
		restapi.TempPass = uniuri.New()
		restapi.TempUser = cfg.User
		fmt.Printf("temporary password for rest api connection with user %s is %s\n", restapi.TempUser, restapi.TempPass)
	}

	restapi.RestApi(cfg.Port, cfg.Host)
}

func handleShellMode(cfg *config.Config) {
	fmt.Println("your input commands will be written to create a YAML structure")
	fmt.Println("enter 'exit' + '\\n' to stop interactive recording")

	// Create a new environment instance
	env := cli.NewEnvironment()

	// Add current environment variables
	for _, envVar := range os.Environ() {
		parts := strings.SplitN(envVar, "=", 2)
		env.Set(parts[0], parts[1])
	}

	commands := cli.InteractiveShell(cfg.ShellType)
	tempmap := functions.PrintShellCommandsAsYaml(commands, env.GetVariables())
	tempyaml, err := yaml.Marshal(tempmap)
	if err != nil {
		fmt.Println("error by marshaling temporary map to yaml")
		return
	}
	fmt.Println(string(tempyaml))
}
