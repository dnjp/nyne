package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func Load(path string) (*Config, error) {
	conf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %s", err)
	}
	var cfg = new(Config)
	if err := json.Unmarshal(conf, cfg); err != nil {
		return nil, fmt.Errorf("error decoding config file: %s", err)
	}
	return cfg, nil
}

type Config struct {
	Spec []Spec `json:"spec"`
	Menu []string `json:"menu"`
}

type Spec struct {
	Ext []string `json:"ext"`
	Cmd []Command `json:"cmd"`
	Fmt Format `json:"fmt"`
}

type Format struct {
	Indent int `json:"indent"`
	Expand bool `json:"expand"`
}

type Command struct {
	Exec string `json:"exec"`
	Args []string `json:"args"`
}