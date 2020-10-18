package main

import (
	"fmt"
	"os"
	"os/user"
	"text/template"	

	"git.sr.ht/~danieljamespost/nyne/util/config"
)

type TC struct {
	Ext string
	Indent int
	Tabexpand bool
	CommentStyle string
}

var tmpl string = `
package gen

type Spec struct {
	Ext string
	Indent int
	Tabexpand bool
	CommentStyle string
}

var  map[string]Spec = {
{{ range . }}
	"{{ .Ext}}": {
		Indent: {{ .Indent }},
		Tabexpand: {{ .Tabexpand }},
		CommentStyle: "{{ .CommentStyle }}",
	},
{{ end }}
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
				Ext: ext,
				CommentStyle: spec.CommentStyle,
				Indent: spec.Indent,
				Tabexpand: spec.Tabexpand,
			}
			specs = append(specs, ts)
		}
	}	

	t, err := template.New("test").Parse(tmpl)
	if err != nil { panic(err) }
	err = t.Execute(os.Stdout, specs)
	if err != nil { panic(err) }
}