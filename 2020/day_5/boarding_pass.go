package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"sort"
	"strings"
)

const inputFile string = "boarding_pass.input"

// Rows are 0 .. 127 = 2^0 .. 2^7 - 1
const rowMaxPower int = 6

// Columns are 0 .. 7 = 2^0 .. 2^3 - 1
const columnMaxPower int = 2

type BoardingPass struct {
	RowSpec    string
	ColumnSpec string
}

func (pass BoardingPass) Row() int {
	return boardingPassPatternToDecimal(pass.RowSpec, rowMaxPower, 'F', 'B')
}

func (pass BoardingPass) Column() int {
	return boardingPassPatternToDecimal(pass.ColumnSpec, columnMaxPower, 'L', 'R')
}

func (pass BoardingPass) Id() int {
	return pass.Row()*8 + pass.Column()
}

func boardingPassPatternToDecimal(pattern string, maxPower int, zeroChar rune, oneChar rune) int {
	id := 0
	pow := maxPower

	for _, char := range pattern {
		switch char {
		case zeroChar:
		case oneChar:
			id += int(math.Pow(2, float64(pow)))
		default:
			panic(fmt.Sprintf(
				"Invalid character in specification: '%c' in '%s'\n",
				char, pattern,
			))
		}

		pow -= 1
	}

	return id
}

func main() {
	taskOne()
	taskTwo()
}

func taskOne() {
	fmt.Println("== Task one ==")

	maxId := 0
	passes := loadBoardingPasses()
	for _, pass := range passes {
		if id := pass.Id(); id > maxId {
			maxId = id
		}
	}

	fmt.Printf("Max ID: %d\n", maxId)
}

func taskTwo() {
	fmt.Println("== Task two ==")

	// Track which seats were seen
	seats := make(map[int]bool)
	for i := 0; i < 1024; i++ {
		seats[i] = false
	}

	passes := loadBoardingPasses()
	for _, pass := range passes {
		seats[pass.Id()] = true
	}

	var missingSeats []int
	for seat, seen := range seats {
		if !seen {
			missingSeats = append(missingSeats, seat)
		}
	}

	// Sort missing seats in ascending order
	sort.Ints(missingSeats)
	for _, seat := range missingSeats {
		fmt.Printf("Missing seat: %d\n", seat)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func loadBoardingPasses() []BoardingPass {
	var passes []BoardingPass

	data, err := ioutil.ReadFile(inputFile)
	check(err)
	// Remove trailing newline
	data = data[:len(data)-1]

	for _, line := range strings.Split(string(data), "\n") {
		// Indices 0..6 are row specification, 7..9 column
		// specification.
		// Mind that Go's sub-slicing is a half-open interval,
		// including the start but excluding the stop index.

		passes = append(
			passes,
			BoardingPass{RowSpec: line[:7], ColumnSpec: line[7:]},
		)
	}

	return passes
}
