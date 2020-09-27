package event

import (
	"fmt"

	"9fans.net/go/acme"
)

type Listener interface {
	Listen() error
	RegisterHook(hook Hook)
	RegisterOpenHook(hook OpenHook)
}
type Acme struct {
	hooks map[AcmeOp][]Hook
	openHooks map[AcmeOp][]OpenHook
	windows map[int]string
}

func NewListener() Listener {
	return &Acme{
		hooks: make(map[AcmeOp][]Hook),
		openHooks: make(map[AcmeOp][]OpenHook),
		windows: make(map[int]string),
	}
}

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

		// create listener on new window events
		if event.Op == "new" {
			err := a.mapWindows()
			if err != nil {
				fmt.Println(err)
			}
			a.startEventListener(event.ID)
		}
	}
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
	fmt.Println("starting event listener")
	w, err := acme.Open(id, nil)
	if err != nil {
		fmt.Println(err)
	}

	a.runOpenHooks(&Win{
		File: a.windows[id],
		ID: id,
		handle: w,
	})

	for e := range w.EventChan() {
		fmt.Printf("RAW: %+v\n", *e)

		// empty event received on delete
		if e.C1 == 0 && e.C2 == 0 {
			w.CloseFiles()
			break
		}

		event, err := a.tokenizeEvent(w, e, id)
		if err != nil {
			w.WriteEvent(e)
			fmt.Println(err)
			return
		}

		fmt.Printf("TOKEN: %+v\n", *event)
		fmt.Printf("\n")

		newEvent := a.runEventHooks(event)
		w.WriteEvent(newEvent.raw)
	}
}