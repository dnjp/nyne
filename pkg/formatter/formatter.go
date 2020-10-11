package formatter

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	// "path/filepath"

	"git.sr.ht/~danieljamespost/nyne/pkg/event"
	"git.sr.ht/~danieljamespost/nyne/util/io"
	"git.sr.ht/~danieljamespost/nyne/util/config"
)

// Formatter listens for Acme events and applies formatting rules to the active buffer
type Formatter interface {
	Run()
	ExecCmds(event event.Event, cmds []config.Command, ext string) error
	WriteMenu(w *event.Win) error
	SetupFormatting(*event.Win, *config.Spec) error
	Refmt(event.Event, config.Command, string) ([]byte, error)
}

// Nyne implements the Formatter inferface for $NYNERULES
type Nyne struct {
	specs         map[string]*config.Spec
	menu          []string
	listener      event.Listener
	debug         bool
	extWithoutDot []string
}

// New constructs a Formatter that uses $NYNERULES for formatting
func New(conf *config.Config) Formatter {
	n := &Nyne{
		specs:         make(map[string]*config.Spec),
		menu:          conf.Tag.Menu,
		listener:      event.NewListener(),
		debug:         len(os.Getenv("DEBUG")) > 0,
		extWithoutDot: []string{},
	}

	for _, spec := range conf.Format {
		copy := &config.Spec{
			Indent: spec.Indent,
			Tabexpand: spec.Tabexpand,
			Extensions: spec.Extensions,
			Commands: spec.Commands,
		}		
		for _, ext := range copy.Extensions {
			if !strings.Contains(ext, ".") {
				n.extWithoutDot = append(n.extWithoutDot, ext)
			}
			n.specs[ext] = copy
		}
	}

	n.listener.RegisterWinHook(event.WinHook{
		Handler: func(w *event.Win) {
			spec, _ := n.getSpec(w.File)
			if spec != nil {
				n.SetupFormatting(w, spec)
			}
			err := n.WriteMenu(w)
			if err != nil {
				io.Error(err)
			}
		},
	})

	n.listener.RegisterPutHook(event.PutHook{
		Handler: func(evt event.Event) event.Event {
			spec, ext := n.getSpec(evt.File)
			if spec == nil {
				return evt
			}
			err := n.ExecCmds(evt, spec.Commands, ext)
			if err != nil {
				io.Error(err)
				return evt
			}
			return evt
		},
	})

	km := &Keymap{
		GetWinFn: func(id int) (*event.Win, error) {
			l := n.listener.GetBufListener(id)
			if l == nil {
				return nil, fmt.Errorf("could not find event loop")
			}
			return l.GetWin(), nil
		},
		GetIndentFn: func(evt event.Event) int {
			spec, _ := n.getSpec(evt.File)
			if spec == nil {
				return 8 // default
			}
			return spec.Indent
		},
	}

	// Tabexpand
	n.listener.RegisterKeyCmdHook(km.Tabexpand(func(evt event.Event) bool {
		spec, _ := n.getSpec(evt.File)
		if spec == nil {
			return false
		}
		return spec.Tabexpand
	}))

	return n
}

func (n *Nyne) getSpec(file string) (*config.Spec, string) {
	ext := n.getExt(file, ".txt")
	spec := n.specs[ext]
	return spec, ext
}

// Run tells the Formatter to begin listening for Acme events
func (n *Nyne) Run() {
	log.Fatal(n.listener.Listen())
}

// ExecCmds executes commands that operate on stdin/stdout against the Acme buffer
func (n *Nyne) ExecCmds(evt event.Event, cmds []config.Command, ext string) error {
	updates := [][]byte{}
	for _, cmd := range cmds {
		new, err := n.Refmt(evt, cmd, ext)
		if err != nil {
			fmt.Println("ERROR: ", err)
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
func (n *Nyne) SetupFormatting(w *event.Win, spec *config.Spec) error {
	if w == nil {
		return fmt.Errorf("state has drifted: *event.Win is nil")
	}
	if spec.Indent == 0 {
		return nil
	}
	if err := w.WriteToTag("\n"); err != nil {
		return err
	}
	if err := w.ExecInTag("Tab", strconv.Itoa(spec.Indent)); err != nil {
		return err
	}
	return nil
}

// Refmt executes a command to the Acme buffer and refreshes the buffer with updated contents
func (n *Nyne) Refmt(evt event.Event, cmd config.Command, ext string) ([]byte, error) {
	l := n.listener.GetBufListener(evt.ID)
	if l == nil {
		return []byte{}, fmt.Errorf("no event loop found")
	}
	
	// get current body
	old, err := l.GetWin().ReadBody()
	if err != nil {
		return []byte{}, err
	}
	
	// write current body to temporary file
	tmp, err := ioutil.TempFile("", fmt.Sprintf("*%s", ext))
	if err != nil {
		return []byte{}, err
	}
	defer os.Remove(tmp.Name())
	
	if _, err := tmp.Write(old); err != nil {
		return []byte{}, err
	}
	
	// replace name with the temporary file
	nargs := replaceName(cmd.Args, tmp.Name())
	out, err := exec.Command(cmd.Exec, nargs...).CombinedOutput()
	if err != nil {
		return []byte{}, err
	}

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
	l := n.listener.GetBufListener(evt.ID)
	if l == nil {
		return fmt.Errorf("no event loop found")
	}
	w := l.GetWin()
	for _, update := range updates {
		if err := w.SetAddr(","); err != nil {
			return err
		}
		if err := w.SetData(update); err != nil {
			return err
		}
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

func (n *Nyne) getExt(in string, def string) string {
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
