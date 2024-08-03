// Package main solves the Advent of Code challenges. Or at least some of them. :)
package main

import (
	"log"

	"github.com/lavode/adventofcode/2023/pkg/balls"
	"github.com/lavode/adventofcode/2023/pkg/data"
	"github.com/lavode/adventofcode/2023/pkg/nlp"
)

const dataRoot = "data/"

func main() {
	// one()
	two()
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

func two() {
	input, err := data.LinesFromFile(dataRoot + "2.txt")
	if err != nil {
		log.Fatalf("Error reading input file: %v", err)
	}

	games := make([]balls.Game, 0)
	for _, line := range input {
		game, err := balls.ParseGame(line)
		if err != nil {
			log.Fatalf("Error parsing game: %v", err)
		}
		games = append(games, game)
	}

	sumOfMatching := 0
	target := map[string]int{"red": 12, "green": 13, "blue": 14}
	for _, game := range games {
		bound := game.BoundOnBalls()

		if target["red"] >= bound["red"] && target["green"] >= bound["green"] && target["blue"] >= bound["blue"] {
			// Game might have happened, with target bag
			sumOfMatching += int(game.Id)
		}
	}

	log.Printf("Sum of plausible IDs: %d", sumOfMatching)
}
