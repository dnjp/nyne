package format

import (
	"fmt"
	"strings"
)

// Command contains options for executing a given command against an
// acme window
type Command struct {
	Exec           string
	Args           []string
	PrintsToStdout bool
}

// Filetype contains the formatting specification for a given file extension
type Filetype struct {
	Name       string
	Extensions []string
	Tabwidth   int
	Tabexpand  bool
	Comment    string
	Commands   []Command
}

// Extension parses the file extension given a file name
func Extension(in string, def string) string {
	filename := Filename(in)
	if !strings.Contains(filename, ".") {
		return filename
	}
	pts := strings.Split(filename, ".")
	if len(pts) == len(in) {
		return def
	}
	return "." + pts[len(pts)-1]
}

// Filename takes the absolute path to a file and returns just the name
// of the file
func Filename(in string) string {
	path := strings.Split(in, "/")
	return path[len(path)-1]
}

// FillFiletypes updates the filetypes in the given map
func FillFiletypes(dst map[string]Filetype, src []Filetype) error {
	for _, ft := range src {
		for _, ext := range ft.Extensions {
			if ft2, ok := dst[ext]; ok {
				return fmt.Errorf("duplicate extension for filetype: original=%+v new=%+v", ft2, ft)
			}
			dst[ext] = ft
		}
	}
	return nil
}
