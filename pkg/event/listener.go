package event

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"9fans.net/go/acme"
	"git.sr.ht/~danieljamespost/nyne/util/io"
)

// Listener can listen for acme Event and Window hooks
type Listener interface {
	Listen() error
	RegisterPutHook(hook PutHook)
	RegisterWinHook(hook WinHook)
	RegisterKeyCmdHook(KeyCmdHook)
	GetBufListener(id int) *Buf
}

// Acme implements the Listener interface for acme events
type Acme struct {
	eventHooks        map[AcmeOp][]PutHook
	winHooks          map[AcmeOp][]WinHook
	keyCmdHooks       map[rune]*KeyCmdHook
	windows           map[int]string
	eventBufListeners map[int]*Buf
	debug             bool
	mux               sync.Mutex
}

// BufListener processes hooks on acme events
type BufListener interface {
	GetWin() *Win
	Start() error
	RegisterPutHook(hook PutHook)
	RegisterWinHook(hook WinHook)
	RegisterKeyCmdHook(KeyCmdHook)
}

// Buf implements the BufListener interface and runs on opened
// acme buffers
type Buf struct {
	ID          int
	File        string
	Win         *Win
	debug       bool
	eventHooks  map[AcmeOp][]PutHook
	winHooks    map[AcmeOp][]WinHook
	keyCmdHooks map[rune]*KeyCmdHook
}

// NewListener constructs an Acme Listener
func NewListener() Listener {
	return &Acme{
		eventHooks:        make(map[AcmeOp][]PutHook),
		winHooks:          make(map[AcmeOp][]WinHook),
		keyCmdHooks:       make(map[rune]*KeyCmdHook),
		windows:           make(map[int]string),
		eventBufListeners: make(map[int]*Buf),
	}
}

// NewBufListener constructs an event loop
func NewBufListener(id int, file string) BufListener {
	return &Buf{
		ID:          id,
		File:        file,
		eventHooks:  make(map[AcmeOp][]PutHook),
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
			return err
		}
		// skip directory windows
		if strings.HasSuffix(event.Name, "/") {
			continue
		}
		// create listener on new window events
		if event.Op == "new" {
			go a.handleNewOp(event.ID)
		}
	}
}

func (a *Acme) handleNewOp(id int) {
	err := a.mapWindows()
	if err != nil {
		io.Error(err)
		return
	}
	if a.isDisabled(id) {
		return
	}
	f := &Buf{
		ID:          id,
		File:        a.windows[id],
		debug:       a.debug,
		eventHooks:  a.eventHooks,
		winHooks:    a.winHooks,
		keyCmdHooks: a.keyCmdHooks,
	}
	a.eventBufListeners[id] = f
	err = f.Start()
	if err != nil {
		io.Error(err)
		return
	}
}

func (a *Acme) isDisabled(id int) bool {
	filename := a.windows[id]
	disabledNames := []string{"/-", "Del", "xplor", "+Errors"}
	for _, name := range disabledNames {
		if strings.Contains(filename, name) {
			return true
		}
	}
	return false
}

func (a *Acme) mapWindows() error {
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

// GetBufListener returns the running Buf by its ID
func (a *Acme) GetBufListener(id int) *Buf {
	return a.eventBufListeners[id]
}

// GetWin returns the active acme Window
func (b *Buf) GetWin() *Win {
	return b.Win
}

// Start begins the event listener for the window
func (b *Buf) Start() error {
	w, err := OpenWin(b.ID, b.File)
	if err != nil {
		return err
	}
	b.Win = w

	// runs hooks for acme 'new' event
	b.runWinHooks(b.Win)

	for e := range b.Win.OpenEventChan() {
		if b.debug {
			log.Printf("RAW: %+v\n", *e)
		}

		event, err := TokenizeEvent(e, b.ID, b.File)
		if err != nil {
			return err
		}

		if event.Origin == Keyboard {
			w.Lastpoint = event.SelBegin
			event = b.runKeyCmdHooks(event)
		} else {
			if event.Origin == DelOrigin && event.Type == DelType {
				b.Win.WriteEvent(event)
				b.Win.Close()
				return nil
			}
			event = b.runPutHooks(event)
		}

		if b.debug {
			log.Printf("TOKEN: %+v\n", event)
		}

		b.Win.WriteEvent(event)

		for _, h := range event.PostHooks {
			if err := h(event); err != nil {
				return err
			}
		}

		// maintain current address after formatting buffer
		if event.Builtin == PUT {
			body, err := b.Win.ReadBody()
			if err != nil {
				return err
			}
			if len(body) < w.Lastpoint {
				w.Lastpoint = len(body)
			}
			if err := b.Win.SetAddr(fmt.Sprintf("#%d", w.Lastpoint)); err != nil {
				return err
			}
			if err := b.Win.SetTextToAddr(); err != nil {
				return err
			}
			if err := b.Win.ExecShow(); err != nil {
				return err
			}
		}
	}
	return nil
}
