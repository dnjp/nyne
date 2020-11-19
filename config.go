package main

import (
	"git.sr.ht/~danieljamespost/nyne/util/config"
)

// Cfg specifies the configuration for nyne and its bundled utilities
var Cfg config.Config = config.Config{

	// Tag configures options for the Acme tag
	Tag: config.Tag{

		// Menu contain commands to be written to the acme
		// scratch area. Along with the default "Put", "Undo", and
		// "Redo" commands, these menu options will be written
		// to the acme scratch area when a new window is opened.
		Menu: []string{
			"|fmt",
			"|com",
			"|a-",
			"|a+",
			"Ldef",
			"Lrefs",
			"Lcomp",
			"win",
		},
	},

	// Format maps an identifier ("c", "go", etc.) to its formatting
	// specification. The identifier is an arbitrary name that useful
	// mostly to logically group formatting directives.
	Format: map[string]config.Spec{

		"c": config.Spec{

			// A string that contains the comment style for
			// the given language.  If the comment style has a
			// defined start and end comment structure (/* */
			// in C), then set commentstyle to the complete
			// comment structure like this: `commentstyle =
			// "/* */"`. com will infer that this means /*
			// should be placed at the beginning and */ should
			// be placed at the end.
			CommentStyle: "/* */",

			// The tab width used for indentation
			Indent: 8,

			// Determines whether to use hard tabs or spaces
			// for indentation
			Tabexpand: false,

			// A list of file extensions that nyne should apply
			// the given formatting rules to
			Extensions: []string{".c", ".h"},

			// The "commands" blocks is used to define the
			// external program to be run against against your
			// buffer on file save. Any number of these blocks
			// may be defined.
			Commands: []config.Command{
				{
					// A string representing the
					// executable used to format
					// the buffer
					Exec: "indent",

					// An array of strings containing the
					// arguments to the executable. $NAME
					// is a macro that will be replaced
					// with the absolute path to the
					// file you are working on. This
					// is a required argument.
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

					// A boolean representing whether
					// the executable will print to
					// stdout. If the command writes
					// the file in place, be sure to
					// set this to false.
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
		"cpp": config.Spec{
			CommentStyle: "// ",
			Indent:       2,
			Tabexpand:    true,
			Extensions: []string{
				".cc",
				".cpp",
				".hpp",
				".cxx",
				".hxx",
			},
			Commands: []config.Command{},
		},
		"go": config.Spec{
			CommentStyle: "// ",
			Indent:       8,
			Tabexpand:    false,
			Extensions:   []string{".go"},
			Commands: []config.Command{
				{
					Exec:           "gofmt",
					Args:           []string{"$NAME"},
					PrintsToStdout: true,
				},
			},
		},
		"html": config.Spec{
			CommentStyle: "<!-- -->",
			Indent:       2,
			Tabexpand:    true,
			Extensions:   []string{".html"},
			Commands:     []config.Command{},
		},
		"java": config.Spec{
			CommentStyle: "// ",
			Indent:       2,
			Tabexpand:    true,
			Extensions:   []string{".java"},
			Commands: []config.Command{
				{
					Exec:           "google-java-format",
					Args:           []string{"$NAME"},
					PrintsToStdout: true,
				},
			},
		},
		"js": config.Spec{
			CommentStyle: "// ",
			Indent:       2,
			Tabexpand:    true,
			Extensions:   []string{".js", ".ts"},
			Commands: []config.Command{
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
		"json": config.Spec{
			CommentStyle: "",
			Indent:       2,
			Tabexpand:    true,
			Extensions:   []string{".json"},
			Commands:     []config.Command{},
		},
		"makefile": config.Spec{
			CommentStyle: "# ",
			Indent:       8,
			Tabexpand:    false,
			Extensions:   []string{"Makefile"},
			Commands:     []config.Command{},
		},
		"markdown": config.Spec{
			CommentStyle: "",
			Indent:       2,
			Tabexpand:    true,
			Extensions:   []string{".md"},
			Commands: []config.Command{
				{
					Exec: "prettier",
					Args: []string{
						"--print-width", "80",
						"--prose-wrap", "always",
						"--write",
						"$NAME",
					},
					PrintsToStdout: false,
				},
			},
		},
		"shell": config.Spec{
			CommentStyle: "# ",
			Indent:       8,
			Tabexpand:    false,
			Extensions:   []string{".rc", ".sh"},
			Commands:     []config.Command{},
		},
		"tf": config.Spec{
			CommentStyle: "# ",
			Indent:       2,
			Tabexpand:    true,
			Extensions:   []string{".tf"},
			Commands: []config.Command{
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
		"toml": config.Spec{
			CommentStyle: "# ",
			Indent:       8,
			Tabexpand:    false,
			Extensions:   []string{".toml"},
			Commands:     []config.Command{},
		},
		"yaml": config.Spec{
			CommentStyle: "# ",
			Indent:       2,
			Tabexpand:    true,
			Extensions:   []string{".yml", ".yaml"},
			Commands:     []config.Command{},
		},
	},
}
