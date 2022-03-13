package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/dnjp/nyne"
)

func main() {
	os.Unsetenv("winid") // do not trust the execution environment

	winid, err := nyne.FindFocusedWinID()
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

	q0, q1, err := w.CurrentAddr()
	if err != nil {
		panic(err)
	}

	dat, err := w.ReadData(q0, q1)
	if err != nil {
		panic(err)
	}

	var in, out bytes.Buffer
	_, err = in.Write(dat)
	if err != nil {
		panic(err)
	}

	com := exec.Command("com")
	com.Stdin = &in
	com.Stdout = &out
	com.Env = os.Environ()
	com.Env = append(com.Env, fmt.Sprintf("winid=%d", winid))
	com.Env = append(com.Env, fmt.Sprintf("%=%s", w.File))
	com.Env = append(com.Env, fmt.Sprintf("samfile=%s", w.File))

	err = com.Run()
	if err != nil {
		panic(err)
	}

	err = w.SetAddr(fmt.Sprintf("#%d;#%d", q0, q1))
	if err != nil {
		panic(err)
	}

	b := out.Bytes()
	w.SetData(b)

	err = w.SetAddr(fmt.Sprintf("#%d;#%d", q0, q0+len(b)))
	if err != nil {
		panic(err)
	}

	err = w.SetTextToAddr()
	if err != nil {
		panic(err)
	}
}
