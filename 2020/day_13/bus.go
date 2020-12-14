package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

const inputFile string = "bus.input"

func main() {
	taskOne()
	taskTwo()
}

type Bus struct {
	// Bus' ID is also equal to its roundtrip time
	id        int
	inService bool
}

func (bus Bus) NextStart(now int) int {
	// As every bus started at 0, with roundtripe-time `id`, starts will be
	// at `k * id` for integer values of k.

	leftSecondsAgo := now % bus.id
	nextStartIn := bus.id - leftSecondsAgo

	return now + nextStartIn
}

func taskOne() {
	fmt.Println("== Task one ==")

	now, buses := loadBuses()
	nextStart := make(map[int]int)

	fmt.Println("Current TS = ", now)
	// Each bus' numerical ID is equal to its roundtrip time.
	for _, bus := range buses {
		if bus.inService {
			start := bus.NextStart(now)
			fmt.Printf("Bus %d: Next start at %d\n", bus.id, start)
			nextStart[bus.id] = start
		}
	}

	earliestBus, earliestStart := -1, -1

	for bus, start := range nextStart {
		if earliestBus == -1 || nextStart[bus] < earliestStart {
			earliestBus, earliestStart = bus, start
		}
	}

	fmt.Printf("Earliest bus = %d, leaving at %d\n", earliestBus, earliestStart)
	waitingTime := earliestStart - now
	fmt.Printf("Waiting time: %d, Bus ID * waiting time = %d\n", waitingTime, waitingTime*earliestBus)
}

func taskTwo() {
	fmt.Println("== Task two ==")

	_, buses := loadBuses()

	fmt.Println("Solve the following set of linear equations, for integer solutions of `k_i`:")
	for idx, bus := range buses {
		if bus.inService {
			fmt.Printf("t + %d = k_%d * %d\n", idx, idx, bus.id)
		}
	}

	fmt.Println("Or equivalently the following congruence equation system: (Notation: Take == to mean the equivalence relation)")
	for idx, bus := range buses {
		if bus.inService {
			fmt.Printf("t == -%d mod %d\n", idx, bus.id)
		}
	}

	fmt.Println("As the bus IDs are pairwise coprime, you can easily solve this with the CRT")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func loadBuses() (timestamp int, busIDs []Bus) {
	data, err := ioutil.ReadFile(inputFile)
	// Remove trailing newline
	data = data[:len(data)-1]
	check(err)

	var buses []Bus

	// First line is current timestamp
	// Second line is Travel times / IDs of bus lines
	lines := strings.Split(string(data), "\n")
	if len(lines) != 2 {
		panic(fmt.Sprintf("Input has invalid number of lines: %d\n", len(lines)))
	}

	ts, err := strconv.Atoi(lines[0])
	check(err)

	for _, busID := range strings.Split(lines[1], ",") {
		bus := Bus{inService: true}

		if busID == "x" {
			bus.inService = false
		} else {
			id, err := strconv.Atoi(busID)
			check(err)
			bus.id = id
		}

		buses = append(buses, bus)
	}

	return ts, buses
}
