package main

import (
	"fmt"
	"os"
	"os/user"

	"git.sr.ht/~danieljamespost/nyne/pkg/formatter"
	"git.sr.ht/~danieljamespost/nyne/util/config"
	"git.sr.ht/~danieljamespost/nyne/util/io"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		io.Error(err)
		return
	}
	cfgPath := fmt.Sprintf("%s/.config/nyne/nyne.toml", usr.HomeDir)
	npath := os.Getenv("NYNERULES")
	if len(npath) > 0 {
		cfgPath = npath
	}
	conf, err := config.Load(cfgPath)
	if err != nil {
		io.Error(err)
		return
	}
	f := formatter.New(conf)
	f.Run()
}
