package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

const inputFile string = "emulator.input"
const debug bool = false

type Instruction struct {
	Command  string
	Argument int
}

func (instr Instruction) String() string {
	return fmt.Sprintf("%s %d", instr.Command, instr.Argument)
}

type Emulator struct {
	Accumulator        int
	InstructionPointer int
	Instructions       []Instruction
}

// Thanks to the simplified instruction set we can actually solve the halting
// problem, and find infinite loops.
//
// Returns whether an infinite loop is found, along with the value of the
// accumulator before the looping instruction, and the index of the looping
// instruction.
func (emu *Emulator) FindLoops() (loopFound bool, accumulator int, instructionPointer int) {
	emu.Accumulator = 0
	emu.InstructionPointer = 0

	// Track, for each instruction based on its index, whether it was visited previously
	visitedInstructions := make(map[int]bool)

	for emu.InstructionPointer < len(emu.Instructions) {
		// We are about to evaluate this instruction
		visitedInstructions[emu.InstructionPointer] = true

		emu.Evaluate(emu.Instructions[emu.InstructionPointer])

		// If the instruction pointer now points to an instruction we
		// have already evaluated, we've encountered a loop.
		if _, ok := visitedInstructions[emu.InstructionPointer]; ok {
			return true, emu.Accumulator, emu.InstructionPointer
		}
	}

	return false, emu.Accumulator, emu.InstructionPointer
}

func (emu *Emulator) Evaluate(instr Instruction) {
	switch instr.Command {
	case "acc":
		emu.Accumulator += instr.Argument

		if debug {
			fmt.Printf("%s => Accumulator = %d\n", instr, emu.Accumulator)
		}
	case "jmp":
		emu.InstructionPointer += instr.Argument
		if debug {
			fmt.Printf("%s => Instruction pointer = %d\n", instr, emu.InstructionPointer)
		}
		// Skip the auto-increment of the instruction pointer at the end
		return
	case "nop":
		if debug {
			fmt.Printf("%s\n", instr)
		}
		// No-op => nothing to do :)
	default:
		panic(fmt.Sprintf("Invalid instruction: %v %v\n", instr.Command, instr.Argument))
	}

	emu.InstructionPointer += 1
}

func main() {
	taskOne()
	taskTwo()
}

func taskOne() {
	fmt.Println("== Task one ==")

	instructions := loadInstructions()
	emulator := Emulator{Instructions: instructions}

	loopFound, acc, instr := emulator.FindLoops()
	if loopFound {
		fmt.Printf("Loop found at instruction %d, accumulator = %d\n", instr, acc)
	} else {
		fmt.Println("No loop found")
	}
}

func taskTwo() {
	fmt.Println("== Task two ==")

	instructions := loadInstructions()
	mutatedInstructions := make([]Instruction, len(instructions))
	emulator := Emulator{Instructions: instructions}

	// For each NOP and JMP instruction, see if, when changing either to
	// the other (only one mutation at a time), a non-looping program
	// results.
	for i := 0; i < len(instructions); i += 1 {
		// Get fresh copy of instructions to work wtih, so as to only change one instruction.
		copy(mutatedInstructions, instructions)

		switch instructions[i].Command {
		case "acc":
			// ACC is guaranteed to not be corrupted
			continue
		case "jmp":
			// Mutate JMP => NOP
			mutatedInstructions[i].Command = "nop"
			fmt.Printf("Changing JMP => NOP at %d\n", i)
		case "nop":
			// Mutate NOP => JMP
			mutatedInstructions[i].Command = "jmp"
			fmt.Printf("Changing NOP => JMP at %d\n", i)
		}

		emulator.Instructions = mutatedInstructions

		loopFound, acc, instr := emulator.FindLoops()
		if loopFound {
			fmt.Println("Mutated instructions still loop")
		} else {
			fmt.Printf(
				"Success: Mutated instructions terminate normally. Accumulator = %d, Instruction pointer = %d\n",
				acc,
				instr,
			)
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func loadInstructions() []Instruction {
	var instructions []Instruction

	data, err := ioutil.ReadFile(inputFile)
	// Remove trailing newline
	data = data[:len(data)-1]
	check(err)

	for _, line := range strings.Split(string(data), "\n") {
		parts := strings.Split(line, " ")
		if len(parts) != 2 {
			panic(fmt.Sprintf("Invalid input: %v\n", line))
		}

		arg, err := strconv.Atoi(parts[1])
		if err != nil {
			panic(fmt.Sprintf("Invalid argument for instruction: %v %v, %v\n", parts[0], parts[1], err))
		}

		instructions = append(
			instructions,
			Instruction{
				Command:  parts[0],
				Argument: arg,
			},
		)
	}

	return instructions
}
