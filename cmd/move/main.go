package main

import (
	"flag"
	"fmt"
	"strings"
	"unicode"

	"github.com/dnjp/nyne"
)

var direction = flag.String("d", "", "the direction to move: up, down, left, right, start, end")
var word = flag.Bool("w", false, "move by word (only valid for left and right)")
var paragraph = flag.Bool("p", false, "move by paragraph (only valid for left and right)")
var sel = flag.Bool("s", false, "select text while moving")

func update(w *nyne.Win, cb func(w *nyne.Win, q0 int) (nq0 int)) {
	q0, q1, err := w.CurrentAddr()
	if err != nil {
		panic(err)
	}

	if *sel {
		var a, b, nq0 int
		switch *direction {
		case "left", "up", "start":
			nq0 = cb(w, q0)
			a = q0
			if nq0 < a {
				a = nq0
			}
			b = nq0
			if q0 > b {
				b = q0
			}
			if q0 != q1 {
				b = q1
			}
		case "right", "down", "end":
			nq0 = cb(w, q1)
			a = q0
			b = nq0
		}
		err = w.SetAddr(fmt.Sprintf("#%d;#%d", a, b))
		if err != nil {
			panic(err)
		}
	} else {
		nq0 := cb(w, q0)
		err = w.SetAddr(fmt.Sprintf("#%d", nq0))
		if err != nil {
			panic(err)
		}
	}

	err = w.SetTextToAddr()
	if err != nil {
		panic(err)
	}
	if err := w.ExecShow(); err != nil {
		panic(err)
	}
}

