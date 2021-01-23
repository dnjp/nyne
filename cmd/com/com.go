package main

import (
	"flag"
	"github.com/dnjp/nyne/gen"
	"github.com/dnjp/nyne/util/io"
	"os"
	"strings"
)

func main() {
	var fflag = flag.String("f", "", "file name to operate on")
	flag.Parse()
	samfile := os.Getenv("samfile")
	if samfile == "" && fflag != nil {
		samfile = *fflag
	}
	filename := gen.GetFileName(samfile)
	ext := gen.GetExt(filename, ".txt")
	comment := gen.Conf[ext].CmtStyle
	if len(comment) == 0 {
		comment = "# "
	}
	in, err := io.PipeIn()
	if err != nil {
		panic(err)
	}

	// parse starting/ending comment parts if present
	parts := strings.Split(strings.TrimSuffix(comment, " "), " ")
	multipart := len(parts) > 1
	var startcom string
	var endcom string
	if multipart {
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

		if multipart {
			// uncomment multipart commented line
			hasbegin := strings.Contains(line, startcom)
			hasend := strings.Contains(line, endcom)
			if hasbegin && hasend {
				nline := strings.Replace(line, startcom, "", 1)
				nline = strings.Replace(nline, endcom, "", 1)
				return nline
			}
		}

		// find first non-indentation character
		first := 0
		for _, ch := range line {
			if ch == ' ' || ch == '\t' {
				first++
				continue
			}
			break
		}

		// uncomment line if beginning charcters are the comment
		comstart := first + len(comment)
		if len(line) > comstart && line[first:comstart] == comment {
			nline := strings.Replace(line, comment, "", 1)
			return nline
		}

		// comment line using appropriate comment structure
		if multipart {
			return line[:first] + startcom + line[first:] + endcom
		}
		return line[:first] + comment + line[first:]
	})
}
