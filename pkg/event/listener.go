package event

import (
	"log"
	"fmt"
	"os"
	"strings"
	"sync"
	"unicode/utf8"
	
	"9fans.net/go/acme"
)

// Listener can listen for acme Event and Window hooks
type Listener interface {
	Listen() error
	RegisterPHook(hook EventHook)
	RegisterNHook(hook WinHook)
	SetTabexpand(w *Win, width int)	
}

// Acme implements the Listener interface for acme events
type Acme struct {
	eventHooks map[AcmeOp][]EventHook
	winHooks   map[AcmeOp][]WinHook
	windows    map[int]string
	debug      bool
	mux        sync.Mutex
}

// NewListener constructs an Acme Listener
func NewListener() Listener {
	return &Acme{
		eventHooks: make(map[AcmeOp][]EventHook),
		winHooks:   make(map[AcmeOp][]WinHook),
		windows:    make(map[int]string),
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
			a.startEventListener(event.ID)
		}
	}
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

func (a *Acme) startEventListener(id int) {
	if a.debug {
		log.Println("opening acme window")
	}
	// open window for modification
	w, err := acme.Open(id, nil)
	if err != nil {
		if a.debug {
			log.Println("failed to open acme window")
		}
		log.Println(err)
		return
	}

	// runs hooks for acme 'new' event
	a.runWinHooks(&Win{
		File:   a.windows[id],
		ID:     id,
		handle: w,
	})

	if w == nil {
		if a.debug {
			log.Printf("lost window handle")
		}
		return
	}
	for e := range w.EventChan() {
		if a.debug {
			log.Printf("RAW: %+v\n", *e)
		}

		// empty event received on delete
		if e.C1 == 0 && e.C2 == 0 {
			if a.debug {
				log.Println("received empty event: treating as del")
			}
			w.CloseFiles()
			go a.mapWindows()
			w.WriteEvent(e)
			break
		}

		event, err := a.tokenizeEvent(w, e, id)
		if err != nil {
			w.WriteEvent(e)
			if a.debug {
				log.Println(err)
			}
			w.WriteEvent(e)
			return
		}

		if a.debug {
			log.Printf("TOKEN: %+v\n", *event)
			log.Printf("\n")
		}

		newEvent := a.runEventHooks(event)
		w.WriteEvent(newEvent.raw)
	}
}

func (a *Acme) SetTabexpand(w *Win, width int) {
	var tab []byte
	for i := 0; i < width; i++ {
		tab = append(tab, ' ')
	}

	for e := range w.handle.EventChan() {
		evtType := fmt.Sprintf("%s%s", string(e.C1), string(e.C2))
		switch (evtType) {
		case "KI":
    			if string(e.Text) == "	" {
				err := w.handle.Addr("#%d;+#1", e.Q0)
				if err != nil {
					log.Print(err)
				}
				w.handle.Write("data", tab)

				e.C1 = 70
				e.C2 = 73
				e.Q1 = e.Q0 + utf8.RuneCount(tab)
				e.OrigQ1 = e.Q0 + utf8.RuneCount(tab)
				e.Nr = utf8.RuneCount(tab)
				e.Text = tab
				w.handle.WriteEvent(e)
			}
		default:
			w.handle.WriteEvent(e)
		}
	}
}