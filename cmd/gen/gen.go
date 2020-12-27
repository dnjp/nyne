package main

import (
	"bytes"
	"fmt"
	"go/format"
	"text/template"

	"github.com/dnjp/nyne/util/config"
)

// GenConf contains the specification for generating a static
// configuration
type GenConf struct {
	Menu  []string
	Specs []GenSpec
}

// GenSpec is the configuration for the generated formatting specification
type GenSpec struct {
	Ext       string
	Indent    int
	Tabexpand bool
	CmtStyle  string
	Cmds      []config.Command
}

func main() {
	specs := []GenSpec{}
	for _, spec := range Cfg.Format {
		for _, ext := range spec.Extensions {
			ts := GenSpec{
				Ext:       ext,
				CmtStyle:  spec.CommentStyle,
				Indent:    spec.Indent,
				Tabexpand: spec.Tabexpand,
				Cmds:      spec.Commands,
			}
			specs = append(specs, ts)
		}
	}
	cfg := GenConf{
		Menu:  Cfg.Tag.Menu,
		Specs: specs,
	}

	t, err := template.New("").Parse(tmpl)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, cfg)
	if err != nil {
		panic(err)
	}
	out, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
}

var tmpl string = `
// Package gen CONTAINS GENERATED CODE - DO NOT EDIT
package gen

import "strings"

// Menu contains the menu options that should be written to the scratch
// buffer
var Menu = []string{
	{{ range .Menu }}
	"{{ . }}",
	{{ end }}
}

// Cmd contains options for executing a given command against an
// acme window
type Cmd struct {
	Exec           string
	Args           []string
	PrintsToStdout bool
}

// Spec contains the formatting specification for a given file extension
type Spec struct {
	Ext string
	Indent int
	Tabexpand bool
	CmtStyle string
	Cmds []Cmd
}

// Conf maps file extensions to their formatting specification
var Conf = map[string]Spec{
{{ range .Specs }}
	"{{ .Ext}}": {
		Indent: {{ .Indent }},
		Tabexpand: {{ .Tabexpand }},
		CmtStyle: "{{ .CmtStyle }}",
		Cmds: []Cmd{
			{{ range .Cmds }}
			{
				Exec: "{{ .Exec }}",
				Args: []string{
					{{ range .Args }}
					"{{ . }}",
					{{ end }}
				},
				PrintsToStdout: {{ .PrintsToStdout }},
			},
			{{ end }}
		},
	},
{{ end }}
}

// GetExt parses the file extension given a file name
func GetExt(in string, def string) string {
	filename := GetFileName(in)
	if !strings.Contains(filename, ".") {
		return filename
	}
	pts := strings.Split(filename, ".")
	if len(pts) == len(in) {
		return def
	}
	return "." + pts[len(pts)-1]
}

// GetFileName takes the absolute path to a file and returns just the name
// of the file
func GetFileName(in string) string {
	path := strings.Split(in, "/")
	return path[len(path)-1]
}`
