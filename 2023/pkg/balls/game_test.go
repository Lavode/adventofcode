package balls

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseGame(t *testing.T) {
	{
		game, err := ParseGame("Game 2: 1 blue, 2 green; 3 green, 4 blue, 1 red; 1 green, 1 blue")
		assert.NoError(t, err)
		t1 := Turn{Balls: map[string]int{"blue": 1, "green": 2}}
		t2 := Turn{Balls: map[string]int{"blue": 4, "green": 3, "red": 1}}
		t3 := Turn{Balls: map[string]int{"blue": 1, "green": 1}}

		assert.ElementsMatch(t, []Turn{t1, t2, t3}, game.turns)
	}
}

func TestBoundOnBalls(t *testing.T) {
	game, err := ParseGame("Game 2: 1 blue, 2 green; 3 green, 4 blue, 1 red; 1 green, 1 blue")
	assert.NoError(t, err)

	bounds := game.BoundOnBalls()
	assert.Equal(t, 4, bounds["blue"])
	assert.Equal(t, 1, bounds["red"])
	assert.Equal(t, 3, bounds["green"])
}
