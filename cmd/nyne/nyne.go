package main

import (
	"log"

	"github.com/dnjp/nyne"
	"github.com/dnjp/nyne/format"
)

func main() {
	f, err := format.NewFormatter(nyne.Filetypes, nyne.Menu, nyne.Tag)
	if err != nil {
		log.Fatal(err)
	}
	err = f.Run()
	if err != nil {
		log.Fatal(err)
	}
}
