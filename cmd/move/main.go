package main

import (
	"flag"
	"fmt"
	"strings"
	"unicode"

	"github.com/dnjp/nyne"
)

var direction = flag.String("d", "", "the direction to move: up, down, left, right")
var word = flag.Bool("w", false, "move by word (only valid for left and right)")

func update(w *nyne.Win, cb func(w *nyne.Win, q0 int) (nq0 int)) {
	q0, _, err := w.CurrentAddr()
	if err != nil {
		panic(err)
	}

	nq0 := cb(w, q0)
	err = w.SetAddr(fmt.Sprintf("#%d", nq0))
	if err != nil {
		panic(err)
	}

	err = w.SetTextToAddr()
	if err != nil {
		panic(err)
	}
}

func readp(w *nyne.Win, q0 int) (nq0 int, c byte) {
	addr := fmt.Sprintf("#%d;#%d", q0-1, q0)
	err := w.SetAddr(addr)
	if err != nil {
		panic(fmt.Errorf("could not set address to '%s': %w", addr, err))
	}
	dat, err := w.ReadData(q0-1, q0)
	if err != nil {
		panic(err)
	}
	if len(dat) == 0 {
		panic("no data")
	}
	return q0 - 1, dat[0]
}

func readn(w *nyne.Win, q0 int) (nq0 int, c byte) {
	err := w.SetAddr(fmt.Sprintf("#%d;#%d", q0, q0+1))
	if err != nil {
		panic(err)
	}
	dat, err := w.ReadData(q0, q0+1)
	if err != nil {
		panic(err)
	}
	if len(dat) == 0 {
		panic("no data")
	}
	return q0 + 1, dat[0]
}

func skip(c byte) bool {
	return c == '\n' || c == ' ' || (!unicode.IsLetter(rune(c)) && !unicode.IsNumber(rune(c)))
}

func left(w *nyne.Win, q0 int) (nq0 int) {
	if *word {
		nq0 = q0
		var c byte
		for {
			nq0, c = readp(w, nq0)
			if skip(c) {
				return nq0
			}
		}
	}
	if nq0 = q0 - 1; nq0 <= 0 {
		return 0
	}
	return nq0
}

func right(w *nyne.Win, q0 int) (nq0 int) {
	if *word {
		nq0 = q0
		var c byte
		for {
			nq0, c = readn(w, nq0)
			if skip(c) {
				return nq0
			}
		}
	}
	return q0 + 1
}

func up(w *nyne.Win, q0 int) (nq0 int) {
	nq0, c := readp(w, q0)
	var nl, cc, fromstart int
	for nq0 > 0 {
		if nl == 1 && c != '\n' {
			cc++
		}
		if c == '\n' {
			nl++
			if nl == 1 {
				fromstart = q0 - nq0
			} else {
				break
			}
		}
		nq0, c = readp(w, nq0)
	}
	if cc < fromstart {
		return nq0 + cc
	}
	return nq0 + fromstart
}

func down(w *nyne.Win, q0 int) (nq0 int) {
	// find beginning of line
	nq0 = q0
	var c byte
	for {
		nq0, c = readp(w, nq0)
		if c == '\n' {
			nq0++
			break
		}
	}

	// find next line with offset
	nq0 = q0
	fromstart := q0 - nq0
	var nl int
	var atnl bool
	if fromstart == 0 {
		atnl = true
	}
	for {
		nq0, c = readn(w, nq0)
		if c == '\n' {
			nl++
		}
		if nl > 1 {
			// we went over the next newline - backtrack
			nq0--
			break
		} else if fromstart <= 0 && !atnl {
			// reached offset
			break
		} else if nl > 0 {
			if atnl {
				// starting point was already at a newline
				// so we just need to move down by 1 line
				break
			}
			fromstart--
		}
	}
	return nq0
}

func main() {
	flag.Parse()

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

	switch strings.ToLower(*direction) {
	case "up":
		update(w, up)
	case "down":
		update(w, down)
	case "left":
		update(w, left)
	case "right":
		update(w, right)
	}
}
