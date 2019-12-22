package golang

import (
	"fmt"
	"io/ioutil"
	"log"
	"bytes"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"unicode/utf8"

	"9fans.net/go/acme"
)

func Reformat(name string, ext string, w *acme.Win, old []byte, new []byte) {
	if ext == ".go" {
		oldTop, err := readImports(bytes.NewReader(old), true)
		if err != nil {
			//log.Print(err)
			return
		}
		newTop, err := readImports(bytes.NewReader(new), true)
		if err != nil {
			//log.Print(err)
			return
		}
		if bytes.Equal(oldTop, newTop) {
			return
		}
		w.Addr("0,#%d", utf8.RuneCount(oldTop))
		w.Write("data", newTop)
		return
	}

	f, err := ioutil.TempFile("", "nyne")
	if err != nil {
		log.Print(err)
		return
	}
	if _, err := f.Write(new); err != nil {
		log.Print(err)
		return
	}
	tmp := f.Name()
	f.Close()
	defer os.Remove(tmp)

	diff, _ := exec.Command("9", "diff", name, tmp).CombinedOutput()

	latest, err := w.ReadAll("body")
	if err != nil {
		log.Print(err)
		return
	}
	if !bytes.Equal(old, latest) {
		log.Printf("skipped update to %s: window modified since Put\n", name)
		return
	}
	w.Write("ctl", []byte("mark"))
	w.Write("ctl", []byte("nomark"))
	diffLines := strings.Split(string(diff), "\n")
	for i := len(diffLines) - 1; i >= 0; i-- {
		line := diffLines[i]
		fmt.Println(line);
		if line == "" {
			continue
		}
		if line[0] == '<' || line[0] == '-' || line[0] == '>' {
			continue
		}
		j := 0
		for j < len(line) && line[j] != 'a' && line[j] != 'c' && line[j] != 'd' {
			j++
		}
		if j >= len(line) {
			log.Printf("cannot parse diff line: %q", line)
			break
		}
		oldStart, oldEnd := parseSpan(line[:j])
		newStart, newEnd := parseSpan(line[j+1:])
		if oldStart == 0 || newStart == 0 {
			continue
		}
		switch line[j] {
		case 'a':
			err := w.Addr("%d+#0", oldStart)
			if err != nil {
				log.Print(err)
				break
			}
			w.Write("data", findLines(new, newStart, newEnd))
		case 'c':
			err := w.Addr("%d,%d", oldStart, oldEnd)
			if err != nil {
				log.Print(err)
				break
			}
			w.Write("data", findLines(new, newStart, newEnd))
		case 'd':
			err := w.Addr("%d,%d", oldStart, oldEnd)
			if err != nil {
				log.Print(err)
				break
			}
		}
	}
	
	// update buffer
	w.Write("ctl", []byte("clean"))
	w.Write("ctl", []byte("get"))

}

func parseSpan(text string) (start, end int) {
	i := strings.Index(text, ",")
	if i < 0 {
		n, err := strconv.Atoi(text)
		if err != nil {
			log.Printf("cannot parse span %q", text)
			return 0, 0
		}
		return n, n
	}
	start, err1 := strconv.Atoi(text[:i])
	end, err2 := strconv.Atoi(text[i+1:])
	if err1 != nil || err2 != nil {
		log.Printf("cannot parse span %q", text)
		return 0, 0
	}
	return start, end
}

func findLines(text []byte, start, end int) []byte {
	i := 0

	start--
	for ; i < len(text) && start > 0; i++ {
		if text[i] == '\n' {
			start--
			end--
		}
	}
	startByte := i
	for ; i < len(text) && end > 0; i++ {
		if text[i] == '\n' {
			end--
		}
	}
	endByte := i
	return text[startByte:endByte]
}