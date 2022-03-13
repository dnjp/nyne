package event

import (
	"log"
	"os"
	"strings"
	"sync"

	"9fans.net/go/acme"
	// "github.com/dnjp/nyne/formatter"
	"github.com/dnjp/nyne/util/io"
)

// Acme implements the Listener interface for acme events
type Acme struct {
	eventHooks  map[AcmeOp][]PutHook
	winHooks    map[AcmeOp][]WinHook
	keyCmdHooks map[rune]*KeyCmdHook
	windows     map[int]string
	bufs        map[int]*Buf
	debug       bool
	mux         sync.Mutex
}

// NewListener constructs an Acme Listener
func NewAcme() *Acme {
	return &Acme{
		eventHooks:  make(map[AcmeOp][]PutHook),
		winHooks:    make(map[AcmeOp][]WinHook),
		keyCmdHooks: make(map[rune]*KeyCmdHook),
		windows:     make(map[int]string),
		bufs:        make(map[int]*Buf),
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
		// TODO
		// ext := formatter.Ext(event.Name, "NONE")
		// if ext == "NONE" || formatter.Conf[ext].Indent == 0 {
		// 	continue
		// }
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
		id:          id,
		file:        a.windows[id],
		debug:       a.debug,
		eventHooks:  a.eventHooks,
		winHooks:    a.winHooks,
		keyCmdHooks: a.keyCmdHooks,
	}
	a.bufs[id] = f
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

// BufListener returns the running Buf by its ID
func (a *Acme) BufListener(id int) *Buf {
	return a.bufs[id]
}

// RegisterPutHook registers hook on acme 'Put' events
func (a *Acme) RegisterPutHook(hook PutHook) {
	hooks := a.eventHooks[Put]
	hooks = append(hooks, hook)
	a.eventHooks[Put] = hooks
}

// RegisterWinHook registers the hook on acme 'New' events
func (a *Acme) RegisterWinHook(hook WinHook) {
	hooks := a.winHooks[New]
	hooks = append(hooks, hook)
	a.winHooks[New] = hooks
}

// RegisterKeyCmdHook registers hook for key events
func (a *Acme) RegisterKeyCmdHook(hook KeyCmdHook) {
	a.keyCmdHooks[hook.Key] = &hook
}
