package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

const inputFile string = "customs.input"

func main() {
	taskOne()
	taskTwo()
}

func taskOne() {
	fmt.Println("== Task one ==")

	answers := loadAnswers()
	totalYesCounts := 0
	for _, group := range answers {
		// We care about the total yes counts *per group*, ie two 'yes'
		// counts for answer 'x' within the same group only count once.
		// This is equivalent to the number of keys within the map
		// minus the key which is used to track the number of people in
		// the group.
		totalYesCounts += len(group) - 1
	}

	fmt.Printf("Total yes counts (grouped by group): %d\n", totalYesCounts)
}

func taskTwo() {
	fmt.Println("== Task two ==")

	answers := loadAnswers()
	totalYesCounts := 0
	for _, group := range answers {
		memberCount, ok := group['_']
		// This count must be present
		if !ok {
			panic(fmt.Sprintf("Group member count not present for group: %+v\n", group))
		}

		for answer, yesCount := range group {
			// Ignore the meta group-size key
			if answer != '_' {
				if memberCount == yesCount {
					totalYesCounts += 1
				}
			}
		}
	}

	fmt.Printf("Total yes counts (grouped by group, all members must have said yes): %d\n", totalYesCounts)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Prepare per-group count of how many people answered 'yes' to each question.
//
// Questions to which none answered are not contained in the inner map. Further
// each map contains a '_' key which stores the number of members the group
// has.
func loadAnswers() []map[rune]int {
	var answers []map[rune]int

	data, err := ioutil.ReadFile(inputFile)
	check(err)
	// Keep trailingn newline, as those separate groups of answers

	group := make(map[rune]int)
	for _, line := range strings.Split(string(data), "\n") {
		// Empty lines separate individual groups
		if line == "" {
			answers = append(answers, group)
			group = make(map[rune]int)
		} else {
			// Each line is the list of one-char answers the person
			// answered 'yes' to, eg 'abjz'.
			for _, answer := range line {
				group[answer] += 1
			}
			// And each line itself is the answers of one person
			group['_'] += 1
		}
	}

	return answers
}
