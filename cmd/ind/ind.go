package main

import (
	"fmt"
	"github.com/dnjp/nyne/gen"
	"github.com/dnjp/nyne/util/io"
	"os"
	"strconv"
)

func main() {
	filename := gen.GetFileName(os.Getenv("samfile"))
	ext := gen.GetExt(filename, ".txt")
	spec := gen.Conf[ext]
	ts := spec.Indent
	te := spec.Tabexpand
	if ts == 0 {
		nts, err := strconv.Atoi(os.Getenv("tabstop"))
		if err != nil {
			panic(fmt.Errorf("invalid $tabstop: %v", err))
		}
		ts = nts
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
