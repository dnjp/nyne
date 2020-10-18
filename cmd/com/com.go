package main

import (
	"fmt"
	"git.sr.ht/~danieljamespost/nyne/gen"
	"git.sr.ht/~danieljamespost/nyne/util/io"
	"os"
	"strings"
	"unicode"
)

func main() {
	var comment string
	if len(os.Args) > 1 {
		comment = os.Args[1]
	} else {
		filename := gen.GetFileName(os.Getenv("samfile"))
		ext := gen.GetExt(filename, ".txt")
		comment = gen.Conf[ext].CommentStyle

	}
	if len(comment) == 0 {
		panic(fmt.Errorf("no comment type supplied, " +
			"expected arg or $COMMENTC"))
	}
	in, err := io.PipeIn()
	if err != nil {
		panic(err)
	}

	io.PipeOut(in, func(line string) string {
		if len(line) == 0 {
			return line
		}
		if strings.Contains(line, comment) {
			nline := strings.Replace(line, comment+" ", "", 1)
			return nline
		}
		first := 0
		for _, ch := range line {
			if unicode.IsLetter(ch) {
				break
			}
			first += 1
		}
		nline := line[:first] + comment + " " + line[first:]
		return nline
	})
}
