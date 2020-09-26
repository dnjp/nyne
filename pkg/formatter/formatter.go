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
	"strconv"

	"9fans.net/go/acme"
	"git.sr.ht/~danieljamespost/nyne/pkg/event"
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
		ops: make(map[string]*Op),
		menu: conf.Menu,
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
			if !n.shouldProcessEvent(evt) {
				return
			}
			n.Cmd(evt, op.Cmd, ext)
 		},
 	})

  	n.listener.RegisterHook(event.Hook{
 		Op: event.NEW,
 		Handler: func(evt *event.Event) {
			if !n.shouldProcessEvent(evt) {
				return
			}
			n.Fmt(evt, op.Fmt)
			n.WMenu(evt)
		},
 	})

	return n
}

func (n *NFmt) shouldProcessEvent(evt *event.Event) bool {
 	ext := getExt(evt.File, ".txt")
	op := n.ops[ext]
	if op == nil {
		return false
	}
	return true
}

// Run tells the Formatter to begin listening for Acme events
func (n *NFmt) Run() {
	log.Fatal(n.listener.Listen())
}

// Cmd executes commands that operate on stdin/stdout against the Acme buffer
// TODO: this should read the file once, create a unified diff, and apply the diff
//             to the buffer instead of doing so for each command
func (n *NFmt) Cmd(evt *event.Event, commands []config.Command, ext string) {
	for _, cmd := range commands {
		args := replaceName(cmd.Args, evt.File)
		n.Refmt(evt, cmd.Exec, args, ext)
    	}
}

// WMenu writes the specified menu options to the Acme buffer
func (n *NFmt) WMenu(evt *event.Event) {
	if err := evt.Win.WriteToTag("\n"); err != nil {
		printErr(err)
		return
	}
	for _, opt := range n.menu {
		cmd := fmt.Sprintf(" (%s)", opt)
		if err := evt.Win.WriteToTag(cmd); err != nil {
			printErr(err)
		}
	}
}

// Fmt opens the Acme buffer for writing and applies the indentation and
// tab expansion options provided in $NYNERULES
func (f* NFmt) Fmt(evt *acme.Event, format config.Format) {
	if format.Indent == 0 {
		return
	}

	if err := evt.Win.ClearTagText(); err != nil {
		printErr(err)
		return
	}

	if err := evt.Win.ExecInTag("Tab", strconv.Itoa(format.Indent)); err != nil {
		printErr(err)
		return
	}


	if format.Expand {
		if err := evt.Win.ClearTagText(); err != nil {
			printErr(err)
			return
		}

		if err := evt.Win.ExecInTag("nynetab", strconv.Itoa(format.Indent)); err != nil {
			printErr(err)
			return
		}
	}
}


// Refmt executes a command to the Acme buffer and refreshes the buffer with updated contents
// TODO: this implementation introduces a bug that breaks undo
func (n *NFmt) Refmt(evt *event.Event, x string, args []string, ext string) {

	old, err := evt.Win.ReadBody()
	if err != nil {
		printErr(err)
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

	evt.Win.SetAddr(",", nil)
	evt.Win.SetData(new)
// 	evt.Win.MarkWinClean()
// 	evt.Win.ExecGet()

}

func printErr(err error) {
	fmt.Fprintf(os.Stderr, "%v", err)
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

