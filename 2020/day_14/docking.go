package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

const inputFile string = "docking.input"

func main() {
	taskOne()
	taskTwo()
}

type Instruction struct {
	command string

	// Arguments for MEM command
	address int
	value   int64

	// Arguments for MASK command
	forceOneMask  int64
	forceZeroMask int64
}

type Emulator struct {
	forceOneMask  int64
	forceZeroMask int64
	// Map rather than slice as access is random, sparse, and the highest
	// address is not known beforehand.
	memory map[int]int64
}

func DefaultEmulator() Emulator {
	emu := Emulator{}
	emu.memory = make(map[int]int64)

	return emu
}

func (emu *Emulator) Process(instr Instruction) {
	switch instr.command {
	case "MEM":
		// Compute masked value
		value := (instr.value | emu.forceOneMask) & emu.forceZeroMask
		emu.memory[instr.address] = value
	case "MASK":
		// Update current mask
		emu.forceOneMask = instr.forceOneMask
		emu.forceZeroMask = instr.forceZeroMask
	default:
		panic(fmt.Sprintf("Invalid instruction: %+v\n", instr))
	}
}

func taskOne() {
	fmt.Println("== Task one ==")

	instructions := loadInstructions()
	emulator := DefaultEmulator()

	for _, instr := range instructions {
		emulator.Process(instr)
	}

	// Get sum of values from memory
	var sum int64 = 0
	for _, val := range emulator.memory {
		sum += val
	}

	fmt.Printf("Sum of values in memory: %d\n", sum)
}

func taskTwo() {
	fmt.Println("== Task two ==")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func loadInstructions() []Instruction {
	data, err := ioutil.ReadFile(inputFile)
	// Remove trailing newline
	data = data[:len(data)-1]
	check(err)

	var instructions []Instruction

	maskMatcher := regexp.MustCompile("^mask = ([X01]+)$")
	memMatcher := regexp.MustCompile("^mem\\[([0-9]+)\\] = ([0-9]+)$")

	for _, line := range strings.Split(string(data), "\n") {
		match := maskMatcher.FindStringSubmatch(line)
		if match != nil {
			instructions = append(instructions, parseMaskCommand(match))
		} else {
			match = memMatcher.FindStringSubmatch(line)
			if match == nil {
				panic(fmt.Sprintf("Invalid input line: %q\n", line))
			}

			instructions = append(instructions, parseMemCommand(match))
		}

	}

	return instructions
}

func parseMaskCommand(match []string) Instruction {
	instr := Instruction{command: "MASK"}
	mask := match[1]

	// For each char of the mask, the corresponding digit of subsequent
	// values will be changed as follows:
	// - If the char is X, it will not be changed
	// - If the char is 1, it will be set to 1
	// - If the char is 0, it will be set to 0
	// This can be achieved by creating two masks:
	// - A mask which is 1 everywhere except where the input mask is zero
	// - A mask which is 0 everywhere except where the input mask is 1
	// Then, combining an input with the mask is equal to `input || setOneMask && setTwoMask`

	var forceOneMask int64
	var forceZeroMask int64

	// Mask has a length of 36 characters. We'll start at the most
	// significant digit, and compensate by doubling the value every
	// iteration.
	for i := 0; i < 36; i += 1 {
		forceOneMask = forceOneMask << 1
		forceZeroMask = forceZeroMask << 1

		switch mask[i] {
		case '0':
			// Force corresponding bit to zero:
			// - forceZeroMask must be zero at this digit (do force zero)
			// - forceOneMask must be zero at this digit (do not force one)
		case '1':
			// Force corresponding bit to one:
			// - forceZeroMask must be one at this digit (do not force zero)
			// - forceOneMask must be one at this digit (do force one)
			forceZeroMask += 1
			forceOneMask += 1
		case 'X':
			// Do not force corresponding bit:
			// - forceZeroMask must be one at this digit (do not force zero)
			// - forceOneMask must be zero at this digit (do not force one)
			forceZeroMask += 1
		default:
			panic(fmt.Sprintf("Invalid mask input: %q\n", mask))
		}
	}

	instr.forceOneMask = forceOneMask
	instr.forceZeroMask = forceZeroMask

	return instr
}

func parseMemCommand(match []string) Instruction {
	instr := Instruction{command: "MEM"}

	addr, err := strconv.Atoi(match[1])
	check(err)

	value, err := strconv.Atoi(match[2])
	check(err)

	instr.address = addr
	instr.value = int64(value)

	return instr
}
