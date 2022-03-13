package format

import (
	"fmt"
)

var DefaultFiletypes = []Filetype{
	{
		Name:       "cpp",
		Extensions: []string{".cc", ".cpp", ".hpp", ".cxx", ".hxx"},
		Tabwidth:   2,
		Tabexpand:  true,
		Comment:    "// ",
		Commands:   []Command{},
	},
	{
		Name:       "java",
		Extensions: []string{".java"},
		Tabwidth:   2,
		Tabexpand:  true,
		Comment:    "// ",
		Commands:   []Command{
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
		Commands: []Command{
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
		Commands:   []Command{},
	},
	{
		Name:       "makefile",
		Extensions: []string{"Makefile"},
		Tabwidth:   8,
		Tabexpand:  false,
		Comment:    "# ",
		Commands:   []Command{},
	},
	{
		Name:       "shell",
		Extensions: []string{".rc", ".sh"},
		Tabwidth:   8,
		Tabexpand:  false,
		Comment:    "# ",
		Commands:   []Command{},
	},
	{
		Name:       "c",
		Extensions: []string{".c", ".h"},
		Tabwidth:   8,
		Tabexpand:  false,
		Comment:    "/* */",
		Commands:   []Command{},
	},
	{
		Name:       "html",
		Extensions: []string{".html"},
		Tabwidth:   2,
		Tabexpand:  true,
		Comment:    "<!-- -->",
		Commands:   []Command{},
	},
	{
		Name:       "markdown",
		Extensions: []string{".md"},
		Tabwidth:   2,
		Tabexpand:  true,
		Comment:    "",
		Commands:   []Command{
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
		Commands: []Command{
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
		Commands:   []Command{},
	},
	{
		Name:       "yaml",
		Extensions: []string{".yml", ".yaml"},
		Tabwidth:   2,
		Tabexpand:  true,
		Comment:    "# ",
		Commands:   []Command{},
	},
	{
		Name:       "go",
		Extensions: []string{".go", "go.mod", "go.sum"},
		Tabwidth:   8,
		Tabexpand:  false,
		Comment:    "// ",
		Commands:   []Command{
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

// DefaultMenu contains the menu options that should be written to the scratch
// buffer
var DefaultMenu = []string{
	"|fmt",
	"|com",
	"|a-",
	"|a+",
	"Ldef",
	"Lrefs",
	"Lcomp",
	"win",
}

var UseDefaultFiletypes = true
var UseDefaultMenu = true

var Menu = func() []string {
	if UseDefaultMenu {
		return DefaultMenu
	}
	return make([]string, 0)
}()

// Config maps file extensions to their formatting specification
var Config = func() map[string]Filetype {
	config := make(map[string]Filetype)
	if UseDefaultFiletypes {
		err := UpdateConfig(DefaultFiletypes, config)
		if err != nil {
			panic(err)
		}
	}
	return config
}()

// UpdateConfig updates the filetypes in the global map
func UpdateConfig(filetypes []Filetype, config map[string]Filetype) error {
	for _, ft := range filetypes {
		for _, ext := range ft.Extensions {
			if ft2, ok := config[ext]; ok {
				return fmt.Errorf("duplicate extension for filetype: original=%+v new=%+v", ft2, ft)
			}
			config[ext] = ft
		}
	}
	return nil
}
