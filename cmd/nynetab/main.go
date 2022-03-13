package main

import (
	"fmt"
	"os"

	"github.com/dnjp/nyne"
)

func main() {
	os.Unsetenv("winid") // do not trust the execution environment

	winid, err := nyne.FindFocusedWinID()
	if err != nil {
		panic(err)
	}

	wins, err := nyne.Windows()
	if err != nil {
		panic(err)
	}

	w, ok := wins[winid]
	if !ok {
		panic(fmt.Errorf("could not find window with id %d", winid))
	}

	_, _, err = w.CurrentAddr()
	if err != nil {
		panic(err)
	}

	ft, _ := nyne.FindFiletype(nyne.Filename(w.File))
	w.SetData(nyne.Tab(ft.Tabwidth, ft.Tabexpand))
}
