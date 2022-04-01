package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/dnjp/nyne"
)

func main() {
	os.Unsetenv("winid") // do not trust the execution environment

	winid, err := nyne.FindFocusedWinID()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not find focused window: %+v", err)
		os.Exit(1)
	}

	w, err := nyne.OpenWin(winid, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open window: %+v", err)
		os.Exit(1)
	}

	_, font, err := w.Font()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get window font: %+v", err)
		os.Exit(1)
	}

	fs, err := nyne.FontSize(font)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get font size: %+v", err)
		os.Exit(1)
	}

	nsize := fs + 1
	if nyne.IsHiDPI(font) {
		nsize = (fs / 2) + 1
	}

	fontName := strings.ReplaceAll(
		font.Name,
		strconv.Itoa(fs),
		strconv.Itoa(nsize))

	err = w.SetFont(fontName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not set font: %+v", err)
		os.Exit(1)
	}
}
