package main

import (
	"git.sr.ht/~danieljamespost/nyne/gen"
	"git.sr.ht/~danieljamespost/nyne/util/io"
	"os"
	"strings"
)

func main() {
	filename := gen.GetFileName(os.Getenv("samfile"))
	ext := gen.GetExt(filename, ".txt")
	comment := gen.Conf[ext].CommentStyle
	if len(comment) == 0 {
		comment = "# "
	}
	in, err := io.PipeIn()
	if err != nil {
		panic(err)
	}

	var startcom string
	var endcom string
	parts := strings.Split(comment, " ")
	if len(parts) > 1 {
		if len(parts[0]) > 0 {
			startcom = parts[0] + " "
		}
		if len(parts[1]) > 0 {
			endcom = " " + parts[1]
		}
	}
	io.PipeOut(in, func(line string) string {
		if len(line) == 0 {
			return line
		}

		if len(parts) > 0 {
			hasbegin := strings.Contains(line, startcom)
			hasend := strings.Contains(line, endcom)
			if hasbegin && hasend {
				nline := strings.Replace(line, startcom, "", 1)
				nline = strings.Replace(nline, endcom, "", 1)
				return nline
			}
		} else {
			if strings.Contains(line, comment) {
				nline := strings.Replace(line, comment, "", 1)
				return nline
			}
		}

		first := 0
		for _, ch := range line {
			if ch == ' ' || ch == '\t' {
				first++
				continue
			}
			break
		}
		var nline string
		if len(parts) > 1 {
			nline = line[:first] + startcom + line[first:] + endcom
		} else {
			nline = line[:first] + comment + line[first:]
		}
		return nline
	})
}
