package balls

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTurn(t *testing.T) {
	{
		turn, err := ParseTurn("3 blue, 4 red")
		assert.NoError(t, err)
		assert.Equal(t, map[string]int{"blue": 3, "red": 4}, turn.Balls)
	}

	{
		turn, err := ParseTurn("5 green, 1 red")
		assert.NoError(t, err)
		assert.Equal(t, map[string]int{"green": 5, "red": 1}, turn.Balls)

	}
}
