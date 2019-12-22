package main

import (
	"fmt"
	"log"
	"os/user"

	"git.sr.ht/~danieljamespost/nyne/pkg/nyne"
	"git.sr.ht/~danieljamespost/nyne/pkg/util/config"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	cfgPath := fmt.Sprintf("%s/.nyne", usr.HomeDir)
	conf, err := config.Load(cfgPath)
	if err != nil {
		log.Fatal(err)
	}
	nyne.New(conf)
}
