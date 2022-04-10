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

func isword(c byte) bool {
	r := rune(c)
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

func prevword(body []byte, tw, start, q int) (nq int) {
	leftv := q - start
	for i := leftv; i >= 0; i-- {
		pc := body[i]
		if i-1 <= 0 {
			return start - 1
		}
		c := body[i-1]
		if !isword(pc) && isword(c) {
			if start+i != q {
				return start + i
			}
		}
	}
	return q
}

func nextword(body []byte, tw, start, q int) (nq int) {
	leftv := q - start
	length := len(body)
	for i := leftv; i < length; i++ {
		off := -1
		if i+off < 0 {
			off = 0
		}
		pc := body[i+off]
		if i+1 >= length {
			return start + length
		}
		c := body[i]
		if !isword(pc) && isword(c) {
			if start+i != q {
				return start + i
			}
		}
	}
	return q
}

func left(body []byte, tw, start, q int) (nq int) {
	nq = q - 1
	if nq <= 0 {
		return 0
	}
	return nq
}

func right(body []byte, tw, start, q int) (nq int) {
	return q + 1
}

func up(body []byte, tw, start, q int) (nq int) {
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

func down(body []byte, tw, start, q int) (nq int) {
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

func blankline(w *nyne.Win, q int, up bool) (nq int) {
	off := 1
	regex := "+/^$/"
	if up {
		off = -1
		regex = "-/^$/"
	}
	err := w.SetAddr("#%d", q+off)
	if err != nil {
		panic(err)
	}
	err = w.SetAddr(regex)
	if err != nil {
		panic(err)
	}
	nq, _, err = w.Addr()
	if err != nil {
		panic(err)
	}
	return nq
}

func curline(w *nyne.Win, sel, incQ1 bool) (body []byte, start, q0, q1 int, err error) {
	q0, q1, err = w.CurrentAddr()
	if err != nil {
		return
	}

	if incQ1 {
		err = w.SetAddr("#%d", q1)
		if err != nil {
			return
		}
	}

	err = w.SetAddr("-+")
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

func prevline(w *nyne.Win, sel bool) (body []byte, start, q0, q1 int, err error) {
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

func nextline(w *nyne.Win, sel bool) (body []byte, start, q0, q1 int, err error) {
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

func update(w *nyne.Win, sel bool, q0, q1 int) {
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
		body, start, q0, q1, err := prevline(w, *sel)
		if err != nil {
			panic(err)
		}
		q0 = up(body, tw, start, q0)
		update(w, *sel, q0, q1)
	case "down":
		body, start, q0, q1, err := nextline(w, *sel)
		if err != nil {
			panic(err)
		}
		if *sel {
			q1 = down(body, tw, start, q1)
		} else {
			q0 = down(body, tw, start, q0)
		}
		update(w, *sel, q0, q1)
	case "left":
		body, start, q0, q1, err := curline(w, *sel, false)
		if err != nil {
			panic(err)
		}
		if *word {
			q0 = prevword(body, tw, start, q0)
		} else if *paragraph {
			q0, _, err = w.CurrentAddr()
			if err != nil {
				return
			}
			q0 = blankline(w, q0, true)
		} else {
			q0 = left(body, tw, start, q0)
		}
		update(w, *sel, q0, q1)
	case "right":
		incQ1 := false
		if *sel && *word {
			incQ1 = true
		}
		body, start, q0, q1, err := curline(w, *sel, incQ1)
		if err != nil {
			panic(err)
		}
		if *sel {
			if *word {
				q1 = nextword(body, tw, start, q1)
			} else if *paragraph {
				_, q1, err = w.CurrentAddr()
				if err != nil {
					return
				}
				q1t := q1
				q1 = blankline(w, q1, false)
				if q1 < q1t {
					// wraparound
					return
				}
			} else {
				q1 = right(body, tw, start, q1)
			}
		} else if *word {
			q0 = nextword(body, tw, start, q0)
		} else if *paragraph {
			q0, _, err = w.CurrentAddr()
			if err != nil {
				return
			}
			q0t := q0
			q0 = blankline(w, q0, false)
			if q0 < q0t {
				// wraparound
				return
			}
		} else {
			q0 = right(body, tw, start, q0)
		}
		update(w, *sel, q0, q1)
	case "start":
		_, start, _, q1, err := curline(w, *sel, false)
		if err != nil {
			panic(err)
		}
		update(w, *sel, start, q1)
	case "end":
		body, start, q0, q1, err := curline(w, *sel, false)
		if err != nil {
			panic(err)
		}
		if *sel {
			q1 = (start + len(body)) - 1
		} else {
			q0 = (start + len(body)) - 1
		}
		update(w, *sel, q0, q1)
	}
}
