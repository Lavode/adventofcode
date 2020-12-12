package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

const inputFile string = "adapters.input"

func main() {
	taskOne()
	taskTwo()
}

func taskOne() {
	fmt.Println("== Task one ==")

	// The charging outlet has an implicit Joltate rating of 0. As we don't
	// known the rating of the first adapter we cannot compensate for it by
	// eg starting `oneJoltDifferences` at 1, rather we simply add it as a
	// 'fake' adapter.
	adapters := loadAdapters()
	adapters = prependInt(adapters, 0)

	oneJoltDifferences := 0
	// Our device is implicitly three Jolts higher than the highest
	// adapter, starting the 3-jolt-difference at 1.
	threeJoltDifferences := 1

	for i := 0; i < len(adapters)-1; i += 1 {
		difference := adapters[i+1] - adapters[i]

		switch difference {
		case 1:
			oneJoltDifferences += 1
		case 3:
			threeJoltDifferences += 1
		default:
		}
	}

	fmt.Printf(
		"1-Jolt differences: %d, 3-Jolt differences: %d, Product: %d\n",
		oneJoltDifferences,
		threeJoltDifferences,
		oneJoltDifferences*threeJoltDifferences,
	)
}

func taskTwo() {
	fmt.Println("== Task two ==")

	// Prepend 'fake' 0-rated adapter of outlet
	adapters := loadAdapters()
	adapters = prependInt(adapters, 0)

	// We require the outlet to be part of the sequence, so can simply
	// start iteration at it.
	count := findSolutions(adapters, 0)
	fmt.Printf("Found valid solutions: %d\n", count)
}

// Find valid arrangments of adapters which connect outlet to personal device.
func findSolutions(adapters []int, currentIdx int) int {
	if currentIdx == len(adapters)-1 {
		// Reached end of list => Found valid solution
		return 1
	}

	currentCount := 0

	currentRating := adapters[currentIdx]
	// There are no duplicates, so our next valid adapter will be 1 to 3
	// higher than the current one. As entries are further sorted, we only
	// care about the maximum rating.
	maxRating := currentRating + 3

	// Start with the next possible candidate
	for i := currentIdx + 1; i < len(adapters) && adapters[i] <= maxRating; i += 1 {
		currentCount += findSolutions(adapters, i)
	}

	return currentCount
}

func prependInt(slice []int, x int) []int {
	// Increase capacity of slice by 1
	slice = append(slice, 0)
	// Move all items one over
	copy(slice[1:], slice)
	// And prepend item
	slice[0] = x

	return slice
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func loadAdapters() []int {
	data, err := ioutil.ReadFile(inputFile)
	// Remove trailing newline
	data = data[:len(data)-1]
	check(err)

	var adapters []int

	// ioutil.ReadFile returns a byte slice, strings.Split expects a string
	for _, s := range strings.Split(string(data), "\n") {
		i, err := strconv.Atoi(s)
		check(err)

		adapters = append(adapters, i)
	}

	sort.Ints(adapters)

	return adapters
}
