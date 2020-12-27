// Package gen CONTAINS GENERATED CODE - DO NOT EDIT
package gen

import "strings"

// Menu contains the menu options that should be written to the scratch
// buffer
var Menu = []string{

	"|fmt",

	"|com",

	"|a-",

	"|a+",

	"Ldef",

	"Lrefs",

	"Lcomp",

	"win",
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
	Ext       string
	Indent    int
	Tabexpand bool
	CmtStyle  string
	Cmds      []Cmd
}

// Conf maps file extensions to their formatting specification
var Conf = map[string]Spec{

	".cc": {
		Indent:    2,
		Tabexpand: true,
		CmtStyle:  "// ",
		Cmds:      []Cmd{},
	},

	".cpp": {
		Indent:    2,
		Tabexpand: true,
		CmtStyle:  "// ",
		Cmds:      []Cmd{},
	},

	".hpp": {
		Indent:    2,
		Tabexpand: true,
		CmtStyle:  "// ",
		Cmds:      []Cmd{},
	},

	".cxx": {
		Indent:    2,
		Tabexpand: true,
		CmtStyle:  "// ",
		Cmds:      []Cmd{},
	},

	".hxx": {
		Indent:    2,
		Tabexpand: true,
		CmtStyle:  "// ",
		Cmds:      []Cmd{},
	},

	".go": {
		Indent:    8,
		Tabexpand: false,
		CmtStyle:  "// ",
		Cmds: []Cmd{

			{
				Exec: "gofmt",
				Args: []string{

					"$NAME",
				},
				PrintsToStdout: true,
			},
		},
	},

	".md": {
		Indent:    2,
		Tabexpand: true,
		CmtStyle:  "",
		Cmds: []Cmd{

			{
				Exec: "prettier",
				Args: []string{

					"--print-width",

					"80",

					"--prose-wrap",

					"always",

					"--write",

					"$NAME",
				},
				PrintsToStdout: false,
			},
		},
	},

	".rc": {
		Indent:    8,
		Tabexpand: false,
		CmtStyle:  "# ",
		Cmds:      []Cmd{},
	},

	".sh": {
		Indent:    8,
		Tabexpand: false,
		CmtStyle:  "# ",
		Cmds:      []Cmd{},
	},

	".tf": {
		Indent:    2,
		Tabexpand: true,
		CmtStyle:  "# ",
		Cmds: []Cmd{

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

	".toml": {
		Indent:    8,
		Tabexpand: false,
		CmtStyle:  "# ",
		Cmds:      []Cmd{},
	},

	".c": {
		Indent:    8,
		Tabexpand: false,
		CmtStyle:  "/* */",
		Cmds: []Cmd{

			{
				Exec: "indent",
				Args: []string{

					"-ncdb",

					"-cp33",

					"-c33",

					"-cd33",

					"-sc",

					"-nsaf",

					"-nsai",

					"-nsaw",

					"-npcs",

					"-nprs",

					"-psl",

					"-br",

					"-bls",

					"-ce",

					"-nsob",

					"-nss",

					"-nbad",

					"-nbap",

					"-bbo",

					"-bc",

					"-hnl",

					"-ts8",

					"-ci4",

					"-cli0",

					"-cbi0",

					"-sbi0",

					"-bli0",

					"-di16",

					"-i8",

					"-ip8",

					"-l75",

					"-lp",

					"-st",

					"$NAME",
				},
				PrintsToStdout: true,
			},

			{
				Exec: "sed",
				Args: []string{

					"s/) {/){/g",

					"$NAME",
				},
				PrintsToStdout: true,
			},
		},
	},

	".h": {
		Indent:    8,
		Tabexpand: false,
		CmtStyle:  "/* */",
		Cmds: []Cmd{

			{
				Exec: "indent",
				Args: []string{

					"-ncdb",

					"-cp33",

					"-c33",

					"-cd33",

					"-sc",

					"-nsaf",

					"-nsai",

					"-nsaw",

					"-npcs",

					"-nprs",

					"-psl",

					"-br",

					"-bls",

					"-ce",

					"-nsob",

					"-nss",

					"-nbad",

					"-nbap",

					"-bbo",

					"-bc",

					"-hnl",

					"-ts8",

					"-ci4",

					"-cli0",

					"-cbi0",

					"-sbi0",

					"-bli0",

					"-di16",

					"-i8",

					"-ip8",

					"-l75",

					"-lp",

					"-st",

					"$NAME",
				},
				PrintsToStdout: true,
			},

			{
				Exec: "sed",
				Args: []string{

					"s/) {/){/g",

					"$NAME",
				},
				PrintsToStdout: true,
			},
		},
	},

	".java": {
		Indent:    2,
		Tabexpand: true,
		CmtStyle:  "// ",
		Cmds: []Cmd{

			{
				Exec: "google-java-format",
				Args: []string{

					"$NAME",
				},
				PrintsToStdout: true,
			},
		},
	},

	".js": {
		Indent:    2,
		Tabexpand: true,
		CmtStyle:  "// ",
		Cmds: []Cmd{

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

	".ts": {
		Indent:    2,
		Tabexpand: true,
		CmtStyle:  "// ",
		Cmds: []Cmd{

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

	".json": {
		Indent:    2,
		Tabexpand: true,
		CmtStyle:  "",
		Cmds:      []Cmd{},
	},

	"Makefile": {
		Indent:    8,
		Tabexpand: false,
		CmtStyle:  "# ",
		Cmds:      []Cmd{},
	},

	".yml": {
		Indent:    2,
		Tabexpand: true,
		CmtStyle:  "# ",
		Cmds:      []Cmd{},
	},

	".yaml": {
		Indent:    2,
		Tabexpand: true,
		CmtStyle:  "# ",
		Cmds:      []Cmd{},
	},

	".html": {
		Indent:    2,
		Tabexpand: true,
		CmtStyle:  "<!-- -->",
		Cmds:      []Cmd{},
	},
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
}
