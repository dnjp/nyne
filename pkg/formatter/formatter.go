package formatter

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"git.sr.ht/~danieljamespost/nyne/pkg/event"
	"git.sr.ht/~danieljamespost/nyne/util/config"
)

// Formatter listens for Acme events and applies formatting rules to the active buffer
type Formatter interface {
	Run()
	ExecCmds(event *event.Event, commands []config.Command, ext string) error
	WriteMenu(w *event.Win) error
	SetupFormatting(*event.Win, config.Format) error
	Refmt(*event.Event, string, []string, string) ([]byte, error)
}

// NFmt implements the Formatter inferface for $NYNERULES
type NFmt struct {
	ops      map[string]*Op
	menu     []string
	listener event.Listener
	debug bool
	extWithoutDot []string
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
		debug: len(os.Getenv("DEBUG")) > 0,
		extWithoutDot: []string{"Makefile"},
	}

	for _, spec := range conf.Spec {
		for _, ext := range spec.Ext {
			n.ops[ext] = &Op{
				Fmt: spec.Fmt,
				Cmd: spec.Cmd,
			}
		}
	}

	n.listener.RegisterOpenHook(event.OpenHook{
		Op: event.NEW,
		Handler: func(w *event.Win) {
			op, _ := n.getOp(w.File)
			if op != nil {
				n.SetupFormatting(w, op.Fmt)
			}
			n.WriteMenu(w)
		},
	})

	n.listener.RegisterHook(event.Hook{
		Op: event.PUT,
		Handler: func(evt *event.Event) *event.Event {
			op, ext := n.getOp(evt.File)
			if op == nil {
				return evt
			}
			n.ExecCmds(evt, op.Cmd, ext)
			return evt
		},
	})


	return n
}


func (n *NFmt) getOp(file string) (*Op, string) {
	ext := n.getExt(file, ".txt")
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
func (n *NFmt) WriteMenu(w *event.Win) error {
	if err := w.WriteToTag("\n"); err != nil {
		return err
	}
	for _, opt := range n.menu {
		cmd := fmt.Sprintf("  %s", opt)
		if err := w.WriteToTag(cmd); err != nil {
			return err
		}
	}
	return nil
}

// SetupFormatting opens the Acme buffer for writing and applies the indentation and
// tab expansion options provided in $NYNERULES
func (f *NFmt) SetupFormatting(w *event.Win, format config.Format) error {
	if format.Indent == 0 {
		return nil
	}

	if err := w.ExecInTag("Tab", strconv.Itoa(format.Indent)); err != nil {
		return err
	}

	if format.Expand {
		if err := w.ExecInTag("nynetab", strconv.Itoa(format.Indent)); err != nil {
			return err
		}
	}
	return nil
}

// Refmt executes a command to the Acme buffer and refreshes the buffer with updated contents
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
	for _, update := range updates {
		if err := evt.Win.SetAddr(","); err != nil {
			return err
		}
		if err := evt.Win.SetData(update); err != nil {
			return err
		}
	}
	return nil
}

func (n *NFmt) resetView(evt *event.Event) error {
	if err := evt.Win.SetAddr("0,0"); err != nil {
		return err
	}
	if err := evt.Win.SetTextToAddr(); err != nil {
		return err
	}
	if err := evt.Win.ExecShow(); err != nil {
		return err
	}
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

func (n *NFmt) getExt(in string, def string) string {
	filename := getFileName(in)
	if includes(filename, n.extWithoutDot) {
		return filename
	}
	pts := strings.Split(filename, ".")
	if len(pts) == len(in) {
		return def
	}
	return "." + pts[len(pts)-1]
}

func getFileName(in string) string {
	path := strings.Split(in, "/")
	return path[len(path)-1]
}

func includes(in string, dat []string) bool {
	for _, val := range dat {
		if val == in {
			return true
		}
	}
	return false
}
