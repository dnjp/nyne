package main

import (
	"bytes"
	"flag"
	"fmt"
	"strings"

	"github.com/dnjp/nyne"
)

var op = flag.String("op", "", "the operation to perform: link, bold, italic")

func islink(dat []byte) bool {
	if bytes.HasPrefix(dat, []byte("http")) {
		return true
	}
	if bytes.Contains(dat, []byte("/")) {
		return true
	}
	return false
}

func update(w *nyne.Win, cb func(w *nyne.Win, q0, q1 int) (nq0, nq1, curs int, out []byte)) {
	q0, q1, err := w.CurrentAddr()
	if err != nil {
		panic(err)
	}

	nq0, nq1, curs, out := cb(w, q0, q1)
	err = w.SetAddr(fmt.Sprintf("#%d;#%d", nq0, nq1))
	if err != nil {
		panic(err)
	}

	w.SetData(out)

	err = w.SetAddr(fmt.Sprintf("#%d", curs))
	if err != nil {
		panic(err)
	}

	err = w.SetTextToAddr()
	if err != nil {
		panic(err)
	}
}

func link(w *nyne.Win, q0, q1 int) (nq0, nq1, curs int, out []byte) {
	if q0 == q1 {
		out = []byte("[]()")
		curs = q0 + 1
		nq0 = q0
		nq1 = q1
		return
	}
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
		curs = q0 + 1
	} else {
		out = []byte("[")
		out = append(out, dat...)
		out = append(out, "]()"...)
		curs = q0 + len(out) - 1
	}
	nq0 = q0
	nq1 = q1
	return
}

func bold(w *nyne.Win, q0, q1 int) (nq0, nq1, curs int, out []byte) {
	if q0 == q1 {
		out = []byte("**")
		curs = q0 + 1
		nq0 = q0
		nq1 = q1
		return
	}
	dat, err := w.ReadData(q0, q1)
	if err != nil {
		panic(err)
	}
	if dat[len(dat)-1] == '\n' {
		dat = dat[:len(dat)-1]
		q1--
	}
	out = []byte("*")
	out = append(out, dat...)
	out = append(out, "*"...)
	curs = q0 + len(out)
	nq0 = q0
	nq1 = q1
	return
}

func italic(w *nyne.Win, q0, q1 int) (nq0, nq1, curs int, out []byte) {
	if q0 == q1 {
		out = []byte("__")
		curs = q0 + 1
		nq0 = q0
		nq1 = q1
		return
	}
	dat, err := w.ReadData(q0, q1)
	if err != nil {
		panic(err)
	}
	if dat[len(dat)-1] == '\n' {
		dat = dat[:len(dat)-1]
		q1--
	}
	out = []byte("_")
	out = append(out, dat...)
	out = append(out, "_"...)
	curs = q0 + len(out)
	nq0 = q0
	nq1 = q1
	return
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
		update(w, link)
	case "bold":
		update(w, bold)
	case "italic":
		update(w, italic)
	default:
		return
	}
}
