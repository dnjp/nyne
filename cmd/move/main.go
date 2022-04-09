/*
Shortcuts for moving the cursor

	Usage of move:
	  -d string
	    	the direction to move: up, down, left, right, start, end
	  -p	move by paragraph (only valid for left and right)
	  -s	select text while moving
	  -w	move by word (only valid for left and right)
*/
package main

import (
	"flag"
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
		err = w.SetAddr("#%d;#%d", a, b)
		if err != nil {
			panic(err)
		}
	} else {
		nq0 := cb(w, q0)
		err = w.SetAddr("#%d", nq0)
		if err != nil {
			panic(err)
		}
	}

	err = w.SelectionFromAddr()
	if err != nil {
		panic(err)
	}
	if !*sel {
		if err := w.Show(); err != nil {
			panic(err)
		}
	}
}

func start(w *nyne.Win, q0 int) (nq0, tabs int) {
	var c byte
	nq0 = q0
	for nq0 >= 0 {
		nq0--
		c, _ = w.Char(nq0)
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
		c, _ = w.Char(nq0)
		nq0++
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
		var pc, c byte
		for {
			pc, _ = w.Char(nq0)
			nq0--
			c, _ = w.Char(nq0)
			if nq0 == 0 {
				return nq0
			}
			if (!isword(pc) && isword(c)) || (isword(pc) && !isword(c)) {
				return nq0 + 1
			}
		}
	}
	if *paragraph {
		err := w.SetAddr("#%d", q0-1)
		if err != nil {
			panic(err)
		}
		err = w.SetAddr("-/^$/")
		if err != nil {
			panic(err)
		}
		nq0, _, err = w.Addr()
		if err != nil {
			panic(err)
		}
		return nq0
	}
	if nq0 = q0 - 1; nq0 <= 0 {
		return 0
	}
	return nq0
}

func right(w *nyne.Win, q0 int) (nq0 int) {
	if *word {
		nq0 = q0 + 1
		var pc, c, nc byte
		for {
			pc, _ = w.Char(nq0 - 1)
			c, _ = w.Char(nq0)
			nc, _ = w.Char(nq0 + 1)
			nq0++
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
		err := w.SetAddr("#%d", q0+1)
		if err != nil {
			panic(err)
		}
		err = w.SetAddr("+/^$/")
		if err != nil {
			panic(err)
		}
		nq0, _, err = w.Addr()
		if err != nil {
			panic(err)
		}
		return nq0
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

	prev := func(a int) (int, byte) {
		a--
		c, _ := w.Char(a)
		return a, c
	}

	for nq0, c = prev(q0); nq0 >= 0; nq0, c = prev(nq0) {
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
			return q0
		}
	}

	// something bad happened, don't move
	return q0
}

func movedown(body []byte, tw, startQ0, currentQ0 int) (nq0 int) {
	var (
		i, nl, starttabs, off  int
		hasc, hasc2, atq0, set bool
		c                      byte
	)

	fromstart := currentQ0 - startQ0
	for i, c = range body {
		if i == int(fromstart) {
			atq0 = true
		}
		if !set && atq0 && nl == 1 {
			set = true
			tabchars := starttabs * tw
			off = (fromstart - starttabs) + tabchars
		}
		switch nl {
		case 0:
			// only count offset once we are at current
			// position
			if c == '\t' {
				if !hasc && !atq0 {
					starttabs++
				}
			} else {
				hasc = true
			}
		case 1:
			if off <= 0 {
				return startQ0 + i
			}
			if c == '\t' {
				// offset starting tabs
				if !hasc2 {
					starttabs--
					off -= tw
				}
			} else {
				hasc2 = true
				off--
				if starttabs > 0 && off%tw == 0 {
					starttabs--
				}
			}
		case 2:
			return startQ0 + (i - 1)
		}

		if c == '\n' {
			nl++
		}
	}
	return startQ0 + i
}

func textdown(w *nyne.Win, sel bool) (body []byte, startQ0, currentQ0, currentQ1 int, err error) {
	currentQ0, currentQ1, err = w.CurrentAddr()
	if err != nil {
		return
	}

	if sel && currentQ1 > currentQ0 {
		// must set addr to q1 so that the
		// regex below will be in refernece
		// to q1 instead of q0
		err = w.SetAddr("#%d", currentQ1)
		if err != nil {
			return
		}
	}
	err = w.SetAddr("-/^/;+2")
	if err != nil {
		return
	}

	var nq1 int
	startQ0, nq1, err = w.Addr()
	if err != nil {
		return
	}
	body, err = w.Data(startQ0, nq1)
	if err != nil {
		return
	}
	return
}

func tabwidth(w *nyne.Win) int {
	tab, font, err := w.Font()
	if err != nil {
		panic(err)
	}
	cw := font.StringWidth("0")
	return tab / cw
}

func main() {
	flag.Parse()

	winid, err := nyne.FocusedWinID(nyne.FocusedWinAddr())
	if err != nil {
		panic(err)
	}

	w, err := nyne.OpenWin(winid, "")
	if err != nil {
		panic(err)
	}

	switch strings.ToLower(*direction) {
	case "up":
		update(w, up)
	case "down":
		tw := tabwidth(w)
		body, startQ0, currentQ0, currentQ1, err := textdown(w, *sel)
		if err != nil {
			panic(err)
		}
		if *sel {
			nq1 := movedown(body, tw, startQ0, currentQ1)
			err = w.SetAddr("#%d;#%d", currentQ0, nq1)
			if err != nil {
				panic(err)
			}
		} else {
			nq0 := movedown(body, tw, startQ0, currentQ0)
			err = w.SetAddr("#%d", nq0)
			if err != nil {
				panic(err)
			}
		}
		err = w.SelectionFromAddr()
		if err != nil {
			panic(err)
		}
		if !*sel {
			if err := w.Show(); err != nil {
				panic(err)
			}
		}
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
