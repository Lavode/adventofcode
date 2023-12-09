// Package main solves the Advent of Code challenges. Or at least some of them. :)
package main

import (
	"log"
	"regexp"

	"github.com/lavode/adventofcode/2023/pkg/data"
)

const dataRoot = "data/"

func main() {
	one()
}

func one() {
	knownDigitsPattern := regexp.MustCompile("(1|2|3|4|5|6|7|8|9|one|two|three|four|five|six|seven|eight|nine)")
	knownDigits := map[string]int{
		"1":     1,
		"one":   1,
		"2":     2,
		"two":   2,
		"3":     3,
		"three": 3,
		"4":     4,
		"four":  4,
		"5":     5,
		"five":  5,
		"6":     6,
		"six":   6,
		"7":     7,
		"seven": 7,
		"8":     8,
		"eight": 8,
		"9":     9,
		"nine":  9,
	}
	sum := 0

	input, err := data.LinesFromFile(dataRoot + "1.txt")
	if err != nil {
		log.Fatalf("Error reading input file: %v", err)
	}

	for _, line := range input {
		// Find first digit
		for i := 0; i < len(line); i++ {
			uptoHere := string(line[0 : i+1])
			if match := knownDigitsPattern.FindString(uptoHere); match != "" {
				sum += knownDigits[match] * 10
				break
			}

		}

		// And the last digit
		for i := len(line) - 1; i >= 0; i-- {
			uptoHere := string(line[i:len(line)])
			if match := knownDigitsPattern.FindString(uptoHere); match != "" {
				sum += knownDigits[match]
				break
			}

		}
	}

	log.Printf("Sum of first+last digit of all lines: %v", sum)
}
