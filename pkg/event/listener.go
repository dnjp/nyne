package event

import (
	"fmt"
	"os"
	"strings"
	"log"

	"9fans.net/go/acme"
)

type Listener interface {
	Listen() error
	RegisterHook(hook Hook)
	RegisterOpenHook(hook OpenHook)
}
type Acme struct {
	hooks     map[AcmeOp][]Hook
	openHooks map[AcmeOp][]OpenHook
	windows   map[int]string
	debug     bool
}

func NewListener() Listener {
	return &Acme{
		hooks:     make(map[AcmeOp][]Hook),
		openHooks: make(map[AcmeOp][]OpenHook),
		windows:   make(map[int]string),
	}
}

func (a *Acme) Listen() error {
	if len(os.Getenv("DEBUG")) > 0 {
		a.debug = true
	}

	l, err := acme.Log()
	if err != nil {
		return err
	}
	for {
		event, err := l.Read()
		if err != nil {
			return err
		}

		// create listener on new window events
		if event.Op == "new" {
			err := a.mapWindows()
			if err != nil {
				fmt.Println(err)
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
	ws, err := acme.Windows()
	if err != nil {
		return err
	}
	for _, w := range ws {
		a.windows[w.ID] = w.Name
	}
	return nil
}

func (a *Acme) startEventListener(id int) {
	if a.debug {
		fmt.Println("starting event listener")
	}
	w, err := acme.Open(id, nil)
	if err != nil {
		fmt.Println(err)
	}

	a.runOpenHooks(&Win{
		File:   a.windows[id],
		ID:     id,
		handle: w,
	})

	if w == nil {
		log.Printf("lost window handler")
		return
	}
	for e := range w.EventChan() {
		if a.debug {
			fmt.Printf("RAW: %+v\n", *e)
		}

		// empty event received on delete
		if e.C1 == 0 && e.C2 == 0 {
			w.CloseFiles()
			go a.mapWindows()
			break
		}

		event, err := a.tokenizeEvent(w, e, id)
		if err != nil {
			w.WriteEvent(e)
			if a.debug {
				fmt.Println(err)
			}
			return
		}

		if a.debug {
			fmt.Printf("TOKEN: %+v\n", *event)
			fmt.Printf("\n")
		}

		newEvent := a.runEventHooks(event)
		w.WriteEvent(newEvent.raw)
	}
}
