package config

import (
	"github.com/BurntSushi/toml"
)

func Load(path string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

type Config struct {
	Format map[string]Spec   
	Tag Tag
}

type Tag struct {
	Menu []string 
}

type Spec struct {
	Indent int 
	Tabexpand bool
	Extensions []string  
	Commands []Command
}

type Command struct {
	Exec string 
	Args []string
}
