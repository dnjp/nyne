package main

import (
	"log"
	"fmt"
	"os"
	"strconv"

	"git.sr.ht/~danieljamespost/nyne/pkg/formatter"
	"git.sr.ht/~danieljamespost/nyne/pkg/event"
)

func main() {
	debug := len(os.Getenv("DEBUG")) > 0
	wId, err := strconv.Atoi(os.Getenv("winid"))
	if err != nil {
		log.Print(err)
	}	
	tabWidth, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Print(err)
		return
	}
	loop := event.NewEventLoop(wId, os.Getenv("$samfile"))
	km := formatter.Keymap{
		GetWinFn: func(id int) (*event.Win, error) {
			if id != wId {
				return nil, fmt.Errorf("id did not match win")
			}
			return loop.GetWin(), nil
		},
		GetIndentFn: func(_ event.Event) int {
			return tabWidth
		},
	}
	loop.RegisterKeyCmdHook(km.Tabexpand(func(evt event.Event) bool {
		return true		
	}))	
	log.Fatal(loop.Start())
}