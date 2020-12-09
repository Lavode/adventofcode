package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

const inputFile string = "xmas_crypto.input"

// 'Manual' queue, which has to be told when to shift its window.
type Queue struct {
	slice []int

	start int
	end   int
}

func (queue Queue) String() string {
	return fmt.Sprintf("%+v", queue.slice[queue.start:queue.end])
}

func (queue Queue) Items() []int {
	return queue.slice[queue.start:queue.end]
}

func (queue *Queue) Shift() error {
	if queue.end <= len(queue.slice)-1 {
		// Able to shift at least one more space
		queue.start += 1
		queue.end += 1
		return nil
	} else {
		return fmt.Errorf("End of queue reached")
	}
}

func main() {
	taskOne()
	taskTwo()
}

func taskOne() {
	fmt.Println("== Task one ==")

	numbers := loadNumbers()
	queue := Queue{slice: numbers, end: 25}

	// We'll start with the 26th items, with the 'queue' containing the 25
	// preceding it.
	for i := 25; i < len(numbers); i += 1 {
		sum := numbers[i]
		_, _, err := findPairWithSum(queue.Items(), sum)
		if err != nil {
			fmt.Printf("Error: %d can not be expressed as sum of %+v!\n", sum, queue)
		}

		e := queue.Shift()
		if e != nil {
			// This shouldn't happen, as each iteration will only
			// advance the queue's window by one.
			panic(fmt.Sprintf("Ran out of items while shifting queue\n"))
		}
	}
}

func taskTwo() {
	fmt.Println("== Task two ==")

	numbers := loadNumbers()
	// We now need to find a continguous sequence of numbers in the input
	// which sums to the number found above.
	// We'll do so by starting with the pair of the first two items, adding
	// the following items until we exceed the sum.  If we exceed the sum,
	// we start with the pair of items 2 and 3, and so on.
	// This is an O(n^2) solution, but for this case good enough.
	goal := 90433990
	for start := 0; start <= len(numbers)-2; start += 1 {
		for end := start + 2; end <= len(numbers); end += 1 {
			sum := 0
			for _, x := range numbers[start:end] {
				sum += x
			}

			if sum == goal {
				// We got the goal :)
				fmt.Printf("Success: Indices %d through %d sum to %d\n", start, end, goal)
				fmt.Printf("Numbers: %+v\n", numbers[start:end])
			}
		}
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

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func loadNumbers() []int {
	data, err := ioutil.ReadFile(inputFile)
	// Remove trailing newline
	data = data[:len(data)-1]
	check(err)

	var numbers []int

	// ioutil.ReadFile returns a byte slice, strings.Split expects a string
	for _, s := range strings.Split(string(data), "\n") {
		i, err := strconv.Atoi(s)
		check(err)

		numbers = append(numbers, i)
	}

	return numbers
}
