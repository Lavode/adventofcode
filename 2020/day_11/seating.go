package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

const inputFile string = "seating.input"
const debug bool = false

type Field struct {
	kind     string
	occupied bool
}

func (field *Field) IsSeat() bool {
	return field.kind == "seat"
}

func (field *Field) IsFloor() bool {
	return field.kind == "floor"
}

func (field Field) String() string {
	switch field.kind {
	case "floor":
		return "."
	case "seat":
		if field.occupied {
			return "#"
		} else {
			return "L"
		}
	default:
		panic(fmt.Sprintf("Invalid field type: %s\n", field.kind))
	}
}

type Parameters struct {
	// If the number of occupied neighbours is >= this parameter, a
	// previously occupied seat will be freed.
	FreeSeatThreshold int

	// If the number of occupied neighbours is <= this parameter, a
	// previously freed seat will be occupied.
	OccupySeatThreshold int

	// A function which, for a given row/column coordinate, returns the
	// fields which are considered a neighbour of this field.
	GetNeighbours func(row int, col int, board *Board) []Field
}

// Return the eight adjacent neighbours of the field at (row, col). If the
// field is at the border, fewer than eight neighbours might be returned.
func getAdjacentNeighbours(row int, col int, board *Board) []Field {
	// A limit of 1 ensures we only consider directly adjacent tiles.
	return getCardinalNeighbours(row, col, board, 1)
}

// Return the eight neighbours of the field at (row, col) which are closest in
// each of the eight directions. If the field is at the border, fewer than
// eight neighbours might be returned.
// Unlike `getAdjacentNeighbours`, such neighbours may be more than than one
// field away.
func getLineOfSightNeighbours(row int, col int, board *Board) []Field {
	return getCardinalNeighbours(row, col, board, -1)
}

// Return non-floor neighbours in each of the eight cardinal directions, up to
// a Chebyshev distance of `limit`.
func getCardinalNeighbours(row int, col int, board *Board, limit int) []Field {
	var neighbours []Field

	for x := -1; x <= 1; x += 1 {
		for y := -1; y <= 1; y += 1 {
			if x == 0 && y == 0 {
				// A field isn't its own neighbour
				continue
			}

			neighbour := getNeighbourInLine(row, col, board, x, y, limit)
			if neighbour != nil {
				neighbours = append(neighbours, *neighbour)
			}
		}
	}

	return neighbours
}

// Get closest seat neighbour in direction indicated by (dX, dY) vector.
//
// `limit` limits the number of tiles which to search in the given direction. A
// limit of `-1` indicates no limit.
//
// If no seat is found in the given direction within the confines of the field,
// nil is returned.
func getNeighbourInLine(row int, col int, board *Board, dX int, dY int, limit int) *Field {
	if debug {
		fmt.Printf("Finding neighbour of (%d, %d) in direction (%d, %d) with limit %d\n", row, col, dX, dY, limit)
	}

	noDistanceLimit := limit == -1

	// Using Chebyshev distance to limit how far we travel
	for step := 1; step <= limit || noDistanceLimit; step += 1 {
		x := row + step*dX
		y := col + step*dY

		neighbour := board.GetField(x, y)
		if neighbour == nil {
			// We've escaped the confines of the board
			return nil
		} else if neighbour.IsSeat() {
			// Found first seat in the given direction
			return neighbour
		}
	}

	// No seat found in the requested direction and sight limit.
	return nil
}

type Board struct {
	fields [][]Field
}

func (board Board) String() string {
	out := ""
	for _, fields := range board.fields {
		for _, field := range fields {
			out += field.String()
		}
		out += "\n"
	}

	return out
}

// Returns the field at (row, col), or `nil` if out of bounds.
func (board *Board) GetField(row int, col int) *Field {
	if row < 0 || col < 0 || row > board.RowCount()-1 || col > board.ColumnCount()-1 {
		return nil
	} else {
		return &board.fields[row][col]
	}
}

func (board *Board) OccupiedSeatsCount() int {
	count := 0

	for _, row := range board.fields {
		for _, field := range row {
			if field.IsSeat() && field.occupied {
				count += 1
			}
		}
	}

	return count
}

func (board *Board) RowCount() int {
	return len(board.fields)
}

func (board *Board) ColumnCount() int {
	if board.RowCount() == 0 {
		return 0
	} else {
		return len(board.fields[0])
	}
}

func (board *Board) Step(params Parameters) bool {
	boardChanged := false
	var newFields [][]Field = make([][]Field, len(board.fields))

	for row, fields := range board.fields {
		newFields[row] = make([]Field, len(fields))

		for col, oldField := range fields {
			if oldField.IsSeat() {
				occupiedCount := 0
				for _, neighbour := range params.GetNeighbours(row, col, board) {
					if neighbour.IsSeat() && neighbour.occupied {
						occupiedCount += 1
					}
				}

				if occupiedCount <= params.OccupySeatThreshold {
					// Seat becomes occupied
					newFields[row][col] = Field{kind: "seat", occupied: true}

					if !oldField.occupied {
						boardChanged = true
					}
				} else if occupiedCount >= params.FreeSeatThreshold {
					// Seat becomes empty
					newFields[row][col] = Field{kind: "seat", occupied: false}

					if oldField.occupied {
						boardChanged = true
					}
				} else {
					// More than 0 but less than 4 neighbours => Seat stays as-is
					newFields[row][col] = oldField
				}
			} else {
				// Floors never change
				newFields[row][col] = Field{kind: "floor", occupied: false}
			}
		}
	}

	board.fields = newFields

	return boardChanged
}

func main() {
	taskOne()
	taskTwo()
}

func taskOne() {
	fmt.Println("== Task one ==")

	board := loadBoard()

	params := Parameters{
		FreeSeatThreshold:   4,
		OccupySeatThreshold: 0,
		GetNeighbours:       getAdjacentNeighbours,
	}

	runUntilStable(&board, params)

	fmt.Printf("Occupied seat count of stable board: %d\n", board.OccupiedSeatsCount())
}

func taskTwo() {
	fmt.Println("== Task two ==")

	board := loadBoard()

	params := Parameters{
		FreeSeatThreshold:   5,
		OccupySeatThreshold: 0,
		GetNeighbours:       getLineOfSightNeighbours,
	}

	runUntilStable(&board, params)

	fmt.Printf("Occupied seat count of stable board: %d\n", board.OccupiedSeatsCount())
}

func runUntilStable(board *Board, params Parameters) {
	for i := 0; ; i += 1 {
		fmt.Printf("\nStep: %d\n", i)
		fmt.Println(board)
		boardChanged := board.Step(params)
		if boardChanged {
			fmt.Printf("Change observed\n")
		} else {
			fmt.Printf("No change observed => Board is stable\n")
			break
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func loadBoard() Board {
	data, err := ioutil.ReadFile(inputFile)
	// Remove trailing newline
	data = data[:len(data)-1]
	check(err)

	var board Board

	// ioutil.ReadFile returns a byte slice, strings.Split expects a string
	for _, line := range strings.Split(string(data), "\n") {
		var fields []Field

		for _, thing := range line {
			switch thing {
			case 'L':
				fields = append(fields, Field{kind: "seat"})
			case '.':
				fields = append(fields, Field{kind: "floor"})
			default:
				panic(fmt.Sprintf("Invalid input '%c' in line: %s\n", thing, line))
			}
		}

		board.fields = append(board.fields, fields)
	}

	return board
}
