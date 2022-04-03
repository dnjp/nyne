/*
Indent selected source code
*/
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/dnjp/nyne"
)

func main() {
	filename := os.Getenv("samfile")
	if filename == "" {
		filename = os.Getenv("%")
	}
	if filename == "" {
		fmt.Fprintf(os.Stderr, "$samfile and $%% are empty. are you sure you're in acme?")
		os.Exit(1)
	}

	ft, _ := nyne.FindFiletype(nyne.Filename(filename))
	tw := ft.Tabwidth
	te := ft.Tabexpand
	if tw == 0 {
		tab := os.Getenv("tabstop")
		if tab == "" {
			tw = 8
		} else {
			ntw, err := strconv.Atoi(tab)
			if err != nil {
				panic(fmt.Errorf("invalid $tabstop: %v", err))
			}
			tw = ntw
		}
	}

	tab := nyne.Tab(tw, te)
	var i, lastnl int
	lastnl = -1
	reader := bufio.NewReader(os.Stdin)
	for {
		b, err := reader.ReadByte()
		if err != nil && err == io.EOF {
			break
		}
		if b == '\n' {
			lastnl = i
		}
		if i == 0 || (i > 0 && i-1 == lastnl) {
			for _, c := range tab {
				fmt.Fprintf(os.Stdout, "%c", c)
			}
		}
		fmt.Fprintf(os.Stdout, "%c", b)
		i++
	}
}
