package format

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/dnjp/nyne/event"
	"github.com/dnjp/nyne/util/io"
)

// Nyne implements the Formatter inferface for $NYNERULES
type Nyne struct {
	listener      event.Listener
	debug         bool
	extWithoutDot []string
}

// New constructs a Formatter that uses $NYNERULES for formatting
func New(filetypes []Filetype) (*Nyne, error) {
	err := UpdateConfig(filetypes, Config)
	if err != nil {
		return nil, err
	}
	n := &Nyne{
		listener:      event.NewListener(),
		debug:         len(os.Getenv("DEBUG")) > 0,
		extWithoutDot: []string{},
	}
	n.registerHooks(n.listener)
	return n, nil
}

func (n *Nyne) registerHooks(l event.Listener) {
	l.RegisterWinHook(event.WinHook{
		Handler: func(w *event.Win) {
			spec, _ := n.getSpec(w.File)
			if spec.Tabwidth != 0 {
				n.SetupFormatting(w, spec)
			}
			err := n.WriteMenu(w)
			if err != nil {
				io.Error(err)
			}
		},
	})

	l.RegisterPutHook(event.PutHook{
		Handler: func(evt event.Event) event.Event {
			fn := func(e event.Event) error {
				spec, ext := n.getSpec(evt.File)
				if spec.Tabwidth == 0 {
					return nil
				}
				err := n.ExecCmds(evt, spec.Commands, ext)
				if err != nil {
					log.Println(err)
				}
				return nil
			}
			evt.PostHooks = append(evt.PostHooks, fn)
			return evt
		},
	})

	// Tabexpand
	l.RegisterKeyCmdHook(Tabexpand(
		func(evt event.Event) bool {
			spec, _ := n.getSpec(evt.File)
			return spec.Tabexpand
		},
		func(id int) (*event.Win, error) {
			l := l.BufListener(id)
			if l == nil {
				return nil, fmt.Errorf("could not find event loop")
			}
			return l.Win(), nil
		},
		func(evt event.Event) int {
			spec, _ := n.getSpec(evt.File)
			if spec.Tabwidth == 0 {
				return 8 // default
			}
			return spec.Tabwidth
		},
	))
}

// Run tells the Formatter to begin listening for Acme events
func (n *Nyne) Run() {
	log.Fatal(n.listener.Listen())
}

// ExecCmds executes commands that operate on stdin/stdout against the
// Acme buffer
func (n *Nyne) ExecCmds(evt event.Event, cmds []Command, ext string) error {
	updates := [][]byte{}
	for _, cmd := range cmds {
		new, err := n.Refmt(evt, cmd, ext)
		if err != nil {
			return err
		}
		updates = append(updates, new)
	}
	return n.WriteUpdates(evt, updates)
}

// WriteMenu writes the specified menu options to the Acme buffer
func (n *Nyne) WriteMenu(w *event.Win) error {
	if w == nil {
		return fmt.Errorf("state has drifted: *event.Win is nil")
	}

	builtin := []string{"Put", "Undo", "Redo"}
	for _, opt := range builtin {
		cmd := fmt.Sprintf("  %s", opt)
		if err := w.WriteToTag(cmd); err != nil {
			return err
		}
	}

	if err := w.WriteToTag("\n"); err != nil {
		return err
	}

	for _, opt := range Menu {
		if strings.Contains(opt, " ") {
			opt = fmt.Sprintf("(%s)", opt)
		}
		cmd := fmt.Sprintf("  %s", opt)
		if err := w.WriteToTag(cmd); err != nil {
			return err
		}
	}
	return nil
}

// SetupFormatting opens the Acme buffer for writing and applies the
// indentation and tab expansion options provided in $NYNERULES
func (n *Nyne) SetupFormatting(w *event.Win, spec Filetype) error {
	if w == nil {
		return fmt.Errorf("state has drifted: *event.Win is nil")
	}
	if spec.Tabwidth == 0 {
		return nil
	}
	if err := w.WriteToTag("\n"); err != nil {
		return err
	}
	if err := w.ExecInTag("Tab", strconv.Itoa(spec.Tabwidth)); err != nil {
		return err
	}
	if spec.Tabexpand {
		if err := w.ExecInTag("tabexpand=true"); err != nil {
			return err
		}
	}
	return nil
}

// Refmt executes a command to the Acme buffer and refreshes the buffer
// with updated contents
func (n *Nyne) Refmt(evt event.Event, cmd Command, xt string) ([]byte, error) {
	l := n.listener.BufListener(evt.ID)
	if l == nil {
		return []byte{}, fmt.Errorf("no event loop found")
	}

	// get current body
	old, err := l.Win().ReadBody()
	if err != nil {
		return []byte{}, err
	}

	var nargs []string
	var tmp *os.File
	if cmd.PrintsToStdout {
		nargs = replaceName(cmd.Args, l.File())
	} else {
		// write current body to temporary file
		tmp, err = ioutil.TempFile("", fmt.Sprintf("*%s", xt))
		if err != nil {
			return []byte{}, err
		}
		defer os.Remove(tmp.Name())
		if _, err = tmp.Write(old); err != nil {
			return []byte{}, err
		}

		// replace name with the temporary file
		nargs = replaceName(cmd.Args, tmp.Name())
	}

	// Execute the command
	out, err := exec.Command(cmd.Exec, nargs...).CombinedOutput()
	if err != nil {
		return []byte{}, fmt.Errorf("Error: %+v\n%s", err, string(out))
	}

	// handle formatting commands that both do and do not write to stdout
	var new []byte
	if cmd.PrintsToStdout {
		new = out
	} else {
		// read the temporary file that has been written to
		new, err = ioutil.ReadFile(tmp.Name())
		if err != nil {
			return []byte{}, err
		}
	}
	return new, nil
}

// WriteUpdates writes the updated contents to the file
func (n *Nyne) WriteUpdates(evt event.Event, updates [][]byte) error {
	l := n.listener.BufListener(evt.ID)
	if l == nil {
		return fmt.Errorf("no event loop found")
	}
	w := l.Win()
	for _, update := range updates {
		if err := w.SetAddr(","); err != nil {
			return err
		}
		if err := w.SetData(update); err != nil {
			return err
		}
		// prevent index out of bounds error
		if w.Lastpoint > len(update) {
			w.Lastpoint = len(update)
		}
		if err := w.SetAddr(fmt.Sprintf("#%d", w.Lastpoint)); err != nil {
			return err
		}
		if err := w.SetTextToAddr(); err != nil {
			return err
		}
		if err := w.ExecShow(); err != nil {
			return err
		}
	}
	w.WriteEvent(evt)
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

func (n *Nyne) getSpec(file string) (Filetype, string) {
	ext := Extension(file, ".txt")
	spec := Config[ext]
	return spec, ext
}
