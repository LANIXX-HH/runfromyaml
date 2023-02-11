package config

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

// YAMLOptions represents the options section in the YAML file
type YAMLOptions struct {
	Options []struct {
		Key   string      `yaml:"key"`
		Value interface{} `yaml:"value"`
	} `yaml:"options"`
}

// LoadFromYAML loads configuration from YAML data
func (c *Config) LoadFromYAML(data []byte) error {
	var yamlOpts YAMLOptions
	if err := yaml.Unmarshal(data, &yamlOpts); err != nil {
		return fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	for _, opt := range yamlOpts.Options {
		switch opt.Key {
		case "debug", "rest", "no-auth", "restout", "no-file", "ai":
			if val, ok := opt.Value.(bool); ok {
				switch opt.Key {
				case "debug":
					c.Debug = val
				case "rest":
					c.Rest = val
				case "no-auth":
					c.NoAuth = val
				case "restout":
					c.RestOut = val
				case "no-file":
					c.NoFile = val
				case "ai":
					c.AI = val
				}
			}
		case "file", "host", "user", "ai-key", "ai-model", "ai-cmdtype", "shell-type":
			if val, ok := opt.Value.(string); ok {
				switch opt.Key {
				case "file":
					c.File = val
				case "host":
					c.Host = val
				case "user":
					c.User = val
				case "ai-key":
					c.AIKey = val
				case "ai-model":
					c.AIModel = val
				case "ai-cmdtype":
					c.AICmdType = val
				case "shell-type":
					c.ShellType = val
				}
			}
		case "port":
			if val, ok := opt.Value.(int); ok {
				c.Port = val
			}
		}
	}

	return nil
}
