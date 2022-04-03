/*
Implements tab expansion and indentation.

	Usage of nynetab:
		nynetab [-unindent]

Nynetab is what is used under the hood for tab expansion in nyne.
Executing `nynetab` will insert either a hard or soft tab depending
on (what is
configured) https://github.com/dnjp/nyne/blob/master/config.go .
Executing `nynetab -unindent` will unindent text that is selected.
*/
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/dnjp/nyne"
)

var unindent = flag.Bool("unindent", false, "")

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

	q0, q1, err := w.CurrentAddr()
	if err != nil {
		panic(err)
	}

	// we have selected text to indent/unindent
	if q1 > q0+1 {
		dat, err := w.ReadData(q0, q1)
		if err != nil {
			panic(err)
		}

		var in, out bytes.Buffer
		_, err = in.Write(dat)
		if err != nil {
			panic(err)
		}

		com := exec.Command(func() string {
			if *unindent {
				return "a-"
			}
			return "a+"
		}())
		com.Stdin = &in
		com.Stdout = &out
		com.Env = os.Environ()
		com.Env = append(com.Env, fmt.Sprintf("winid=%d", winid))
		com.Env = append(com.Env, fmt.Sprintf("%%=%s", w.File))
		com.Env = append(com.Env, fmt.Sprintf("samfile=%s", w.File))

		err = com.Run()
		if err != nil {
			panic(err)
		}

		err = w.SetAddr("#%d;#%d", q0, q1)
		if err != nil {
			panic(err)
		}

		b := out.Bytes()
		w.SetData(b)

		err = w.SetAddr("#%d;#%d", q0, q0+len(b))
		if err != nil {
			panic(err)
		}

		err = w.SetTextToAddr()
		if err != nil {
			panic(err)
		}

		return
	}

	ft, _ := nyne.FindFiletype(nyne.Filename(w.File))
	w.SetData(nyne.Tab(ft.Tabwidth, ft.Tabexpand))
}
