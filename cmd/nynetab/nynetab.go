package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/dnjp/nyne/event"
	"github.com/dnjp/nyne/formatter"
	"github.com/dnjp/nyne/gen"
)

func main() {
	wID, err := strconv.Atoi(os.Getenv("winid"))
	if err != nil {
		log.Print(err)
	}

	filename := gen.GetFileName(os.Getenv("samfile"))
	ext := gen.GetExt(filename, ".txt")
	tabwidth := gen.Conf[ext].Indent
	if tabwidth == 0 && len(os.Args) > 1 {
		width, err := strconv.Atoi(os.Args[1])
		if err != nil {
			log.Print(err)
			return
		}
		tabwidth = width
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
			return tabwidth
		},
	}
	buf.RegisterKeyCmdHook(km.Tabexpand(func(evt event.Event) bool {
		return true
	}))
	err = buf.Start()
	if err != nil {
		panic(err)
	}
}
