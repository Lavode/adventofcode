package main

import "fmt"
import "io/ioutil"
import "strconv"
import "strings"

const inputFile string = "expense_report.input"

func main() {
	taskOne()
	taskTwo()
}

func taskOne() {
	fmt.Println("== Task one ==")
	expenses := getExpenses()

	x, y, err := findPairWithSum(expenses, 2020)
	if err == nil {
		fmt.Printf(
			"Match found: %d + %d = 2020, %d * %d = %d\n",
			x, y,
			x, y,
			x*y,
		)
	} else {
		fmt.Println("No match found, something's amiss.")
	}
}

func taskTwo() {
	fmt.Println("== Task two ==")
	expenses := getExpenses()

	x, y, z, err := findTripletWithSum(expenses, 2020)
	if err == nil {
		fmt.Printf(
			"Match found: %d + %d + %d = 2020, %d * %d * %d = %d\n",
			x, y, z,
			x, y, z,
			x*y*z,
		)
	} else {
		fmt.Println("No match found, something's amiss.")
	}
}

// Find a pair of numbers x and y in a slice of numbers with x + y = sum.
//
// Returns an error if no pair found.
// O(n) implementation, utilizing a hashmap as auxiliary data structure.
func findPairWithSum(numbers []int, sum int) (int, int, error) {
	// We don't really care about the value type of the map, but GoLang has
	// no native support for sets
	cache := make(map[int]bool)

	for _, a := range numbers {
		b := sum - a

		if _, ok := cache[b]; ok {
			// A value `b` has been seen in `numbers` earlier.
			return a, b, nil
		} else {
			// No value `b` was seen in `numbers` so far, add
			// current value `a` to cache
			cache[a] = true
		}
	}

	// No matching pair found
	return 0, 0, fmt.Errorf("No pair (x, y) with x + y = %d found", sum)
}

// Find three numbers `x`, `y`, `z` in a slice of numbers with `x + y + z = sum`.
//
// Returns an error if no triplet found.
// O(n^2) implementation, utilizing a hashmap as auxiliary data structure.
func findTripletWithSum(numbers []int, sum int) (int, int, int, error) {
	for _, x := range numbers {
		y, z, err := findPairWithSum(numbers, sum-x)
		if err == nil {
			// x + y + z = sum
			return x, y, z, nil
		}
	}

	// No `y` and `z` found for any of the possible `x`
	return 0, 0, 0, fmt.Errorf("No triplet (x, y, z) with x + y + z = %d found", sum)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getExpenses() []int {
	data, err := ioutil.ReadFile(inputFile)
	// Remove trailing newline
	data = data[:len(data)-1]
	check(err)

	var expenses []int

	// ioutil.ReadFile returns a byte slice, strings.Split expects a string
	for _, s := range strings.Split(string(data), "\n") {
		i, err := strconv.Atoi(s)
		check(err)

		expenses = append(expenses, i)
	}

	return expenses
}
