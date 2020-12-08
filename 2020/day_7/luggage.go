package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

const inputFile string = "luggage.input"

func main() {
	taskOne()
	taskTwo()
}

func taskOne() {
	fmt.Println("== Task one ==")

	// Rules specify, for each colour, which colours *it must contain*,
	// whereas we care about *which colours gold can be contained by*.
	mustContain := loadRules()

	// Which we invert to get a map of what each colour may be contained
	// in.
	canBeContainedIn := invertRules(mustContain)

	validColours := getOuterColoursFor("shiny gold", canBeContainedIn)
	fmt.Printf("Valid colours to contain shiny gold: %d\n", len(validColours))
}

func taskTwo() {
	fmt.Println("== Task two ==")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Return valid bag colours which the given bag colour may be contained in.
func getOuterColoursFor(colour string, canBeContainedIn map[string]map[string]int) map[string]bool {
	return getOuterColoursForRecurse(colour, canBeContainedIn, nil, 0)
}

// This is the recursion body of `getOuterColoursFor`, not to be used directly.
func getOuterColoursForRecurse(colour string, canBeContainedIn map[string]map[string]int, validOuterColours map[string]bool, depth int) map[string]bool {
	if validOuterColours == nil {
		validOuterColours = make(map[string]bool)
	}

	outerColours, ok := canBeContainedIn[colour]
	if !ok {
		// Colour can seemingly not be contained in anything
		return validOuterColours
	}

	fmt.Printf("%sColour %s may be contained in:\n", logPad(depth), colour)
	for outerColour, count := range outerColours {
		fmt.Printf("%s %s (x %d)\n", logPad(depth), outerColour, count)

		if _, ok := validOuterColours[outerColour]; !ok {
			fmt.Printf("%s This is new information, adding to list\n", logPad(depth))
			validOuterColours[outerColour] = true

			// outerColour was not yet known to be able to contain this, recurse to find colours it can be contained by
			newValidColours := getOuterColoursForRecurse(outerColour, canBeContainedIn, validOuterColours, depth+1)
			for colour, _ := range newValidColours {
				validOuterColours[colour] = true
			}
		} else {
			// outerColour already known to be able to contain this, skip it
			fmt.Printf("%s Already known, skipping\n", logPad(depth))
		}
	}

	return validOuterColours
}

func invertRules(canContain map[string]map[string]int) map[string]map[string]int {
	canBeContainedIn := make(map[string]map[string]int)

	for outerColour, outerRules := range canContain {
		// outerColour can contain any of colours listed in outerRules
		// => Any of colours in outerRules can be contained by outerColour

		for innerColour, count := range outerRules {
			if _, ok := canBeContainedIn[innerColour]; !ok {
				canBeContainedIn[innerColour] = make(map[string]int)
			}

			canBeContainedIn[innerColour][outerColour] = count
		}
	}

	return canBeContainedIn
}

func logPad(depth int) string {
	formatString := "%" + strconv.Itoa(depth*2) + "s[%d] "
	return fmt.Sprintf(formatString, "", depth)
}

// A map containing, for each colour, a map of colours it must contain.
//
// As an example, consider the following:
// {
//   "blue": { "red": 2, "yellow": 7 },
//   "gray": { },
//   "yellow": { "black": 1, "gray": 3 },
// }
// This means that blue bags must contain 2 red and 7 yellow bags. Gray bags must
// not contain anything, and yellow bags must contain 1 black and 3 gray bags.
func loadRules() map[string]map[string]int {
	rules := make(map[string]map[string]int)

	data, err := ioutil.ReadFile(inputFile)
	// Remove trailing newline
	data = data[:len(data)-1]
	check(err)

	// Matches the full line, capturing all 'can be contained in y' rules in one group
	ruleMatcher := regexp.MustCompile("^([a-z ]+) bags contain (.*)\\.$")

	for _, line := range strings.Split(string(data), "\n") {
		match := ruleMatcher.FindStringSubmatch(line)
		if match == nil {
			panic(fmt.Sprintf("Line did not match pattern: %s\n", line))
		}

		parentColour := match[1]
		childRules := match[2]

		if match[2] == "no other bags" {
			rules[parentColour] = make(map[string]int)
		} else {
			rules[parentColour] = extractChildRules(childRules)
		}
	}

	return rules
}

func extractChildRules(line string) map[string]int {
	rules := make(map[string]int)

	// Individual rules are split by ", "
	for _, childRule := range strings.Split(line, ", ") {
		// Matches individual 'can be contained in y' rule
		containedInMatcher := regexp.MustCompile("^(\\d+) ([a-z ]+) bags?$")

		innerMatch := containedInMatcher.FindStringSubmatch(childRule)
		if innerMatch == nil {
			panic(fmt.Sprintf("Rule did not match pattern: '%s'\n", childRule))
		}

		childColour := innerMatch[2]
		childCount, err := strconv.Atoi(innerMatch[1])
		check(err)

		rules[childColour] = childCount
	}

	return rules
}
