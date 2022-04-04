package nyne

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
)

// Formatter formats acme windows and buffers
type Formatter struct {
	acme   *Acme
	debug  bool
	config map[string]Filetype
}

// NewFormatter constructs a Formatter
func NewFormatter(filetypes []Filetype, menutag []string) (*Formatter, error) {
	f := &Formatter{
		acme:   NewAcme(),
		debug:  len(os.Getenv("DEBUG")) > 0,
		config: make(map[string]Filetype),
	}
	err := FillFiletypes(f.config, filetypes)
	if err != nil {
		return nil, err
	}

	f.acme.WinHooks = map[Text][]WinHandler{
		New: {
			func(w *Win) {
				ft, _ := f.filetype(w.File)
				if ft.Tabwidth != 0 {
					f.fmt(w, ft)
				}
				for _, opt := range menutag {
					err := w.AppendTag(opt)
					if err != nil {
						panic(err)
					}
				}
			},
		},
	}

	f.acme.EventHooks = map[Text][]Handler{
		Put: {
			func(evt Event) (Event, bool) {
				evt.WriteHooks = append(evt.WriteHooks, func(e Event) error {
					ft, ext := f.filetype(evt.File)
					if ft.Tabwidth == 0 {
						return nil
					}
					err := f.exec(evt, ft.Commands, ext)
					if err != nil {
						log.Println(err)
					}
					return nil
				})
				return evt, true
			},
		},
	}

	key, expand := Tabexpand(
		func(evt Event) bool {
			ft, _ := f.filetype(evt.File)
			return ft.Tabexpand
		},
		func(id int) (*Win, error) {
			l := f.acme.Buf(id)
			if l == nil {
				return nil, fmt.Errorf("could not find event loop")
			}
			return l.Win(), nil
		},
		func(evt Event) int {
			ft, _ := f.filetype(evt.File)
			if ft.Tabwidth == 0 {
				return 8 // default
			}
			return ft.Tabwidth
		})
	f.acme.KeyHooks = map[rune]Handler{
		key: expand,
	}

	return f, nil
}

// Run tells the Formatter to begin listening for Acme events
func (f *Formatter) Run() error {
	return f.acme.Listen()
}

// exec executes commands that operate on stdin/stdout against the
// Acme buffer
func (f *Formatter) exec(evt Event, cmds []Command, ext string) error {
	updates := [][]byte{}
	for _, cmd := range cmds {
		new, err := f.refmt(evt, cmd, ext)
		if err != nil {
			return err
		}
		updates = append(updates, new)
	}
	return f.update(evt, updates)
}

// fmt opens the Acme buffer for writing and applies the
// indentation and tab expansion options provided in $NYNERULES
func (f *Formatter) fmt(w *Win, ft Filetype) error {
	if w == nil {
		return fmt.Errorf("state has drifted: *Win is nil")
	}
	if ft.Tabwidth == 0 {
		return nil
	}
	if err := w.AppendTag("\n"); err != nil {
		return err
	}
	if err := w.Exec("Tab", strconv.Itoa(ft.Tabwidth)); err != nil {
		return err
	}
	if ft.Tabexpand {
		if err := w.Exec("tabexpand=true"); err != nil {
			return err
		}
	}
	return nil
}

// refmt executes a command to the Acme buffer and refreshes the buffer
// with updated contents
func (f *Formatter) refmt(evt Event, cmd Command, xt string) ([]byte, error) {
	l := f.acme.Buf(evt.ID)
	if l == nil {
		return []byte{}, fmt.Errorf("no event loop found")
	}

	// get current body
	old, err := l.Win().Body()
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
		return []byte{}, fmt.Errorf("error: %+v\n%s", err, string(out))
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

// update writes the updated contents to the file
func (f *Formatter) update(evt Event, updates [][]byte) error {
	l := f.acme.Buf(evt.ID)
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
		if err := w.SetAddr("#%d", w.Lastpoint); err != nil {
			return err
		}
		if err := w.SelectionFromAddr(); err != nil {
			return err
		}
		if err := w.Show(); err != nil {
			return err
		}
	}
	w.WriteEvent(evt)
	return nil
}

func (f *Formatter) filetype(file string) (Filetype, string) {
	ext := Extension(file, ".txt")
	ft := f.config[ext]
	return ft, ext
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
