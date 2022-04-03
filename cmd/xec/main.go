package main

import (
	"fmt"
	"os"

	"github.com/dnjp/nyne"
)

func usage(base string) {
	fmt.Fprintf(os.Stderr, "%s [command]", base)
}

func main() {
	if len(os.Args) <= 1 {
		usage(os.Args[0])
		os.Exit(1)
	}
	cmd := os.Args[1]
	args := os.Args[2:]
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

	err = w.ExecInTag(cmd, args...)
	if err != nil {
		panic(err)
	}
}
