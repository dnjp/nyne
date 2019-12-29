package nyne

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"unicode/utf8"

	"9fans.net/go/acme"
	"git.sr.ht/~danieljamespost/nyne/pkg/nyne/golang"
	"git.sr.ht/~danieljamespost/nyne/pkg/util/config"
)

func New(conf *config.Config) {
	l, err := acme.Log()
	if err != nil {
		log.Fatal(err)
	}
	for {
		event, err := l.Read()
		if err != nil {
			log.Fatal(err)
		}
		switch event.Op {
		case "put":
			runCmd(conf, event)
		case "new":
			runFmt(conf, event)
			printMenu(conf, event)
		}
	}
}

func runFmt(conf *config.Config, event acme.LogEvent) {
	for _, spec := range conf.Spec {
		for _, ext := range spec.Ext {
			if strings.HasSuffix(event.Name, ext) {
				format(event, spec.Fmt)
			}
	    	}
    	}
}

func runCmd(conf *config.Config, event acme.LogEvent) {
	for _, spec := range conf.Spec {
		for _, ext := range spec.Ext {
			if strings.HasSuffix(event.Name, ext) {
				for _, cmd := range spec.Cmd {
					args := replaceName(cmd.Args, event.Name)
					reformat(event.ID, event.Name, cmd.Exec, args, ext)
		    		}
		    	}
	    	}
    	}
}

func printMenu(conf *config.Config, event acme.LogEvent) {
	w, err := acme.Open(event.ID, nil)
	if err != nil {
		log.Print(err)
		return
	}
	defer w.CloseFiles()
	if err := w.Fprintf("tag", "%s", "\n"); err != nil {
		log.Print(err)
	}

	for _, opt := range conf.Menu {
		cmd := fmt.Sprintf(" (%s)", opt)
		if err := w.Fprintf("tag", "%s", cmd); err != nil {
			log.Print(err)
		}
	}

}

func format(event acme.LogEvent, f config.Format) {
	w, err := acme.Open(event.ID, nil)
	if err != nil {
		log.Print(err)
		return
	}
	defer w.CloseFiles()

	if f.Indent != 0 {
		tabCmd := fmt.Sprintf("Tab %d", f.Indent)
		if err := w.Ctl("cleartag"); err != nil {
			log.Print(err)
		}
		tag, err := w.ReadAll("tag")
		if err != nil {
			log.Print(err)
		}
		offset := utf8.RuneCount(tag)
		cmdlen := utf8.RuneCountInString(tabCmd)
		if err := w.Fprintf("tag", "%s", tabCmd); err != nil {
			log.Print(err)
		}
		evt := new(acme.Event)
		evt.C1 = 'M'
		evt.C2 = 'x'
		evt.Q0 = offset
		evt.Q1 = offset + cmdlen
		w.WriteEvent(evt)
	}

	if f.Expand == true {
		expCmd := fmt.Sprintf("nynetab %d", f.Indent)
		if err := w.Ctl("cleartag"); err != nil {
			log.Print(err)
		}
		tag, err := w.ReadAll("tag")
		if err != nil {
			log.Print(err)
		}
		offset := utf8.RuneCount(tag)
		cmdlen := utf8.RuneCountInString(expCmd)
		if err := w.Fprintf("tag", "%s", expCmd); err != nil {
			log.Print(err)
		}
		evt := new(acme.Event)
		evt.C1 = 'M'
		evt.C2 = 'x'
		evt.Q0 = offset
		evt.Q1 = offset + cmdlen
		w.WriteEvent(evt)
	}
}

func replaceName(arr []string, name string) []string {
	newArr := make([]string, len(arr))
	for idx, str := range arr {
		if str == "$NAME" {
			newArr[idx] = name
		} else {
			newArr[idx] = arr[idx]
		}
	}
	return newArr
}

func reformat(id int, name string, x string, args []string, ext string) {
	w, err := acme.Open(id, nil)
	if err != nil {
		log.Print(err)
		return
	}
	defer w.CloseFiles()

	old, err := ioutil.ReadFile(name)
	if err != nil {
		return
	}
	new, err := exec.Command(x, args...).CombinedOutput()
	if err != nil {
		if strings.Contains(string(new), "fatal error") {
			fmt.Fprintf(os.Stderr, "%s %s: %v\n%s", x, name, err, new)
		} else {
			fmt.Fprintf(os.Stderr, "%s", new)
		}
		return
	}

	if bytes.Equal(old, new) {
		return
	}

	if ext != ".go" {
		w.Write("ctl", []byte("clean"))
		w.Write("ctl", []byte("get"))
		return
	} else {
		golang.Reformat(name, ext, w, old, new)
	}
}

