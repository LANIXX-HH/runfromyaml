package config

import (
	"flag"
)

// Config holds all configuration for the application
type Config struct {
	Debug     bool
	Rest      bool
	NoAuth    bool
	RestOut   bool
	NoFile    bool
	AI        bool
	Shell     bool
	File      string
	Host      string
	User      string
	AIInput   string
	AIKey     string
	AIModel   string
	AICmdType string
	ShellType string
	Port      int
}

// New creates a new Config instance with default values
func New() *Config {
	return &Config{
		Debug:     false,
		Rest:      false,
		NoAuth:    false,
		RestOut:   false,
		NoFile:    false,
		AI:        false,
		Shell:     false,
		File:      "commands.yaml",
		Host:      "localhost",
		User:      "rest",
		AIInput:   "",
		AIKey:     "",
		AIModel:   "text-davinci-003",
		AICmdType: "shell",
		ShellType: "bash",
		Port:      8080,
	}
}

// ParseFlags parses command line flags into the config
func (c *Config) ParseFlags() {
	flag.BoolVar(&c.Debug, "debug", c.Debug, "debug - activate debug mode to print more informations")
	flag.BoolVar(&c.Rest, "rest", c.Rest, "restapi - start this instance in background mode in rest api mode")
	flag.BoolVar(&c.NoAuth, "no-auth", c.NoAuth, "no-auth - disable rest auth")
	flag.BoolVar(&c.RestOut, "restout", c.RestOut, "rest output - activate output to http response")
	flag.BoolVar(&c.NoFile, "no-file", c.NoFile, "no-file - file option should be disabled")
	flag.BoolVar(&c.AI, "ai", c.AI, "ai - interact with OpenAI")
	flag.BoolVar(&c.Shell, "shell", c.Shell, "shell - interactive shell")

	flag.StringVar(&c.File, "file", c.File, "file - file with all defined commands, descriptions and configuration blocks in yaml fromat")
	flag.StringVar(&c.Host, "host", c.Host, "host - set host for rest api mode (default host is localhost)")
	flag.StringVar(&c.User, "user", c.User, "user - set username for rest api authentication (default username is rest)")
	flag.StringVar(&c.AIInput, "ai-in", c.AIInput, "ai - interact with OpenAI")
	flag.StringVar(&c.AIKey, "ai-key", c.AIKey, "ai - OpenAI API Key")
	flag.StringVar(&c.AIModel, "ai-model", c.AIModel, "ai-model - OpenAI Model for answer generation")
	flag.StringVar(&c.AICmdType, "ai-cmdtype", c.AICmdType, "ai-cmdtype - For which type of code should be examples generated")
	flag.StringVar(&c.ShellType, "shell-type", c.ShellType, "shell-type - which shell type should be used for recording all the commands to generate yaml structure")

	flag.IntVar(&c.Port, "port", c.Port, "port - set http port for rest api mode (default http port is 8080)")

	flag.Parse()
}
