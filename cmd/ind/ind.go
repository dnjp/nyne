package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/dnjp/nyne"
	"github.com/dnjp/nyne/util/io"
)

func main() {
	var fflag = flag.String("f", "", "file name to operage on")
	var iflag = flag.Int("t", 0, "tabwidth in spaces")
	flag.Parse()
	samfile := os.Getenv("samfile")
	if samfile == "" && fflag != nil {
		samfile = *fflag
	}
	ft, _ := nyne.FindFiletype(nyne.Filename(samfile))
	ts := ft.Tabwidth
	te := ft.Tabexpand
	if ts == 0 {
		tab := os.Getenv("tabstop")
		if tab == "" && iflag != nil {
			ts = *iflag
		} else {
			nts, err := strconv.Atoi(tab)
			if err != nil {
				panic(fmt.Errorf("invalid $tabstop: %v", err))
			}
			ts = nts
		}
	}

	in, err := io.PipeIn()
	if err != nil {
		panic(err)
	}

	io.PipeOut(in, func(line string) string {
		if len(line) == 0 {
			return line
		}
		var tab string
		if te {
			for i := 0; i < ts; i++ {
				tab += " "
			}
			return tab + line
		}
		return "\t" + line
	})
}
