package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"os/user"
	"text/template"

	"git.sr.ht/~danieljamespost/nyne/util/config"
)

type TC struct {
	Ext          string
	Indent       int
	Tabexpand    bool
	CommentStyle string
}

var tmpl string = `
// GENERATED CODE - DO NOT EDIT
package gen

import "strings"

// Spec contains the formatting specification for a given file extension
type Spec struct {
	Ext string
	Indent int
	Tabexpand bool
	CommentStyle string
}

// Conf maps file extensions to their formatting specification
var Conf = map[string]Spec{
{{ range . }}
	"{{ .Ext}}": {
		Indent: {{ .Indent }},
		Tabexpand: {{ .Tabexpand }},
		CommentStyle: "{{ .CommentStyle }}",
	},
{{ end }}
}

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

func GetFileName(in string) string {
	path := strings.Split(in, "/")
	return path[len(path)-1]
}
`

func main() {

	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	cfgPath := fmt.Sprintf("%s/.config/nyne/nyne.toml", usr.HomeDir)
	npath := os.Getenv("NYNERULES")
	if len(npath) > 0 {
		cfgPath = npath
	}
	conf, err := config.Load(cfgPath)
	if err != nil {
		panic(err)
	}

	specs := []TC{}
	for _, spec := range conf.Format {
		for _, ext := range spec.Extensions {
			ts := TC{
				Ext:          ext,
				CommentStyle: spec.CommentStyle,
				Indent:       spec.Indent,
				Tabexpand:    spec.Tabexpand,
			}
			specs = append(specs, ts)
		}
	}

	t, err := template.New("").Parse(tmpl)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, specs)
	if err != nil {
		panic(err)
	}
	out, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
}
