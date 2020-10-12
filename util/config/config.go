package config

import (
	"text/template"
	"go/format"
	"bytes"
	
	"github.com/BurntSushi/toml"
)

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
	CommentStyle string
	Indent     int
	Tabexpand  bool
	Extensions []string
	Commands   []Command
}

// Command contains options for executing a given command against an
// acme window
type Command struct {
	Exec           string
	Args           []string
	PrintsToStdout bool
}

type TmplConfig struct {
	SpecProps []TmplSpecType
	Specs []TmplSpec
}

type TmplSpecType struct {
	Name string
	DataType string
}

type TmplSpec struct {
	Ext string
	CommentStyle string
	Indent     int
	Tabexpand  bool
}

// Load parses the TOML configuration file at the specified path
func Load(path string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func Compile(cfgPath, writePath, tmplPath string) error {
	cfg, err := Load(cfgPath)
	if err != nil {
		return err
	}
	
	specs := []TmplSpec{}
	for _, spec := range cfg.Format {
		for _, ext := range spec.Extensions {
			ts := TmplSpec{
				Ext: ext,
				CommentStyle: spec.CommentStyle,
				Indent: spec,Indent,
				Tabexpand: spec.Tabexpand,
			}
			specs = append(specs, ts)
		}
	}
	
	tc := TmplConfig{
		SpecProps: []TmplSpecType{
			{
				Name: "CommentStyle",
				DataType: "string",
			},
			{
				Name: "Indent",
				DataType: "int",
			},
			{
				Name: "Tabexpand",
				DataType: "bool",
			},						
		},
		Specs: specs,
	}
	
	blogf := fmt.Sprintf("%s/%s", writePath, "config.go")
	f, err := os.Create(blobf)
	if err != nil {
		return err
	}
	defer f.Close()
	
	w := &bytes.Buffer{}
	
	t, err := template.New("nyne").ParseFiles(tmplPath)
	if err != nil {
		return err
	}
	if err := t.Execute(w, tc); err != nil {
		return err
	}
	data, err := format.Source(builder.Bytes())
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(blobf, data, os.ModePerm); err != nil {
		return err
	}
	return nil
}