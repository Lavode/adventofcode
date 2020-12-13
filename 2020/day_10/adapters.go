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

	// Generate a list of adapters from which a given adapter can be
	// reached.
	// 5 => [2, 4] would eg mean that adapter 5 can be reached from
	// adapters 2 and 4.
	reachableFrom := make(map[int][]int)
	for idx, destination := range adapters {
		for i := idx - 1; i >= 0; i -= 1 {
			source := adapters[i]
			if source >= destination-3 {
				// Destination can be reached from source
				reachableFrom[destination] = append(reachableFrom[destination], source)
			} else {
				// As the list of adapters is sorted, no
				// earlier adapter is going to be a valid
				// source.
				break
			}
		}
	}

	// Dynamic-programming approach to calculate number of ways in which
	// the final 'adapter' (ie the device) can be reached.
	// Given a node `x` can be reached from nodes `a` and `b`, which can be
	// reached in 3 and 5 ways respectively, then there node `x` can be
	// reached in 3 + 5 = 8 ways.
	//
	// We pre-seed the source node, which can be reached in exactly 1 way
	// (and has to be part of the path).
	pathCount := make(map[int]int)
	pathCount[0] = 1
	for _, destination := range adapters {
		for _, source := range reachableFrom[destination] {
			pathCount[destination] += pathCount[source]
		}

		fmt.Printf("Node %d can be reached in %d ways\n", destination, pathCount[destination])
	}
}

// Find valid arrangments of adapters which connect outlet to personal device.
//
// This is a brute-force - albeit smart - solution which traverses all
// *possible* solutions. It will not scale well with input size.
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
