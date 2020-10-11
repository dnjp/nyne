package event

import (
	"log"
	"fmt"
	"os"
	"strings"
	"sync"

	"9fans.net/go/acme"
)

// Listener can listen for acme Event and Window hooks
type Listener interface {
	Listen() error
	RegisterPHook(hook EventHook)
	RegisterNHook(hook WinHook)
	RegisterKeyCmdHook(KeyCmdHook)
	GetEventLoopByID(id int) *FileLoop
}

// Acme implements the Listener interface for acme events
type Acme struct {
	eventHooks  map[AcmeOp][]EventHook
	winHooks    map[AcmeOp][]WinHook
	keyCmdHooks map[rune]*KeyCmdHook
	windows     map[int]string
	eventLoops  map[int]*FileLoop
	debug       bool
	mux         sync.Mutex
}

type EventLoop interface {
	GetWin() *Win
	Start() error
	RegisterPHook(hook EventHook)
	RegisterNHook(hook WinHook)
	RegisterKeyCmdHook(KeyCmdHook)
}

type FileLoop struct {
	ID          int
	File        string
	Win         *Win
	debug       bool
	eventHooks  map[AcmeOp][]EventHook
	winHooks    map[AcmeOp][]WinHook
	keyCmdHooks map[rune]*KeyCmdHook
}

// NewListener constructs an Acme Listener
func NewListener() Listener {
	return &Acme{
		eventHooks:  make(map[AcmeOp][]EventHook),
		winHooks:    make(map[AcmeOp][]WinHook),
		keyCmdHooks: make(map[rune]*KeyCmdHook),
		windows:     make(map[int]string),
		eventLoops:  make(map[int]*FileLoop),
	}
}

// NewEventLoop constructs an event loop
func NewEventLoop(id int, file string) EventLoop {
	return &FileLoop{
		ID:          id,
		File:        file,
		eventHooks:  make(map[AcmeOp][]EventHook),
		winHooks:    make(map[AcmeOp][]WinHook),
		keyCmdHooks: make(map[rune]*KeyCmdHook),
	}
}

// Listen watches the acme event log for events and executes hooks
// based on those events
func (a *Acme) Listen() error {
	if len(os.Getenv("DEBUG")) > 0 {
		a.debug = true
	}

	if a.debug {
		log.Println("opening acme log")
	}
	l, err := acme.Log()
	if err != nil {
		if a.debug {
			log.Printf("failed to read acme log: %v\n", err)
		}
		return err
	}
	for {
		if a.debug {
			log.Println("reading acme event")
		}
		event, err := l.Read()
		if err != nil {
			if a.debug {
				log.Printf("failed to read acme event: %v\n", err)
			}
			return err
		}
		// skip directory windows
		if strings.HasSuffix(event.Name, "/") {
			continue
		}
		// create listener on new window events
		if event.Op == "new" {

			err := a.mapWindows()
			if err != nil {
				if a.debug {
					log.Println("failed to map win IDs")
				}
				log.Println(err)
				continue
			}
			if a.isDisabled(event.ID) {
				continue
			}
			f := &FileLoop{
				ID:          event.ID,
				File:        a.windows[event.ID],
				debug:       a.debug,
				eventHooks:  a.eventHooks,
				winHooks:    a.winHooks,
				keyCmdHooks: a.keyCmdHooks,
			}
			a.eventLoops[event.ID] = f
			go a.startEventLoop(f)
		}
	}
}

func (a *Acme) startEventLoop(f *FileLoop) {
	log.Fatal(f.Start())
}

func (a *Acme) isDisabled(id int) bool {
	filename := a.windows[id]
	// TODO: this should be decerned in a more intelligent way
	disabledNames := []string{"/-", "Del", "xplor"}
	for _, name := range disabledNames {
		if strings.Contains(filename, name) {
			return true
		}
	}
	return false
}

func (a *Acme) mapWindows() error {
	if a.debug {
		log.Println("mapping win IDs to names")
	}
	ws, err := acme.Windows()
	if err != nil {
		return err
	}
	a.mux.Lock()
	defer a.mux.Unlock()
	a.windows = make(map[int]string)
	for _, w := range ws {
		a.windows[w.ID] = w.Name
	}
	return nil
}

func (a *Acme) GetEventLoopByID(id int) *FileLoop {
	return a.eventLoops[id]
}

func (f *FileLoop) GetWin() *Win {
	return f.Win
}

var lastpoint int = 0
func (f *FileLoop) Start() error {
	if f.debug {
		log.Println("opening acme window")
	}
	// open window for modification
	w, err := OpenWin(f.ID, f.File)
	if err != nil {
		if f.debug {
			log.Println("failed to open acme window: %v", err)
		}
		return err
	}
	f.Win = w

	// runs hooks for acme 'new' event
	f.runWinHooks(f.Win)

	for e := range f.Win.OpenEventChan() {
		if f.debug {
			log.Printf("RAW: %+v\n", *e)
		}

		event, err := TokenizeEvent(e, f.ID, f.File)
		if err != nil {
			return err
		}

		if event.Origin == Keyboard {
			lastpoint = event.SelBegin
			event = f.runKeyCmdHooks(event)
		}

		if event.Origin == DelOrigin && event.Type == DelType {
			if f.debug {
				log.Println("delete event received")
			}
			f.Win.WriteEvent(event)
			f.Win.Close()
			return nil
		}

		if f.debug {
			log.Printf("TOKEN: %+v\n", event)
		}

		event = f.runEventHooks(event)
		if f.debug {
			log.Printf("NewEvent: %+v\n", event)
		}
		f.Win.WriteEvent(event)
		
		// TODO: encapsulate this as an optional post save hook
		if event.Builtin == PUT {
			if err := f.Win.SetAddr(fmt.Sprintf("#%d", lastpoint)); err != nil {
				return err
			}
			if err := f.Win.SetTextToAddr(); err != nil {
				return err
			}
			if err := f.Win.ExecShow(); err != nil {
				return err
			}
		}
	}
	return nil
}
