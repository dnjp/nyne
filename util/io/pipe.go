package io

import (
	"os"
	"fmt"
	"bufio"
	"io"
)

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