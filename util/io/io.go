package io

import (
	"fmt"
	"os"
)

// PrintErr prints to OS stderr
func PrintErr(err error) {
	fmt.Fprintf(os.Stderr, "%v", err)
}
