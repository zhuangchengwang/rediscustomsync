package library

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func Readlines(filename string) (lines []string, err error) {
	_, err = os.Stat(filename)
	if err != nil {
		return
	}

	file, err := os.Open(filename)
	if err != nil {
		return
	}

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		line = strings.Trim(line, "\n \t")
		if err != nil || err == io.EOF {
			break
		}
		if line != "" {
			lines = append(lines, line)
		}
	}
	return
}
