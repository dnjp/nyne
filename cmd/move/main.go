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

func up(body []byte, tw, start, q int) (nq0 int) {
	var (
		i, nl, fromstart, starttabs, off int
		c                                byte
	)

	// find fromstart
	for i = len(body) - 1; i >= 0; i-- {
		c = body[i]
		if c == '\n' {
			nl++
		}
		if nl == 1 {
			fromstart = (len(body) - 1) - i
			break
		}
		if c == '\t' {
			starttabs++
		}
	}

	if fromstart == 0 {
		return q - len(body)
	}

	tabchars := starttabs * tw
	off = (fromstart - starttabs) + tabchars
	nl = 0
	for i, c = range body {
		if c == '\n' {
			nl++
		}
		if nl == 1 {
			return q - (len(body) - i)
		}
		if c == '\t' {
			starttabs--
			off -= tw
		} else {
			off--
			// subtract character from tab offset
			if starttabs > 0 && off%tw == 0 {
				starttabs--
			}
		}
		if off <= 0 {
			return q - ((len(body) - 1) - i)
		}
	}

	// something bad happened, don't move
	return q
}

func down(body []byte, tw, start, q int) (nq0 int) {
	var (
		i, nl, starttabs, off int
		hasc, hasc2, atq, set bool
		c                     byte
	)

	fromstart := q - start
	for i, c = range body {
		if i == int(fromstart) {
			atq = true
		}
		if !set && atq && nl == 1 {
			set = true
			tabchars := starttabs * tw
			off = (fromstart - starttabs) + tabchars
		}
		switch nl {
		case 0:
			// only count offset once we are at current
			// position
			if c == '\t' {
				if !hasc && !atq {
					starttabs++
				}
			} else {
				hasc = true
			}
		case 1:
			if off <= 0 {
				return start + i
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
				// subtract character from tab offset
				if starttabs > 0 && off%tw == 0 {
					starttabs--
				}
			}
		case 2:
			// we are moving from the end of the current
			// line, but the next line has fewer characters
			// than the line we started on. end on the
			// newline for the next line.
			return start + (i - 1)
		}

		if c == '\n' {
			nl++
		}
	}
	return start + i
}

func uptext(w *nyne.Win, sel bool) (body []byte, start, q0, q1 int, err error) {
	q0, q1, err = w.CurrentAddr()
	if err != nil {
		return
	}

	err = w.SetAddr("-1;#%d", q0)
	if err != nil {
		return
	}

	var end int
	start, end, err = w.Addr()
	if err != nil {
		return
	}
	body, err = w.Data(start, end)
	if err != nil {
		return
	}
	return
}

func downtext(w *nyne.Win, sel bool) (body []byte, start, q0, q1 int, err error) {
	q0, q1, err = w.CurrentAddr()
	if err != nil {
		return
	}

	if sel && q1 > q0 {
		// must set addr to q1 so that the
		// regex below will be in refernece
		// to q1 instead of q0
		err = w.SetAddr("#%d", q1)
		if err != nil {
			return
		}
	}
	err = w.SetAddr("-/^/;+2")
	if err != nil {
		return
	}

	var end int
	start, end, err = w.Addr()
	if err != nil {
		return
	}
	body, err = w.Data(start, end)
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

func updatesel(w *nyne.Win, sel bool, q0, q1 int) {
	var err error
	if sel {
		err = w.SetAddr("#%d;#%d", q0, q1)
		if err != nil {
			panic(err)
		}
	} else {
		err = w.SetAddr("#%d", q0)
		if err != nil {
			panic(err)
		}
	}
	err = w.SelectionFromAddr()
	if err != nil {
		panic(err)
	}
	if !sel {
		if err := w.Show(); err != nil {
			panic(err)
		}
	}
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

	tw := tabwidth(w)
	switch strings.ToLower(*direction) {
	case "up":
		body, start, q0, q1, err := uptext(w, *sel)
		if err != nil {
			panic(err)
		}
		q0 = up(body, tw, start, q0)
		updatesel(w, *sel, q0, q1)
	case "down":
		body, start, q0, q1, err := downtext(w, *sel)
		if err != nil {
			panic(err)
		}
		if *sel {
			q1 = down(body, tw, start, q1)
		} else {
			q0 = down(body, tw, start, q0)
		}
		updatesel(w, *sel, q0, q1)
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
