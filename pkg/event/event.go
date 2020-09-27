package event

import (
	"9fans.net/go/acme"
	"fmt"
)

type Event struct {
	ID int
	Op   AcmeOp
	File string
	log  acme.LogEvent
	Win  *Win
}

func (e *Event) ConnectWin() error {
	w, err := acme.Open(e.ID, nil)
	if err != nil {
		return err
	}
	e.Win = &Win{
		handle: w,
	}
	eventlistener(w)
	return nil
}

func eventlistener(w *acme.Win) {
	for e := range w.EventChan() {
		evtType := fmt.Sprintf("%s%s", string(e.C1), string(e.C2))
		switch (evtType) {
		default:
			fmt.Println(evtType)
			fmt.Printf("%+v\n", *e)
			w.WriteEvent(e)
		}
	}
}

func (e *Event) CloseFilesForWin() {
	e.Win.closeFiles()
}