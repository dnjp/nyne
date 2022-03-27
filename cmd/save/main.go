package main

import (
	"fmt"
	"os"
	"path"
	"strings"

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

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	// ignore save for terminal
	_, file := path.Split(w.File)
	file = strings.TrimPrefix(file, "-")
	if strings.Contains(hostname, file) {
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
