package nyne

import (
	"fmt"
	"regexp"
	"strconv"

	"9fans.net/go/draw"
)

var numre = regexp.MustCompile("[0-9]+")

// Menu contains the menu options that should be written to the tag
var Menu = []string{
	" Put  ", "Undo  ", "Redo  ", "win", "\n",
	"|com  ", "|a-  ", "|a+  ", "Ldef  ", "Lrefs  ", "Lcomp",
}

// Config maps file extensions to their formatting specification
var Config = func() map[string]Filetype {
	c := make(map[string]Filetype)
	err := FillFiletypes(c, Filetypes)
	if err != nil {
		panic(err)
	}
	return c
}()

// Nums extracts numbers from the string, returning all matches
func Nums(s string) (nums []int, err error) {
	for _, ns := range numre.FindAllString(s, -1) {
		n, err := strconv.Atoi(ns)
		if err != nil {
			return nil, err
		}
		nums = append(nums, n)
	}
	return
}

// FontSize returns the size of the font
func FontSize(font *draw.Font) (int, error) {
	sizes, err := Nums(font.Name)
	if err != nil {
		return 0, err
	}
	if l := len(sizes); l > 1 || l == 0 {
		return 0, fmt.Errorf("could not parse font size")
	}
	return sizes[0], nil
}

// IsHiDPI returns whether the font is being displayed in a HiDPI
// context. This is absolutely a hack and should not be trusted.
func IsHiDPI(font *draw.Font) bool {
	// I said this was a hack
	size, _ := FontSize(font)
	return size >= 24
}

// FindFiletype returns the filetype in the nyne config if present
func FindFiletype(filename string) (ft Filetype, ok bool) {
	ft, ok = Config[Extension(filename, ".txt")]
	return
}

// Filetypes define file formatting rules that will be applied
var Filetypes = []Filetype{
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
		Name:       "text",
		Extensions: []string{".txt"},
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
