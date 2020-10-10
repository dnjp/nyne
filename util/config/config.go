package config

import (
	"github.com/BurntSushi/toml"
)

// Load parses the TOML configuration file at the specified path
func Load(path string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// Config represents the entire configuration file
type Config struct {
	Format map[string]Spec
	Tag    Tag
}

// Tag configures options for the Acme tag
type Tag struct {
	Menu []string
}

// Spec contains the formatting specified by the config file
type Spec struct {
	Indent     int
	Tabexpand  bool
	Extensions []string
	Commands   []Command
}

// Command contains options for executing a given command against an
// acme window
type Command struct {
	Exec string
	Args []string
}
