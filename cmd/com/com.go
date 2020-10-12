package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

func main() {
	var comment string
	if len(os.Args) > 1 {
		comment = os.Args[1]
	} else {
		comment = os.Getenv("COMMENTC")
	}
	if len(comment) == 0 {
		panic(fmt.Errorf("no comment type supplied, " +
			"expected arg or $COMMENTC"))
	}
	reader := bufio.NewReader(os.Stdin)
	var in []rune
	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		in = append(in, input)
	}
	if len(in) == 0 {
		panic("must be used with pipe")
	}
	out := []string{}
	for _, line := range strings.Split(string(in), "\n") {
		if len(line) == 0 {
			out = append(out, line)
			continue
		}
		if strings.Contains(line, comment) {
			nline := strings.Replace(line, comment+" ", "", 1)
			out = append(out, nline)
		} else {
			first := 0
			for _, ch := range line {
				if unicode.IsLetter(ch) {
					break
				}
				first += 1	
			}
			nline := line[:first] + comment + " " + line[first:]
			out = append(out, nline)		
		}

	}
	fmt.Printf(strings.Join(out, "\n"))
}