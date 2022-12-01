package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// const inputFile = "input.test.txt"
const inputFile = "input.txt"

func main() {
	input := getInput()


	fmt.Printf(
		"Total increasing measures: %d\n",
		countIncreasingMeasures(input),
	)

	fmt.Printf(
		"Total increasing measures (3-sliding window): %d\n",
		countIncreasingSlidingWindow(input),
	)
}

func countIncreasingMeasures(input []int) int {
	// Count the number of measures which have increased compared to the
	// previous one
	increasingMeasuresCount := 0

	previous := input[0]
	for i := 1; i < len(input); i += 1 {
		if input[i] > previous {
			increasingMeasuresCount += 1
		}

		previous = input[i]
	}

	return increasingMeasuresCount
}

func countIncreasingSlidingWindow(input []int) int {
	// Count increases across a sliding window of size three
	increasingMeasuresCount := 0

	// idx will be starting index of first window in comparison
	// That is the first window will be [idx, idx + 1, idx + 2],
	// the second window [idx + 1, idx + 2, idx + 3]
	for idx := 0; idx < len(input) - 3; idx += 1 {
		// We could optimize this by keeping track of the sum,
		// subtracting the old and adding the new values. But this
		// would only save us one operation per window, so seems not
		// worth the effort.
		leading_window := input[idx] + input[idx + 1] + input[idx + 2]
		lagging_window := input[idx + 1] + input[idx + 2] + input[idx + 3]

		if lagging_window > leading_window {
			increasingMeasuresCount += 1
		}
	}

	return increasingMeasuresCount;
}

func getInput() []int {
	content, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Error reading file: %v\n", err)
	}

	lines := strings.Split(string(content), "\n")
	// Drop trailing newline
	lines = lines[:len(lines) - 1]

	input := make([]int, len(lines))

	for idx, line := range lines {
		i, err := strconv.Atoi(line)
		if err != nil {
			log.Fatalf("Error reading file: %v\n", err)
		}

		input[idx] = i
	}

	return input
}
