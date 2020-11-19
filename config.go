package main

import (
	"git.sr.ht/~danieljamespost/nyne/util/config"
)

var Cfg config.Config = config.Config{
	Tag: config.Tag{
		Menu: []string{"|fmt", "|com", "|a-", "|a+", "Ldef", "Lrefs", "Lcomp", "win"},
	},
	Format: map[string]config.Spec{
		"c": config.Spec{
			CommentStyle: "/* */",
			Indent:       8,
			Tabexpand:    false,
			Extensions:   []string{".c", ".h"},
			Commands: []config.Command{
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
					Exec:           "sed",
					Args:           []string{"s/) {/){/g", "$NAME"},
					PrintsToStdout: true,
				},
			},
		},
		"cpp": config.Spec{
			CommentStyle: "// ",
			Indent:       2,
			Tabexpand:    true,
			Extensions:   []string{".cc", ".cpp", ".hpp", ".cxx", ".hxx"},
			Commands:     []config.Command{},
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
					Exec:           "prettier",
					Args:           []string{"$NAME", "--write", "--loglevel", "error"},
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
					Exec:           "terraform",
					Args:           []string{"fmt", "$NAME"},
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
