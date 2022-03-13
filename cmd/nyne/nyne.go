package main

import (
	"log"

	"github.com/dnjp/nyne/format"
)

func main() {
	formatter, err := format.New([]format.Filetype{})
	if err != nil {
		log.Fatal(err)
	}
	formatter.Run()
}
