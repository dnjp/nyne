/*
Shortcuts for working with markdown

	Usage of md:
	  -op string
	    	the operation to perform: link, bold, italic, preview
*/
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/dnjp/nyne"
)

var op = flag.String("op", "", "the operation to perform: link, bold, italic, preview")

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
	err = w.SetAddr("#%d;#%d", nq0, nq1)
	if err != nil {
		panic(err)
	}

	w.SetData(out)

	err = w.SetAddr("#%d", curs)
	if err != nil {
		panic(err)
	}

	err = w.SelectionFromAddr()
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
	dat, err := w.Data(q0, q1)
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
	dat, err := w.Data(q0, q1)
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
	dat, err := w.Data(q0, q1)
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

func preview(w *nyne.Win) {
	outfile := strings.TrimSuffix(path.Base(w.File), ".md")
	dir := path.Dir(w.File)
	outpath := path.Join("/tmp/", outfile+".html")

	var out bytes.Buffer
	tomd := exec.Command(
		"pandoc",
		"--metadata",
		"title="+outfile,
		"-s",
		w.File,
	)
	tomd.Stdout = &out

	err := tomd.Run()
	if err != nil {
		panic(err)
	}

	output := out.Bytes()
	if len(output) == 0 {
		panic("pandoc could not process output")
	}

	// fix relative paths
	output = bytes.ReplaceAll(output, []byte("href=\".."), []byte("href=\""+dir+"/.."))
	output = bytes.ReplaceAll(output, []byte("href=\"./"), []byte("href=\""+dir+"/"))
	output = bytes.ReplaceAll(output, []byte("src=\".."), []byte("src=\""+dir+"/.."))
	output = bytes.ReplaceAll(output, []byte("src=\"./"), []byte("src=\""+dir+"/"))
	err = os.WriteFile(outpath, output, 0644)
	if err != nil {
		panic(err)
	}

	plumb := exec.Command("web", outpath)
	err = plumb.Run()
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()

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
	case "preview":
		preview(w)
	default:
		return
	}
}
