package formatter

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"unicode/utf8"

	"9fans.net/go/acme"
	"git.sr.ht/~danieljamespost/nyne/pkg/golang"
	"git.sr.ht/~danieljamespost/nyne/util/config"
)

// Formatter listens for Acme events and applies formatting rules to the active buffer
type Formatter interface {
	Listen()
	Cmd(event acme.LogEvent, commands []config.Command, ext string)
	WMenu(event acme.LogEvent)
	Fmt(event acme.LogEvent, format config.Format)
	Refmt(id int, name string, x string, args []string, ext string)
}

// NFmt implements the Formatter inferface for $NYNERULES
type NFmt struct {
	ops map[string]*Op
	menu []string
}

// Op specifies a formatting operation to be performed on an Acme buffer
type Op struct {
	Fmt config.Format
	Cmd []config.Command
}

// New constructs a Formatter that uses $NYNERULES for formatting
func New(conf *config.Config) Formatter {
	n := &NFmt{
		ops: make(map[string]*Op),
		menu: conf.Menu,
	}

	for _, spec := range conf.Spec {
    		for _, ext := range spec.Ext {
	    		n.ops[ext] = &Op{
				Fmt: spec.Fmt,
				Cmd: spec.Cmd,
			}
		}
    	}
	return n
}

// Listen tells the Formatter to begin listening for Acme events
func (n *NFmt) Listen() {
	l, err := acme.Log()
	if err != nil {
		log.Fatal(err)
	}
	for {
		event, err := l.Read()
		if err != nil {
			log.Fatal(err)
		}
		ext := getExt(event.Name, ".txt")
		op := n.ops[ext]
		if op == nil {
			continue
		}
		switch event.Op {
		case "put":
			n.Cmd(event, op.Cmd, ext)
		case "new":
			n.Fmt(event, op.Fmt)
			n.WMenu(event)
		}
	}
}

// Cmd executes commands that operate on stdin/stdout against the Acme buffer
// TODO: this should read the file once, create a unified diff, and apply the diff
//             to the buffer instead of doing so for each command
func (n *NFmt) Cmd(event acme.LogEvent, commands []config.Command, ext string) {
	for _, cmd := range commands {
		args := replaceName(cmd.Args, event.Name)
		n.Refmt(event.ID, event.Name, cmd.Exec, args, ext)
    	}
}

// WMenu writes the specified menu options to the Acme buffer
func (n *NFmt) WMenu(event acme.LogEvent) {
	w, err := acme.Open(event.ID, nil)
	if err != nil {
		log.Print(err)
		return
	}
	defer w.CloseFiles()
	if err := w.Fprintf("tag", "%s", "\n"); err != nil {
		log.Print(err)
	}

	for _, opt := range n.menu {
		cmd := fmt.Sprintf(" (%s)", opt)
		if err := w.Fprintf("tag", "%s", cmd); err != nil {
			log.Print(err)
		}
	}
}

// Fmt opens the Acme buffer for writing and applies the indentation and
// tab expansion options provided in $NYNERULES
func (f* NFmt) Fmt(event acme.LogEvent, format config.Format) {
	w, err := acme.Open(event.ID, nil)
	if err != nil {
		log.Print(err)
		return
	}
	defer w.CloseFiles()

	if format.Indent != 0 {
		tabCmd := fmt.Sprintf("Tab %d", format.Indent)
		if err := w.Ctl("cleartag"); err != nil {
			log.Print(err)
		}
		tag, err := w.ReadAll("tag")
		if err != nil {
			log.Print(err)
		}
		offset := utf8.RuneCount(tag)
		cmdlen := utf8.RuneCountInString(tabCmd)
		if err := w.Fprintf("tag", "%s", tabCmd); err != nil {
			log.Print(err)
		}
		evt := new(acme.Event)
		evt.C1 = 'M'
		evt.C2 = 'x'
		evt.Q0 = offset
		evt.Q1 = offset + cmdlen
		w.WriteEvent(evt)
	}

	if format.Expand == true {
		expCmd := fmt.Sprintf("nynetab %d", format.Indent)
		if err := w.Ctl("cleartag"); err != nil {
			log.Print(err)
		}
		tag, err := w.ReadAll("tag")
		if err != nil {
			log.Print(err)
		}
		offset := utf8.RuneCount(tag)
		cmdlen := utf8.RuneCountInString(expCmd)
		if err := w.Fprintf("tag", "%s", expCmd); err != nil {
			log.Print(err)
		}
		evt := new(acme.Event)
		evt.C1 = 'M'
		evt.C2 = 'x'
		evt.Q0 = offset
		evt.Q1 = offset + cmdlen
		w.WriteEvent(evt)
	}
}


// Refmt executes a command to the Acme buffer and refreshes the buffer with updated contents
// TODO: this implementation introduces a bug that breaks undo
func (n *NFmt) Refmt(id int, name string, x string, args []string, ext string) {
	w, err := acme.Open(id, nil)
	if err != nil {
		log.Print(err)
		return
	}
	defer w.CloseFiles()

	// TODO: read from 9p contents instead of raw file
	old, err := ioutil.ReadFile(name)
	if err != nil {
		return
	}
	new, err := exec.Command(x, args...).CombinedOutput()
	if err != nil {
		if strings.Contains(string(new), "fatal error") {
			fmt.Fprintf(os.Stderr, "%s %s: %v\n%s", x, name, err, new)
		} else {
			fmt.Fprintf(os.Stderr, "%s", new)
		}
		return
	}

	if bytes.Equal(old, new) {
		return
	}

	if ext != ".go" {
		w.Write("ctl", []byte("clean"))
		w.Write("ctl", []byte("get"))
		return
	} else {
		golang.Reformat(name, ext, w, old, new)
	}
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

