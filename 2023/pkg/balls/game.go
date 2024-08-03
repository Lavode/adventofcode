package balls

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Game is a single execution of the game,  consisting of multiple turns.
type Game struct {
	Id    int64
	turns []Turn
}

// BoundOnBalls returns a lower bound on the number of balls of each colour
// which will have been present in the game.
func (game Game) BoundOnBalls() map[string]int {
	out := make(map[string]int)

	for _, turn := range game.turns {
		for colour, count := range turn.Balls {
			prevCount, ok := out[colour]
			if !ok || count > prevCount {
				out[colour] = count
			}
		}
	}

	return out
}

var gamePattern = regexp.MustCompile("^Game ([0-9]+): (.*)$")

// ParseGame parses a turn from a string representation thereof.
func ParseGame(input string) (Game, error) {
	game := Game{}
	game.turns = make([]Turn, 0)

	matches := gamePattern.FindStringSubmatch(input)
	if matches == nil {
		return game, fmt.Errorf("Invalid input: %s", input)
	}

	id, err := strconv.ParseInt(matches[1], 10, 32)
	if err != nil {
		return game, fmt.Errorf("Invalid ID in game specification: %s, %v", matches[1], err)
	}
	game.Id = id

	serializedTurns := strings.Split(matches[2], "; ")
	if len(serializedTurns) == 0 {
		return game, fmt.Errorf("Invalid turns in game specification: %s", matches[2])
	}

	for _, serializedTurn := range serializedTurns {
		turn, err := ParseTurn(serializedTurn)
		if err != nil {
			return game, fmt.Errorf("Invalid turn in game specification: %s", serializedTurn)
		}

		game.turns = append(game.turns, turn)
	}

	return game, nil
}
