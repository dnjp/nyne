package nyne

import (
	"log"
	"os"
	"strings"
	"sync"

	"9fans.net/go/acme"
	"github.com/dnjp/nyne/util/io"
)

// Acme implements the Listener interface for acme events
type Acme struct {
	eventHooks map[AcmeOp][]Hook
	winHooks   map[AcmeOp][]WinHook
	keyHooks   map[rune]KeyHook
	wins       map[int]string
	bufs       map[int]*Buf
	debug      bool
	mux        sync.Mutex
}

// NewAcme constructs an Acme event listener
func NewAcme() *Acme {
	return &Acme{
		eventHooks: make(map[AcmeOp][]Hook),
		winHooks:   make(map[AcmeOp][]WinHook),
		keyHooks:   make(map[rune]KeyHook),
		wins:       make(map[int]string),
		bufs:       make(map[int]*Buf),
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
			go a.startBuf(event.ID)
		}
	}
}

func (a *Acme) startBuf(id int) {
	err := a.mapWindows()
	if err != nil {
		io.Error(err)
		return
	}
	if a.isDisabled(id) {
		return
	}

	f := &Buf{
		id:         id,
		file:       a.wins[id],
		debug:      a.debug,
		eventHooks: a.eventHooks,
		winHooks:   a.winHooks,
		keyHooks:   a.keyHooks,
	}
	a.bufs[id] = f

	err = f.Start()
	if err != nil {
		io.Error(err)
		return
	}
}

var disabledNames = []string{"/-", "Del", "xplor", "+Errors"}

func (a *Acme) isDisabled(id int) bool {
	filename := a.wins[id]
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
	a.wins = make(map[int]string)
	for _, w := range ws {
		a.wins[w.ID] = w.Name
	}
	return nil
}

// BufListener returns the running Buf by its ID
func (a *Acme) BufListener(id int) *Buf {
	return a.bufs[id]
}

// RegisterHook registers hook on acme 'Put' events
func (a *Acme) RegisterHook(hook Hook) {
	a.eventHooks[hook.Op] = append(a.eventHooks[hook.Op], hook)
}

// RegisterWinHook registers the hook on acme 'New' events
func (a *Acme) RegisterWinHook(hook WinHook) {
	a.winHooks[hook.Op] = append(a.winHooks[hook.Op], hook)
}

// RegisterKeyHook registers hook for key events
func (a *Acme) RegisterKeyHook(hook KeyHook) {
	a.keyHooks[hook.Key] = hook
}
