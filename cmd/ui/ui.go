package main

import (
	"flag"
	"fmt"
	"github.com/dnjp/nyne/gen"
	"github.com/dnjp/nyne/util/io"
	"os"
	"strings"
	"strconv"
)

func main() {
	var fflag = flag.String("f", "", "file name to operage on")
	var iflag = flag.Int("t", 0, "tabwidth in spaces")
	flag.Parse()
	samfile := os.Getenv("samfile")
	if samfile == "" && fflag != nil {
		samfile = *fflag
	}
	filename := gen.GetFileName(samfile)
	ext := gen.GetExt(filename, ".txt")
	spec := gen.Conf[ext]
	ts := spec.Indent
	te := spec.Tabexpand
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
			return strings.Replace(line, tab, "", 1)
		}
		return strings.Replace(line, "\t", "", 1)
	})
}
