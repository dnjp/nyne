// Acmego watches acme for .go files being written.
// Each time a .go file is written, acmego checks whether the
// import block needs adjustment. If so, it makes the changes
// in the window body but does not write the file.
package nyne

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

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
		if event.Name != "" && event.Op == "put" {
			for _, spec := range conf.Spec {
				for _, ext := range spec.Ext {
					if strings.HasSuffix(event.Name, ext) {
						for _, cmd := range spec.Cmd {
							args := replaceName(cmd.Args, event.Name)
							reformat(event.ID, event.Name, cmd.Exec, args, ext)
					    		event, err = l.Read()
					    		if err != nil {
					    			log.Fatal(err)
					    		}
				    		}
				    	}
			    	}
		    	}
		}
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

