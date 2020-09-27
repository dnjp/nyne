package formatter

import (
	"bytes"
	"fmt"
	"log"
	// "os"
	"os/exec"
	"strconv"
	"strings"
	// "unicode/utf8"

	"git.sr.ht/~danieljamespost/nyne/pkg/event"
	"git.sr.ht/~danieljamespost/nyne/util/config"
	// "git.sr.ht/~danieljamespost/nyne/util/io"
	// "github.com/sergi/go-diff/diffmatchpatch"
)

// Formatter listens for Acme events and applies formatting rules to the active buffer
type Formatter interface {
	Run()
	ExecCmds(event *event.Event, commands []config.Command, ext string) error
	WriteMenu(event *event.Event) error
	SetupFormatting(event *event.Event, format config.Format) error
	Refmt(*event.Event, string, []string, string) ([]byte, error)
}

// NFmt implements the Formatter inferface for $NYNERULES
type NFmt struct {
	ops      map[string]*Op
	menu     []string
	listener event.Listener
}

// Op specifies a formatting operation to be performed on an Acme buffer
type Op struct {
	Fmt config.Format
	Cmd []config.Command
}

// New constructs a Formatter that uses $NYNERULES for formatting
func New(conf *config.Config) Formatter {
	n := &NFmt{
		ops:      make(map[string]*Op),
		menu:     conf.Menu,
		listener: event.NewListener(),
	}

	for _, spec := range conf.Spec {
		for _, ext := range spec.Ext {
			n.ops[ext] = &Op{
				Fmt: spec.Fmt,
				Cmd: spec.Cmd,
			}
		}
	}

	n.listener.RegisterHook(event.Hook{
		Op: event.PUT,
		Handler: func(evt *event.Event) {
			op, ext := n.getOp(evt)
			if op == nil {
				return
			}
			n.ExecCmds(evt, op.Cmd, ext)
		},
	})

	n.listener.RegisterHook(event.Hook{
		Op: event.NEW,
		Handler: func(evt *event.Event) {
			op, _ := n.getOp(evt)
			if op == nil {
				return
			}
			n.SetupFormatting(evt, op.Fmt)
			n.WriteMenu(evt)
		},
	})

	return n
}

func (n *NFmt) getOp(evt *event.Event) (*Op, string) {
	ext := getExt(evt.File, ".txt")
	op := n.ops[ext]
	return op, ext
}

// Run tells the Formatter to begin listening for Acme events
func (n *NFmt) Run() {
	log.Fatal(n.listener.Listen())
}

// ExecCmds executes commands that operate on stdin/stdout against the Acme buffer
// TODO: this should read the file once, create a unified diff, and apply the diff
//             to the buffer instead of doing so for each command
func (n *NFmt) ExecCmds(evt *event.Event, commands []config.Command, ext string) error {
	updates := [][]byte{}
	for _, cmd := range commands {
		args := replaceName(cmd.Args, evt.File)
		new, err := n.Refmt(evt, cmd.Exec, args, ext)
		if err != nil {
			return err
		}
		updates = append(updates, new)
	}
	return n.WriteUpdates(evt, updates)
}

// WriteMenu writes the specified menu options to the Acme buffer
func (n *NFmt) WriteMenu(evt *event.Event) error {
	if err := evt.Win.WriteToTag("\n"); err != nil {
		return err
	}
	for _, opt := range n.menu {
		cmd := fmt.Sprintf(" (%s)", opt)
		if err := evt.Win.WriteToTag(cmd); err != nil {
			return err
		}
	}
	return nil
}

// SetupFormatting opens the Acme buffer for writing and applies the indentation and
// tab expansion options provided in $NYNERULES
func (f *NFmt) SetupFormatting(evt *event.Event, format config.Format) error {
	if format.Indent == 0 {
		return nil
	}

	if err := evt.Win.ClearTagText(); err != nil {
		return err
	}

	if err := evt.Win.ExecInTag("Tab", strconv.Itoa(format.Indent)); err != nil {
		return err
	}

	if format.Expand {
		if err := evt.Win.ClearTagText(); err != nil {
			return err
		}

		if err := evt.Win.ExecInTag("nynetab", strconv.Itoa(format.Indent)); err != nil {
			return err
		}
	}
	return nil
}

// Refmt executes a command to the Acme buffer and refreshes the buffer with updated contents
// TODO: this implementation introduces a bug that breaks undo
func (n *NFmt) Refmt(evt *event.Event, x string, args []string, ext string) ([]byte, error) {
	old, err := evt.Win.ReadBody()
	if err != nil {
		return []byte{}, err
	}
	new, err := exec.Command(x, args...).CombinedOutput()
	if err != nil {
		return []byte{}, err
	}

	if bytes.Equal(old, new) {
		return old, nil
	}

	return new, nil
}

func (n *NFmt) WriteUpdates(evt *event.Event, updates [][]byte) error {
	// 	original, err := n.Win.ReadBody()
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	unified := [][]diffmatchpatch.Diff
	// 	for i, update := range updates {
	// 		old := []byte{}
	// 		if i == 0 {
	// 			old = original
	// 		} else {
	// 			old = updates[i-1]
	// 		}
	//
	// 		dmp := diffmatchpatch.New()
	// 		diffs = dmp.DiffMain(string(old),string(update), false)
	// 		unified = append(unified, diffs)
	// 	}

	update := updates[len(updates)-1]

	if err := evt.Win.SetAddr(","); err != nil {
		return err
	}

	if err := evt.Win.SetData(update); err != nil {
		return err
	}

	if err := evt.Win.SetAddr("0,0"); err != nil {
		return err
	}

	if err := evt.Win.SetTextToAddr(); err != nil {
		return err
	}

	if err := evt.Win.ExecShow(); err != nil {
		return err
	}

	evt.Win.MarkWinClean()

	return nil
}

func replaceName(arr []string, name string) []string {
	newArr := make([]string, len(arr))
	for idx, str := range arr {
		if str == "$NAME" {
			newArr[idx] = name
		} else {
			newArr[idx] = arr[idx]
		}
	}
	return newArr
}

func getExt(in string, def string) string {
	pts := strings.Split(in, ".")
	if len(pts) == len(in) {
		return def
	}
	return "." + pts[len(pts)-1]
}
