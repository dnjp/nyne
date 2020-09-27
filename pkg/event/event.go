package event

import (
	"9fans.net/go/acme"
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
	return nil
}

func (e *Event) CloseFilesForWin() {
	e.Win.closeFiles()
}