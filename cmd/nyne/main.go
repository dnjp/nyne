/*
The core autoformatting engine that is run from within acme
*/
package main

import (
	"log"

	"github.com/dnjp/nyne"
)

func main() {
	f, err := nyne.NewFormatter(nyne.Filetypes, nyne.Menu)
	if err != nil {
		log.Fatal(err)
	}
	err = f.Run()
	if err != nil {
		log.Fatal(err)
	}
}
