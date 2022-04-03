/*
Comments/uncomments piped text

	Usage of com:
		|com
*/
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/dnjp/nyne"
)

func main() {
	filename := os.Getenv("samfile")
	if filename == "" {
		filename = os.Getenv("%")
	}
	if filename == "" {
		fmt.Fprintf(os.Stderr, "$samfile and $%% are empty. are you sure you're in acme?")
		os.Exit(1)
	}

	ft, _ := nyne.FindFiletype(nyne.Filename(filename))
	comment := ft.Comment
	if len(comment) == 0 {
		comment = "# "
	}

	var startcom string = comment
	var endcom string
	parts := strings.Split(strings.TrimSuffix(comment, " "), " ")
	if len(parts) > 1 {
		if len(parts[0]) > 0 {
			startcom = parts[0] + " "
		}
		if len(parts[1]) > 0 {
			endcom = " " + parts[1]
		}
	}
	startcomlen := len(startcom)
	endcomlen := len(endcom)

	var commentedLines, uncommentedLines int
	var commentIdx int
	var hasCommented bool

	in := []byte{}
	reader := bufio.NewReader(os.Stdin)
	for {
		b, err := reader.ReadByte()
		if err != nil && err == io.EOF {
			break
		}
		if commentIdx < len(startcom) && b == startcom[commentIdx] {
			commentIdx++
			if commentIdx == len(startcom) {
				hasCommented = true
				commentIdx = 0
			}
		} else {
			commentIdx = 0
		}
		if b == '\n' {
			if hasCommented {
				commentedLines++
				hasCommented = false
			} else {
				uncommentedLines++
			}
		}
		in = append(in, b)
	}

	shouldComment := uncommentedLines > commentedLines
	var i, comidx int
	var hasstart bool
	fromstart := 0
	prevstartcom := -1
	bufa := -1

	for {
		if i >= len(in) {
			break
		}

		shouldStart := startcomlen > 0 && !hasstart && shouldComment
		inlineWithPrevStart := prevstartcom >= 0 && fromstart >= prevstartcom && !hasstart
		isNotWhitespace := in[i] != '\t' && in[i] != ' '
		inComment := comidx > 0

		if !shouldComment {
			if comidx < len(endcom) && in[i] == endcom[comidx] {
				if bufa < 0 {
					bufa = i
				}
				comidx++
				if comidx == len(endcom) {
					bufa = -1
				}
				i++
				continue
			} else if comidx < len(startcom) && in[i] == startcom[comidx] {
				if bufa < 0 {
					bufa = i
				}
				comidx++
				if comidx == len(startcom) {
					for j := bufa + len(startcom); j < i; j++ {
						fmt.Fprintf(os.Stdout, "%c", in[j])
					}
					bufa = -1
				}
				i++
				continue
			} else {
				if bufa >= 0 {
					for j := bufa; j < i; j++ {
						fmt.Fprintf(os.Stdout, "%c", in[j])
					}
					comidx = 0
					bufa = -1
				} else {
					comidx = 0
					goto Out
				}
			}
		} else if in[i] == '\n' {
			if !hasstart {
				fromstart = -1
				prevstartcom = -1
				goto Out
			}
			if endcomlen > 0 {
				if comidx >= endcomlen {
					comidx = 0
				} else {
					fmt.Fprintf(os.Stdout, "%c", endcom[comidx])
					comidx++
					fromstart++
					continue
				}
			}
			hasstart = false
			fromstart = -1
		} else if inlineWithPrevStart || (shouldStart && (isNotWhitespace || inComment)) {
			if comidx >= startcomlen {
				comidx = 0
				hasstart = true
			} else {
				if prevstartcom < 0 {
					prevstartcom = fromstart
				}
				fmt.Fprintf(os.Stdout, "%c", startcom[comidx])
				comidx++
				fromstart++
				continue
			}
		}

	Out:
		fmt.Fprintf(os.Stdout, "%c", in[i])
		i++
		fromstart++
	}
}
