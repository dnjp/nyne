package nyne

import (
	"fmt"
	"os"
)

// Error prints to stderr
func Error(err error) {
	fmt.Fprintf(os.Stderr, "%+v", err)
}
