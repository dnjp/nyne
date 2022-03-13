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

	filename := os.Getenv("samfile")
	if filename == "" {
		filename = os.Getenv("%")
	}
	if filename == "" {
		fmt.Fprintf(os.Stderr, "$samfile and $%% are empty. are you sure you're in acme?")
		os.Exit(1)
	}

	ft, _ := nyne.FindFiletype(nyne.Filename(filename))
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
	key, expand := nyne.Tabexpand(
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
		})

	buf.KeyHooks = map[rune]nyne.Handler{
		key: expand,
	}

	err = buf.Start()
	if err != nil {
		panic(err)
	}
}
