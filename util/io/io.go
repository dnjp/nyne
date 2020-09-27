package io

import (
	"fmt"
	"os"
)

func PrintErr(err error) {
	fmt.Fprintf(os.Stderr, "%v", err)
}