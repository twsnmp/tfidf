package util

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"
)

// ReadLines : read lines from file with filter
func ReadLines(file, filter string) ([]string, error) {
	return ReadSplitter(file, filter, '\n')
}

// ReadSplitter : read lines from file with separator and filter
func ReadSplitter(file, filter string, splitter byte) (lines []string, err error) {
	var f *regexp.Regexp
	if filter != "" {
		f = regexp.MustCompile(filter)
	}
	fin, err := os.Open(file)
	if err != nil {
		return
	}
	r := bufio.NewReader(fin)
	for {
		line, err := r.ReadString(splitter)
		if err == io.EOF {
			break
		}
		if f != nil && !f.MatchString(line) {
			continue
		}
		line = strings.Replace(line, string(splitter), "", -1)
		lines = append(lines, line)
	}
	return
}
