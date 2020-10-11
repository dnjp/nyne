package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"git.sr.ht/~danieljamespost/nyne/pkg/event"
	"git.sr.ht/~danieljamespost/nyne/pkg/formatter"
)

func main() {
	wID, err := strconv.Atoi(os.Getenv("winid"))
	if err != nil {
		log.Print(err)
	}
	tabWidth, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Print(err)
		return
	}
	buf := event.NewBufListener(wID, os.Getenv("$samfile"))
	km := formatter.Keymap{
		GetWinFn: func(id int) (*event.Win, error) {
			if id != wID {
				return nil, fmt.Errorf("id did not match win")
			}
			return buf.GetWin(), nil
		},
		GetIndentFn: func(_ event.Event) int {
			return tabWidth
		},
	}
	buf.RegisterKeyCmdHook(km.Tabexpand(func(evt event.Event) bool {
		return true
	}))
	log.Fatal(buf.Start())
}
