package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/dnjp/nyne"
)

func main() {
	wid, err := strconv.Atoi(os.Getenv("winid"))
	if err != nil {
		log.Print(err)
	}

	ft, _ := nyne.FindFiletype(nyne.Filename(os.Getenv("samfile")))
	tabwidth := ft.Tabwidth
	if tabwidth == 0 && len(os.Args) > 1 {
		width, err := strconv.Atoi(os.Args[1])
		if err != nil {
			log.Print(err)
			return
		}
		tabwidth = width
	}

	buf := nyne.NewBuf(wid, os.Getenv("$samfile"))
	buf.RegisterKeyHook(nyne.Tabexpand(
		func(evt nyne.Event) bool {
			return true
		},
		func(id int) (*nyne.Win, error) {
			if id != wid {
				return nil, fmt.Errorf("id did not match win")
			}
			return buf.Win(), nil
		},
		func(_ nyne.Event) int {
			return tabwidth
		},
	))
	err = buf.Start()
	if err != nil {
		panic(err)
	}
}
