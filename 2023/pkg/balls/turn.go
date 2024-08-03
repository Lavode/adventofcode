package balls

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var quantityPattern = regexp.MustCompile("^([0-9]+) (red|blue|green)$")

// Turn is a single turn of a game, within which various amounts
// of various colours of balls are revealed.
type Turn struct {
	Balls map[string]int
}

// ParseTurn parses a turn from a string representation thereof.
func ParseTurn(input string) (Turn, error) {
	turn := Turn{}
	turn.Balls = make(map[string]int)

	quantities := strings.Split(input, ", ")
	if len(quantities) == 0 {
		return turn, fmt.Errorf("Invalid input: %s", input)
	}

	for _, quantity := range quantities {
		matches := quantityPattern.FindStringSubmatch(quantity)
		if matches == nil {
			return turn, fmt.Errorf("Invalid input part: %s", quantity)
		}

		count, err := strconv.ParseInt(matches[1], 10, 32)
		if err != nil {
			return turn, fmt.Errorf("Invalid count in turn specification: %s, %v", matches[0], err)
		}

		color := matches[2]

		turn.Balls[color] = int(count)
	}

	return turn, nil
}
