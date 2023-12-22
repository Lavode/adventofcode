// Package main solves the Advent of Code challenges. Or at least some of them. :)
package main

import (
	"log"

	"github.com/lavode/adventofcode/2023/pkg/data"
	"github.com/lavode/adventofcode/2023/pkg/nlp"
)

const dataRoot = "data/"

func main() {
	one()
}

func one() {
	sum := 0

	input, err := data.LinesFromFile(dataRoot + "1.txt")
	if err != nil {
		log.Fatalf("Error reading input file: %v", err)
	}

	for _, line := range input {
		if digit, ok := nlp.FindDigitFromFront(line, true); ok {
			sum += digit * 10
		} else {
			log.Fatalf("Did not find any digit in line: %s", line)
		}

		if digit, ok := nlp.FindDigitFromBack(line, true); ok {
			sum += digit
		} else {
			log.Fatalf("Did not find any digit in line: %s", line)
		}
	}

	log.Printf("Sum of first+last digit of all lines: %v", sum)
}
