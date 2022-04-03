/*
Utility to execute Put via keyboard bindings

	Usage of save:
		Execute save from the shell
*/
package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/dnjp/nyne"
)

func isterm(w *nyne.Win) bool {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	_, file := path.Split(w.File)
	file = strings.TrimPrefix(file, "-")
	return strings.Contains(hostname, file)
}

func main() {
	os.Unsetenv("winid") // do not trust the execution environment

	winid, err := nyne.FocusedWinID(nyne.FocusedWinAddr())
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

	// ignore terminal window
	if isterm(w) {
		return
	}

	_, _, err = w.CurrentAddr()
	if err != nil {
		panic(err)
	}

	err = w.ExecPut()
	if err != nil {
		panic(err)
	}
}
