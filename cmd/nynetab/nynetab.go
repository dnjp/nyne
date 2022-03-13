package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/dnjp/nyne"
	"github.com/dnjp/nyne/event"
	"github.com/dnjp/nyne/format"
)

func main() {
	wid, err := strconv.Atoi(os.Getenv("winid"))
	if err != nil {
		log.Print(err)
	}

	filename := format.Filename(os.Getenv("samfile"))
	ft, _ := nyne.Filetype(filename)
	tabwidth := ft.Tabwidth
	if tabwidth == 0 && len(os.Args) > 1 {
		width, err := strconv.Atoi(os.Args[1])
		if err != nil {
			log.Print(err)
			return
		}
		tabwidth = width
	}

	buf := event.NewBuf(wid, os.Getenv("$samfile"))
	buf.RegisterKeyHook(format.Tabexpand(
		func(evt event.Event) bool {
			return true
		},
		func(id int) (*event.Win, error) {
			if id != wid {
				return nil, fmt.Errorf("id did not match win")
			}
			return buf.Win(), nil
		},
		func(_ event.Event) int {
			return tabwidth
		},
	))
	err = buf.Start()
	if err != nil {
		panic(err)
	}
}
