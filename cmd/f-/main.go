package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"9fans.net/go/acme"
)

func main() {
	numre, err := regexp.Compile("[0-9]+")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create regex: %+v", err)
		os.Exit(1)
	}
	winid := os.Getenv("winid")
	if winid == "" {
		fmt.Fprintf(os.Stderr, "cannot use f+ outside of acme")
		os.Exit(1)
	}
	id, err := strconv.Atoi(winid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse winid: %+v", err)
		os.Exit(1)
	}
	win, err := acme.Open(id, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open window: %+v", err)
		os.Exit(1)
	}
	_, font, err := win.Font()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get window font: %+v", err)
		os.Exit(1)
	}
	sizes := numre.FindAllString(font.Name, -1)
	if l := len(sizes); l > 1 || l == 0 {
		fmt.Fprintf(os.Stderr, "could not parse font size")
		os.Exit(1)
	}
	fs, err := strconv.Atoi(sizes[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse font size: %+v", err)
		os.Exit(1)
	}
	ns := (fs / 2) - 1
	fontName := strings.ReplaceAll(
		font.Name,
		strconv.Itoa(fs),
		strconv.Itoa(ns),
	)
	err = win.Ctl("font %s", fontName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not send ctl: %+v", err)
		os.Exit(1)
	}
}
