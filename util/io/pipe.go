package io

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// PipeIn reads from piped stdin and returns the result
func PipeIn() ([]rune, error) {
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
		return in, fmt.Errorf("must be used with pipe")
	}
	return in, nil
}

// PipeOut applies the given line transformation for each line before
// printing to stdout
func PipeOut(in []rune, fn func(string) string) {
	out := []string{}
	lines := strings.Split(string(in), "\n")
	for _, line := range lines {
		nline := fn(line)
		out = append(out, nline)
	}
	fmt.Print(strings.Join(out, "\n"))
}
