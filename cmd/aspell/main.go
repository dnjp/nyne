package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/dnjp/nyne"
)

func main() {
	var winid int
	var err error
	if wid := os.Getenv("winid"); wid != "" {
		winid, err = strconv.Atoi(wid)
	} else {
		winid, err = nyne.FindFocusedWinID()
	}
	if err != nil {
		panic(err)
	}

	wins, err := nyne.Windows()
	if err != nil {
		panic(err)
	}

	w, ok := wins[winid]
	if !ok {
		panic(fmt.Errorf("could not find window with id %d", winid))
	}

	dir := filepath.Dir(w.File)
	name := dir + "/-spell"
	w2, err := nyne.NewWin()
	if err != nil {
		panic(err)
	}
	err = w2.Name(name)
	if err != nil {
		panic(err)
	}

	body, err := w.ReadBody()
	if err != nil {
		panic(err)
	}

	env := envvars(winid, name)
	var spellflags []string
	if len(os.Args) > 1 {
		spellflags = os.Args[1:]
	}

	var lastnl int
	var start int
	for i, c := range body {
		if c != '\n' {
			continue
		}
		line := body[lastnl : i+1]
		out, err := pipe(bytes.NewBuffer(line), env, "spell", spellflags...)
		if err != nil {
			panic(err)
		}
		words := bufio.NewReader(out)
		for {
			word, _, err := words.ReadLine()
			if err != nil {
				break
			}
			if len(word) == 0 {
				continue
			}
			idx := bytes.Index(line, word)
			addr := w.File + fmt.Sprintf(":#%d:%s\n", start+idx, string(word))
			err = w2.WriteToBody([]byte(addr))
			if err != nil {
				panic(err)
			}
		}
		lastnl = i + 1
		start = lastnl
	}

	err = w2.MarkWinClean()
	if err != nil {
		panic(err)
	}
	err = w2.SetAddr("#0")
	if err != nil {
		panic(err)
	}
	err = w2.SetTextToAddr()
	if err != nil {
		panic(err)
	}
	if err := w2.ExecShow(); err != nil {
		panic(err)
	}
}

func envvars(winid int, filename string) []string {
	env := os.Environ()
	env = append(env, fmt.Sprintf("winid=%d", winid))
	env = append(env, fmt.Sprintf("%%=%s", filename))
	env = append(env, fmt.Sprintf("samfile=%s", filename))
	return env
}

func pipe(in *bytes.Buffer, env []string, x string, args ...string) (*bytes.Buffer, error) {
	// b := in.Bytes()
	// if b[len(b)-1] != '\n' {
	// 	panic(fmt.Sprintf("not newline: '%s'", string(b)))
	// }
	var out bytes.Buffer
	cmd := exec.Command(x, args...)
	if in != nil {
		cmd.Stdin = in
	}
	cmd.Stdout = &out
	cmd.Env = env
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return &out, nil
}
