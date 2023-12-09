// Package data provides various ways to load and interpret input data.
package data

import (
	"bufio"
	"os"
)

// LinesFromFile reads newline-separated lines from the input file.
func LinesFromFile(filePath string) ([]string, error) {
	out := make([]string, 0)

	file, err := os.Open(filePath)
	if err != nil {
		return out, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		out = append(out, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return out, err
	}

	return out, nil
}
