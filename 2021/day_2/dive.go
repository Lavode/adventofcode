package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// const inputFile = "input.test.txt"
const inputFile = "input.txt"

// Movement direction for sub commands
const (
	Down = iota
	Up
	Forward
)

type Command struct {
	Direction int
	Amount    int
}

func main() {
	commands := getInput()
	calculateSimplePosition(commands)
	calculateComplexPosition(commands)
}

func calculateSimplePosition(commands []Command) {
	// As per task 1

	depth := 0
	// Horizontal distance
	distance := 0

	for _, cmd := range commands {
		if cmd.Direction == Up {
			depth -= cmd.Amount
		} else if cmd.Direction == Down {
			depth += cmd.Amount
		} else if cmd.Direction == Forward {
			distance += cmd.Amount
		}
	}

	fmt.Printf(
		"Final depth = %d, distance = %d, product = %d\n",
		depth,
		distance,
		depth*distance,
	)
}

func calculateComplexPosition(commands []Command) {
	// As per task 2

	depth := 0
	// Positive indicates we are aiming towards the sea floor, negative
	// towards the surface.
	// When moving forward by `x`, the depth will also change by `aim * x`.
	aim := 0
	// Horizontal distance
	distance := 0

	for _, cmd := range commands {
		if cmd.Direction == Up {
			aim -= cmd.Amount
		} else if cmd.Direction == Down {
			aim += cmd.Amount
		} else if cmd.Direction == Forward {
			distance += cmd.Amount
			depth += aim * cmd.Amount
		}

		fmt.Printf("Processed command: %v\n", cmd)
		fmt.Printf("Depth = %d, distance = %d, aim = %d\n", depth, distance, aim)
	}

	fmt.Printf(
		"Final depth = %d, distance = %d, product = %d\n",
		depth,
		distance,
		depth*distance,
	)
}

func getInputLines() []string {
	content, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Error reading file: %v\n", err)
	}

	lines := strings.Split(string(content), "\n")
	// Drop trailing newline
	lines = lines[:len(lines)-1]

	return lines
}

func getInput() []Command {
	commands := make([]Command, 0)

	pattern := regexp.MustCompile("^(up|down|forward) ([0-9]+)$")

	for _, line := range getInputLines() {
		match := pattern.FindStringSubmatch(line)
		if len(match) != 3 {
			log.Fatalf("Line did not match pattern: %v\n", line)
		}

		cmdType := match[1]
		cmdArg, err := strconv.Atoi(match[2])
		if err != nil {
			log.Fatalf("Unable to parse command, invalid command argument: %v\n", line)
		}

		if cmdType == "up" {
			commands = append(commands, Command{Up, cmdArg})
		} else if cmdType == "down" {
			commands = append(commands, Command{Down, cmdArg})
		} else if cmdType == "forward" {
			commands = append(commands, Command{Forward, cmdArg})
		} else {
			log.Fatalf("Unable to parse command, invalid command type: %v\n", line)
		}
	}

	return commands
}
