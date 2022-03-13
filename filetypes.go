package nyne

import (
	"github.com/dnjp/nyne/format"
)

// config maps file extensions to their formatting specification
var config = func() map[string]format.Filetype {
	c := make(map[string]format.Filetype)
	err := format.FillFiletypes(c, Filetypes)
	if err != nil {
		panic(err)
	}
	return c
}()

// Filetype returns the filetype in the nyne config if present
func Filetype(ext string) (ft format.Filetype, ok bool) {
	ft, ok = config[ext]
	return
}

// Filetypes define file formatting rules that will be applied
var Filetypes = []format.Filetype{
	{
		Name:       "cpp",
		Extensions: []string{".cc", ".cpp", ".hpp", ".cxx", ".hxx"},
		Tabwidth:   2,
		Tabexpand:  true,
		Comment:    "// ",
		Commands:   []format.Command{},
	},
	{
		Name:       "java",
		Extensions: []string{".java"},
		Tabwidth:   2,
		Tabexpand:  true,
		Comment:    "// ",
		Commands:   []format.Command{
			// {
			// 	Exec: "google-java-format",
			// 	Args: []string{
			// 		"$NAME",
			// 	},
			// 	PrintsToStdout: true,
			// },
		},
	},
	{
		Name:       "javascript",
		Extensions: []string{".js", ".ts"},
		Tabwidth:   2,
		Tabexpand:  true,
		Comment:    "// ",
		Commands: []format.Command{
			{
				Exec: "prettier",
				Args: []string{
					"$NAME",
					"--write",
					"--loglevel",
					"error",
				},
				PrintsToStdout: false,
			},
		},
	},
	{
		Name:       "json",
		Extensions: []string{".json"},
		Tabwidth:   2,
		Tabexpand:  true,
		Comment:    "",
		Commands:   []format.Command{},
	},
	{
		Name:       "makefile",
		Extensions: []string{"Makefile"},
		Tabwidth:   8,
		Tabexpand:  false,
		Comment:    "# ",
		Commands:   []format.Command{},
	},
	{
		Name:       "shell",
		Extensions: []string{".rc", ".sh"},
		Tabwidth:   8,
		Tabexpand:  false,
		Comment:    "# ",
		Commands:   []format.Command{},
	},
	{
		Name:       "c",
		Extensions: []string{".c", ".h"},
		Tabwidth:   8,
		Tabexpand:  false,
		Comment:    "/* */",
		Commands:   []format.Command{},
	},
	{
		Name:       "html",
		Extensions: []string{".html"},
		Tabwidth:   2,
		Tabexpand:  true,
		Comment:    "<!-- -->",
		Commands:   []format.Command{},
	},
	{
		Name:       "markdown",
		Extensions: []string{".md"},
		Tabwidth:   2,
		Tabexpand:  true,
		Comment:    "",
		Commands:   []format.Command{
			// {
			// 	Exec: "prettier",
			// 	Args: []string{
			// 		"--print-width",
			// 		"80",
			// 		"--prose-wrap",
			// 		"always",
			// 		"--write",
			// 		"$NAME",
			// 	},
			// 	PrintsToStdout: false,
			// },
		},
	},
	{
		Name:       "terraform",
		Extensions: []string{".tf"},
		Tabwidth:   2,
		Tabexpand:  true,
		Comment:    "# ",
		Commands: []format.Command{
			{
				Exec: "terraform",
				Args: []string{
					"fmt",
					"$NAME",
				},
				PrintsToStdout: false,
			},
		},
	},
	{
		Name:       "toml",
		Extensions: []string{".toml"},
		Tabwidth:   8,
		Tabexpand:  false,
		Comment:    "# ",
		Commands:   []format.Command{},
	},
	{
		Name:       "yaml",
		Extensions: []string{".yml", ".yaml"},
		Tabwidth:   2,
		Tabexpand:  true,
		Comment:    "# ",
		Commands:   []format.Command{},
	},
	{
		Name:       "go",
		Extensions: []string{".go", "go.mod", "go.sum"},
		Tabwidth:   8,
		Tabexpand:  false,
		Comment:    "// ",
		Commands:   []format.Command{
			// {
			// 	Exec: "gofmt",
			// 	Args: []string{
			// 		"$NAME",
			// 	},
			// 	PrintsToStdout: true,
			// },
		},
	},
}
