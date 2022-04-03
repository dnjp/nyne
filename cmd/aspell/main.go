package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/dnjp/nyne"
)

func envvars(winid int, filename string) []string {
	env := os.Environ()
	env = append(env, fmt.Sprintf("winid=%d", winid))
	env = append(env, fmt.Sprintf("%%=%s", filename))
	env = append(env, fmt.Sprintf("samfile=%s", filename))
	return env
}

func pipe(in *bytes.Buffer, env []string, x string, args ...string) (*bytes.Buffer, error) {
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

func text(w *nyne.Win) (body []byte, q0 int) {
	var q1 int
	var err error
	q0, q1, err = w.CurrentAddr()
	if err != nil {
		panic(err)
	}

	if q1 > q0 {
		body, err = w.ReadData(q0, q1)
		if err != nil {
			panic(err)
		}
	} else {
		q0 = 0
		body, err = w.ReadBody()
		if err != nil {
			panic(err)
		}
	}
	return
}

func winid() (winid int, err error) {
	if wid := os.Getenv("winid"); wid != "" {
		winid, err = strconv.Atoi(wid)
	} else {
		winid, err = nyne.FindFocusedWinID()
	}
	return
}

func inout() (inw, outw *nyne.Win, wid int, name string, err error) {
	var wins map[int]*nyne.Win
	var ok bool

	wid, err = winid()
	if err != nil {
		return
	}

	wins, err = nyne.Windows()
	if err != nil {
		return
	}

	inw, ok = wins[wid]
	if !ok {
		err = fmt.Errorf("could not find window with id %d", wid)
		return
	}

	dir := filepath.Dir(inw.File)
	name = dir + "/-spell"

	for _, win := range wins {
		if win.File == name {
			outw = win
			err = outw.ClearBody()
			if err != nil {
				return
			}
		}
	}
	if outw == nil {
		outw, err = nyne.NewWin()
		if err != nil {
			return
		}

		err = outw.Name(name)
		if err != nil {
			return
		}
	}
	return
}

func main() {
	var (
		lastnl, q0      int
		err             error
		out             string
		corrections     map[int]string
		body, line      []byte
		addrs           []int
		env, spellflags []string
	)

	corrections = make(map[int]string)
	addrs = make([]int, 0)

	inw, outw, wid, name, err := inout()
	if err != nil {
		panic(err)
	}

	body, q0 = text(inw)
	env = envvars(wid, name)
	if len(os.Args) > 1 {
		spellflags = os.Args[1:]
	}

	for i, c := range body {
		if c != '\n' {
			continue
		}
		line = body[lastnl : i+1]
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
			idx := bytes.Index(line, word) + q0
			addr := inw.File + fmt.Sprintf(":#%d:%s", idx, string(word))
			corrections[idx] = addr
			addrs = append(addrs, idx)
		}
		lastnl = i + 1
		q0 += len(line)
	}

	sort.Ints(addrs)
	for _, addr := range addrs {
		out += corrections[addr] + "\n"
	}

	err = outw.WriteToBody([]byte(out))
	if err != nil {
		panic(err)
	}

	err = outw.MarkWinClean()
	if err != nil {
		panic(err)
	}
	err = outw.SetAddr("#0")
	if err != nil {
		panic(err)
	}
	err = outw.SetTextToAddr()
	if err != nil {
		panic(err)
	}
	if err := outw.ExecShow(); err != nil {
		panic(err)
	}
}
