package main

import (
	"log"
	"os"
	"strconv"

	"git.sr.ht/~danieljamespost/nyne/pkg/event"
	"9fans.net/go/acme"
)

func main() {
	wId, err := strconv.Atoi(os.Getenv("winid"))
	if err != nil {
		log.Print(err)
	}
	w, err := acme.Open(wId, nil)
	if err != nil {
		log.Print(err)
	}
	tabWidth, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Print(err)
	}
	lis := event.NewListener()
	lis.SetTabexpand(event.NewWin(w), tabWidth)
}