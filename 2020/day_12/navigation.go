package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

const inputFile string = "navigation.input"

type Instruction struct {
	command  string
	argument int
}

type Ship struct {
	// east-west (east positive)
	longitude int
	// north-south (north positive)
	latitude int

	// Orientation. 0 = North, 90 = East, 180 = South, 270 = West
	heading int

	waypoint Waypoint
}

// Waypoint coordinates are always relative to ship, ie really an offset in lat
// / lon.
type Waypoint struct {
	// east-west (east positive)
	longitude int
	// north-south (north positive)
	latitude int
}

func defaultWaypoint() Waypoint {
	return Waypoint{longitude: 10, latitude: 1}
}

func (waypoint *Waypoint) MoveNorth(x int) {
	waypoint.latitude += x
}

func (waypoint *Waypoint) MoveSouth(x int) {
	waypoint.latitude -= x
}

func (waypoint *Waypoint) MoveEast(x int) {
	waypoint.longitude += x
}

func (waypoint *Waypoint) MoveWest(x int) {
	waypoint.longitude -= x
}

func (waypoint *Waypoint) RotateLeft(deg int) {
	// Rotate left (ccw) deg / 90 times (integer division, will round down)
	for i := 0; i < deg/90; i += 1 {
		// Each ccw rotation by 90 degrees is equivalent to
		// transforming (x, y) into (x', y') as:
		// x' := -y
		// y' := x
		waypoint.longitude, waypoint.latitude = -waypoint.latitude, waypoint.longitude
	}
}

func (waypoint *Waypoint) RotateRight(deg int) {
	// Rotate right (cw) deg / 90 times (integer division, will round down)
	for i := 0; i < deg/90; i += 1 {
		// Each cw rotation by 90 degrees is equivalent to
		// transforming (x, y) into (x', y') as:
		// x' := y
		// y' := -x
		waypoint.longitude, waypoint.latitude = waypoint.latitude, -waypoint.longitude
	}
}

func defaultShip() Ship {
	// Default ship heading is east
	return Ship{heading: 90, waypoint: defaultWaypoint()}
}

// Instructions move ship directly
func ProcessDirectMovement(inst Instruction, ship *Ship) {
	switch inst.command {
	case "N":
		ship.MoveNorth(inst.argument)
	case "S":
		ship.MoveSouth(inst.argument)
	case "E":
		ship.MoveEast(inst.argument)
	case "W":
		ship.MoveWest(inst.argument)
	case "F":
		ship.MoveForward(inst.argument)
	case "L":
		ship.TurnLeft(inst.argument)
	case "R":
		ship.TurnRight(inst.argument)
	default:
		panic(fmt.Sprintf("Unknown command: %s\n", inst.command))
	}
}

// Instructions move ship by means of waypoint
func ProcessIndirectMovement(inst Instruction, ship *Ship) {
	switch inst.command {
	case "N":
		ship.waypoint.MoveNorth(inst.argument)
	case "S":
		ship.waypoint.MoveSouth(inst.argument)
	case "E":
		ship.waypoint.MoveEast(inst.argument)
	case "W":
		ship.waypoint.MoveWest(inst.argument)
	case "F":
		ship.MoveToWaypoint(inst.argument)
	case "L":
		ship.waypoint.RotateLeft(inst.argument)
	case "R":
		ship.waypoint.RotateRight(inst.argument)
	default:
		panic(fmt.Sprintf("Unknown command: %s\n", inst.command))
	}
}

func (ship *Ship) MoveNorth(x int) {
	ship.latitude += x
}

func (ship *Ship) MoveSouth(x int) {
	ship.latitude -= x
}

func (ship *Ship) MoveEast(x int) {
	ship.longitude += x
}

func (ship *Ship) MoveWest(x int) {
	ship.longitude -= x
}

func (ship *Ship) MoveForward(x int) {
	switch ship.heading {
	case 0:
		ship.MoveNorth(x)
	case 90:
		ship.MoveEast(x)
	case 180:
		ship.MoveSouth(x)
	case 270:
		ship.MoveWest(x)
	default:
		fmt.Printf("Unable to move towards current heading: %d\n", ship.heading)
	}
}

func (ship *Ship) MoveToWaypoint(x int) {
	// Waypoint is relative to ship, if ship moves so does the waypoint. As
	// such we treat its lat/lon as offsets directly, and never change
	// them.

	for i := 0; i < x; i += 1 {
		ship.latitude += ship.waypoint.latitude
		ship.longitude += ship.waypoint.longitude
	}
}

func (ship *Ship) TurnLeft(x int) {
	ship.heading = mod(ship.heading-x, 360)
}

func (ship *Ship) TurnRight(x int) {
	ship.heading = mod(ship.heading+x, 360)
}

// Manhattan distance of ship from origin (0, 0)
func (ship Ship) distance() int {
	return abs(ship.longitude) + abs(ship.latitude)
}

func abs(x int) int {
	if x < 0 {
		return -x
	} else {
		return x
	}
}

// Go's built-in % operation is a remainder, not a mathematical modulus. As an
// example, a % b < 0 if a < 0.
// This is modulo as you know and love it from mathematics.
func mod(a, b int) int {
	m := a % b

	if a < 0 {
		m += abs(b)
	}

	return m
}

func main() {
	taskOne()
	taskTwo()
}

func taskOne() {
	fmt.Println("== Task one ==")

	instructions := loadInstructions()
	ship := defaultShip()

	for _, instr := range instructions {
		ProcessDirectMovement(instr, &ship)
	}

	fmt.Printf("Ship's current position: Lat: %d, Lon: %d, Heading: %d\n", ship.latitude, ship.longitude, ship.heading)
	fmt.Printf("Distance from origin: %d\n", ship.distance())
}

func taskTwo() {
	fmt.Println("== Task two ==")

	instructions := loadInstructions()
	ship := defaultShip()

	for _, instr := range instructions {
		ProcessIndirectMovement(instr, &ship)
	}

	fmt.Printf("Ship's current position: Lat: %d, Lon: %d, Heading: %d\n", ship.latitude, ship.longitude, ship.heading)
	fmt.Printf("Distance from origin: %d\n", ship.distance())
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

	matcher := regexp.MustCompile("^([A-Z])([0-9]+)$")

	// ioutil.ReadFile returns a byte slice, strings.Split expects a string
	for _, line := range strings.Split(string(data), "\n") {
		match := matcher.FindStringSubmatch(line)
		if match == nil || len(match) != 3 {
			panic(fmt.Sprintf("Invalid input line: %s\n", line))
		}

		instruction := match[1]
		arg, err := strconv.Atoi(match[2])
		check(err)

		switch instruction {
		case "N", "S", "E", "W", "F":
		case "L", "R":
			if arg < 0 || arg > 360 || arg%90 != 0 {
				panic(fmt.Sprintf("Invalid argument for turn left/right command: %d\n", arg))
			}
		default:
			panic(fmt.Sprintf("Invalid instruction '%s' in line %s\n", instruction, line))
		}

		instructions = append(instructions, Instruction{command: instruction, argument: arg})
	}

	return instructions
}
