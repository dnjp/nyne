/*
Wrapper around f+ or f- intended to be invoked from a tool like skhd
*/
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/dnjp/nyne"
)

var op = flag.String("op", "inc", "font operation to execute: inc, dec")

func main() {
	flag.Parse()

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

	var cmd string
	if strings.ToLower(*op) == "dec" {
		cmd = "f-"
	} else {
		cmd = "f+"
	}
	font := exec.Command(cmd)
	font.Env = os.Environ()
	font.Env = append(font.Env, fmt.Sprintf("winid=%d", winid))
	font.Env = append(font.Env, fmt.Sprintf("%%=%s", w.File))
	font.Env = append(font.Env, fmt.Sprintf("samfile=%s", w.File))
	err = font.Run()
	if err != nil {
		panic(err)
	}
}