func readp(w *nyne.Win, q0 int) (nq0 int, c byte) {
	off := 1
	if q0 == 0 {
		off = 0
	}
	addr := fmt.Sprintf("#%d;#%d", q0-off, q0)
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

func readn(w *nyne.Win, q0 int) (nq0 int, c byte, eof bool) {
	err := w.SetAddr(fmt.Sprintf("#%d;#%d", q0, q0+1))
	if err != nil {
		if err.Error() == "address out of range" {
			eof = true
			return
		}
		panic(err)
	}
	dat, err := w.ReadData(q0, q0+1)
	if err != nil {
		panic(err)
	}
	if len(dat) == 0 {
		panic("no data")
	}
	return q0 + 1, dat[0], false
}

func start(w *nyne.Win, q0 int) (nq0, tabs int) {
	var c byte
	nq0 = q0
	for nq0 >= 0 {
		nq0, c = readp(w, nq0)
		if c == '\t' {
			tabs++
		} else if c == '\n' {
			nq0++
			break
		}
	}
	return nq0, tabs
}

func end(w *nyne.Win, q0 int) (nq0, tabs int) {
	var c byte
	nq0 = q0
	for nq0 >= 0 {
		nq0, c, _ = readn(w, nq0)
		if c == '\t' {
			tabs++
		} else if c == '\n' {
			nq0--
			break
		}
	}
	return nq0, tabs
}

func startline(w *nyne.Win, q0 int) (nq0 int) {
	nq0, _ = start(w, q0)
	return nq0
}

func endline(w *nyne.Win, q0 int) (nq0 int) {
	nq0, _ = end(w, q0)
	return nq0
}

func isword(c byte) bool {
	r := rune(c)
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

func left(w *nyne.Win, q0 int) (nq0 int) {
	if *word {
		nq0 = q0 - 1
		var tnq0 int
		var pc, c, nc byte
		for {
			_, pc, _ = readn(w, nq0)
			tnq0, c = readp(w, nq0)
			_, nc = readp(w, nq0-1)
			nq0 = tnq0
			if nq0 == 0 {
				return nq0
			}
			if !isword(pc) && isword(c) && isword(nc) {
				return nq0 + 1
			}
			if !isword(pc) && isword(c) && !isword(nc) {
				return nq0 + 1
			}
			if isword(pc) && !isword(c) {
				return nq0 + 1
			}
		}
	}
	if *paragraph {
		nq0 = q0
		_, ca := readp(w, nq0)
		_, cb, eof := readn(w, nq0)
		if ca == '\n' && (cb == '\n' || eof) {
			nq0--
		}
		for {
			nq0a, ca := readp(w, nq0)
			_, cb, eof := readn(w, nq0)
			if ca == '\n' && cb == '\n' || eof {
				return nq0
			}
			nq0 = nq0a
		}
	}
	if nq0 = q0 - 1; nq0 <= 0 {
		return 0
	}
	return nq0
}

func right(w *nyne.Win, q0 int) (nq0 int) {
	if *word {
		nq0 = q0 + 1
		var tnq0 int
		var pc, c, nc byte
		for {
			_, pc = readp(w, nq0)
			tnq0, c, _ = readn(w, nq0)
			_, nc, _ = readn(w, nq0+1)
			nq0 = tnq0
			if !isword(pc) && isword(c) && isword(nc) {
				return nq0 - 1
			}
			if !isword(pc) && isword(c) && !isword(nc) {
				return nq0 - 1
			}
			if isword(pc) && !isword(c) {
				return nq0 - 1
			}
		}
	}
	if *paragraph {
		nq0 = q0
		_, ca := readp(w, nq0)
		_, cb, _ := readn(w, nq0)
		if ca == '\n' && cb == '\n' {
			nq0++
		}
		for {
			_, ca := readp(w, nq0)
			nq0b, cb, eof := readn(w, nq0)
			if ca == '\n' && cb == '\n' || eof {
				return nq0
			}
			nq0 = nq0b
		}
	}
	return q0 + 1
}

func up(w *nyne.Win, q0 int) (nq0 int) {
	var (
		nl         int  // newline counter, index
		ch, tabs   int  // current line
		chp, tabsp int  // previous line
		c          byte // current character
	)

	ft, _ := nyne.FindFiletype(nyne.Filename(w.File))
	for nq0, c = readp(w, q0); nq0 >= 0; nq0, c = readp(w, nq0) {
		if c == '\n' {
			nl++
		}
		if nq0 == 0 {
			nl++
		}
		switch nl {
		case 0: // current line
			if c == '\t' {
				tabs++
			} else if c != '\n' {
				ch++
			}
		case 1: // previous line
			if c == '\t' {
				tabsp++
			} else if c != '\n' {
				chp++
			}
		case 2: // start of previous line
			if ch == 0 && tabs == 0 && tabsp == 0 {
				// line only contains newline character,
				// so return the current line
				if nq0 > 0 {
					return nq0 + 1
				}
				return nq0
			}
			nq0++
			nc := ch
			if tabs > tabsp {
				nc += (tabs - tabsp) * ft.Tabwidth
				tabs = tabsp
			} else if tabs > 0 && tabsp > tabs {
				nc -= (tabsp - tabs) * ft.Tabwidth
				if nc < 0 {
					nc = 0
				} else {
					tabs = tabsp
				}
			}
			if nc > chp {
				nc = chp
			}
			return nq0 + tabs + nc
		case 3:
			// line only contained newline, move down a
			// line to previous point
			return nq0 + 1
		default:
			break
		}
	}

	// something bad happened, don't move
	return q0
}

func down(w *nyne.Win, q0 int) (nq0 int) {
	var (
		nl, fromstart, tabs int
		atnl                bool
		c                   byte
	)

	ft, _ := nyne.FindFiletype(nyne.Filename(w.File))
	nq0, tabs = start(w, q0) // find beginning of line
	fromstart = q0 - nq0     // find next line with offset
	nq0 = q0                 // move back to starting position

	if fromstart == 0 {
		atnl = true
	}

	tabsn := 0
	flush := false
	flushstart := 0
	var flushc byte
	for {
		nq0, c, _ = readn(w, nq0)
		if c == '\n' {
			nl++
		}

		switch nl {
		case 0: // current line
			continue
		case 1: // next line
			if c == '\t' {
				tabsn++
			}
			if flush {
				continue
			}
			if atnl {
				// starting point was already at a newline
				// so we just need to move down by 1 line
				return nq0
			} else if fromstart <= 0 || tabs-tabsn == 0 {
				flush = true
				flushstart = nq0
				flushc = c
				continue
			}
			fromstart--
		default: // over next line
			if flush {
				var off int
				if tabs-tabsn > 0 {
					off = ((tabs - tabsn) * ft.Tabwidth)
					if fromstart > 0 {
						off -= fromstart
					} else {
						off -= 1 // newline
					}
				} else if tabsn-tabs > 0 {
					if fromstart < ft.Tabwidth {
						fromstart = 0
					}
				}
				if flushc == '\t' && fromstart >= ft.Tabwidth && tabsn > tabs {
					fromstart -= ft.Tabwidth
					off++
				}
				rt := flushstart + off + fromstart
				if rt >= nq0 {
					rt = nq0 - 1
				}
				if q0 == 0 {
					rt--
				}
				return rt
			}
			// backtrack
			return nq0 - 1
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
	case "start":
		update(w, startline)
	case "end":
		update(w, endline)
	}
}
