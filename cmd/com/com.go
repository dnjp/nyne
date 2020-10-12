package main

import (

	"fmt"
	"os"
	"strings"
	"unicode"
	"git.sr.ht/~danieljamespost/nyne/util/io"
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
	in, err := io.PipeIn()
	if err != nil {
		panic(err)
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