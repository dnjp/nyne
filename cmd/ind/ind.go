package main

import (
	"git.sr.ht/~danieljamespost/nyne/util/io"
	"strconv"
	"fmt"
	"os"
)

func main() {
	ts := 0
	te := false
	if len(os.Getenv("tabexpand")) > 0 {
		te = true
		tss := os.Getenv("tabstop")
		if len(tss) == 0 {
			panic(fmt.Errorf("$tabstop not set"))
		}
		nts, err := strconv.Atoi(tss)	
		if err != nil {
			panic(err)
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
			return tab+line
		}
		return "\t" + line
	})	
}