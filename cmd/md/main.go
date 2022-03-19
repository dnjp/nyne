package main

import (
	"bytes"
	"flag"
	"fmt"
	"strings"

	"github.com/dnjp/nyne"
)

var op = flag.String("op", "", "the operation to perform")

func islink(dat []byte) bool {
	if bytes.HasPrefix(dat, []byte("http")) {
		return true
	}
	if bytes.Contains(dat, []byte("/")) {
		return true
	}
	return false
}

func link(w *nyne.Win) {
	q0, q1, err := w.CurrentAddr()
	if err != nil {
		panic(err)
	}

	var nq0 int
	var out []byte
	if q0 == q1 {
		out = []byte("[]()")
		nq0 = q0 + 1
	} else {
		dat, err := w.ReadData(q0, q1)
		if err != nil {
			panic(err)
		}
		if dat[len(dat)-1] == '\n' {
			dat = dat[:len(dat)-1]
			q1--
		}
		if islink(dat) {
			out = []byte("[](")
			out = append(out, dat...)
			out = append(out, ")"...)
			nq0 = q0 + 1
		} else {
			out = []byte("[")
			out = append(out, dat...)
			out = append(out, "]()"...)
			nq0 = q0 + len(out) - 1
		}
	}

	err = w.SetAddr(fmt.Sprintf("#%d;#%d", q0, q1))
	if err != nil {
		panic(err)
	}

	w.SetData(out)

	err = w.SetAddr(fmt.Sprintf("#%d", nq0))
	if err != nil {
		panic(err)
	}

	err = w.SetTextToAddr()
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()

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

	ft, _ := nyne.FindFiletype(nyne.Filename(w.File))
	if ft.Name != "markdown" {
		return
	}

	switch strings.ToLower(*op) {
	case "link":
		link(w)
		return
	default:
		return
	}
}
