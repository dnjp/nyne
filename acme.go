package nyne

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"9fans.net/go/acme"
)

// Acme implements the Listener interface for acme events
type Acme struct {
	EventHooks map[Text][]Handler
	WinHooks   map[Text][]WinHandler
	KeyHooks   map[rune]Handler
	wins       map[int]string
	bufs       map[int]*Buf
	mux        sync.Mutex
}

// NewAcme constructs an Acme event listener
func NewAcme() *Acme {
	return &Acme{
		EventHooks: make(map[Text][]Handler),
		WinHooks:   make(map[Text][]WinHandler),
		KeyHooks:   make(map[rune]Handler),
		wins:       make(map[int]string),
		bufs:       make(map[int]*Buf),
	}
}

// Listen watches the acme event log for events and executes hooks
// based on those events
func (a *Acme) Listen() error {
	l, err := acme.Log()
	if err != nil {
		return err
	}
	for {
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
		fmt.Fprintf(os.Stderr, "%+v", err)
		return
	}
	if a.isDisabled(id) {
		return
	}

	f := &Buf{
		id:         id,
		file:       a.wins[id],
		EventHooks: a.EventHooks,
		WinHooks:   a.WinHooks,
		KeyHooks:   a.KeyHooks,
	}
	a.bufs[id] = f

	err = f.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
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

// Buf returns the running Buf by its ID
func (a *Acme) Buf(id int) *Buf {
	return a.bufs[id]
}
