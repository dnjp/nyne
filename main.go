package main

import (
	"fmt"
	"log"
	"os/user"
  "os"

	"git.sr.ht/~danieljamespost/nyne/pkg/formatter"
	"git.sr.ht/~danieljamespost/nyne/util/config"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	cfgPath := fmt.Sprintf("%s/lib/nyne", usr.HomeDir)
	npath := os.Getenv("NYNERULES")
	if npath != "" {
	  cfgPath = npath
	}
	conf, err := config.Load(cfgPath)
	if err != nil {
		log.Fatal(err)
	}
	f := formatter.New(conf)
	f.Listen()
}
