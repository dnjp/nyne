package main

import (
	"log"

	"github.com/dnjp/nyne"
)

func main() {
	f, err := nyne.NewFormatter(nyne.Filetypes, nyne.Menu, nyne.Tag)
	if err != nil {
		log.Fatal(err)
	}
	err = f.Run()
	if err != nil {
		log.Fatal(err)
	}
}
