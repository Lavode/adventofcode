package main

import "fmt"
import "io/ioutil"
import "strings"

const inputFile string = "avoiding_trees.input"
const tree rune = '#'
const free rune = '.'

func main() {
	taskOne()
	taskTwo()
}

func taskOne() {
	fmt.Println("== Task one ==")
	treesHit := checkTreesHit(3, 1)

	fmt.Printf("Hit %d trees\n", treesHit)
}

func taskTwo() {
	fmt.Println("== Task two ==")

	treesHit1 := checkTreesHit(1, 1)
	fmt.Printf("Slope 1/1: Hit %d trees\n", treesHit1)

	treesHit2 := checkTreesHit(3, 1)
	fmt.Printf("Slope 3/1: Hit %d trees\n", treesHit2)

	treesHit3 := checkTreesHit(5, 1)
	fmt.Printf("Slope 5/1: Hit %d trees\n", treesHit3)

	treesHit4 := checkTreesHit(7, 1)
	fmt.Printf("Slope 7/1: Hit %d trees\n", treesHit4)

	treesHit5 := checkTreesHit(1, 2)
	fmt.Printf("Slope 1/2: Hit %d trees\n", treesHit5)

	fmt.Printf("Product of trees hit: %d\n", treesHit1*treesHit2*treesHit3*treesHit4*treesHit5)
}

func checkTreesHit(dX int, dY int) int {
	// x is horizontal, y vertical coordinate.
	// Mind that first index of `landscape` is the *vertical* coordinate,
	// ie y.
	landscape := loadLandscape()
	var x, y int = 0, 0

	var width int = len(landscape[0])
	var height int = len(landscape)

	var treesHit int = 0

	for y < height {
		// fmt.Printf("x = %d, y = %d, Field = %c\n", x, y, landscape[y][x])

		if landscape[y][x] == tree {
			treesHit += 1
		}

		// As each row has `width` entries, indices are from 0 to width - 1.
		// Sum modulo width ensures we wrap around.
		x = (x + dX) % width
		y += dY
	}

	return treesHit
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func loadLandscape() [][]rune {
	var landscape [][]rune

	data, err := ioutil.ReadFile(inputFile)
	// Remove trailing newline
	data = data[:len(data)-1]
	check(err)

	for _, line := range strings.Split(string(data), "\n") {
		var row []rune
		for _, thing := range line {
			row = append(row, thing)
		}
		landscape = append(landscape, row)
	}

	return landscape
}
