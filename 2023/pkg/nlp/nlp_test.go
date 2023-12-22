package nlp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindDigitFromFront(t *testing.T) {
	input := "underwater2threesevensix9"

	{
		out, ok := FindDigitFromFront(input, false)
		assert.True(t, ok)
		assert.Equal(t, 3, out)
	}

	{
		out, ok := FindDigitFromFront(input, true)
		assert.True(t, ok)
		assert.Equal(t, 2, out)
	}

	{
		out, ok := FindDigitFromFront("naught to see 5; move along", false)
		assert.False(t, ok)
		assert.Equal(t, 0, out)
	}

	{
		out, ok := FindDigitFromFront("naught to see; move along", true)
		assert.False(t, ok)
		assert.Equal(t, 0, out)
	}
}

func TestFindDigitFromBack(t *testing.T) {
	input := "underwater2threesevensix9"

	{
		out, ok := FindDigitFromBack(input, false)
		assert.True(t, ok)
		assert.Equal(t, 6, out)
	}

	{
		out, ok := FindDigitFromBack(input, true)
		assert.True(t, ok)
		assert.Equal(t, 9, out)
	}

	{
		out, ok := FindDigitFromBack("naught to see 5; move along", false)
		assert.False(t, ok)
		assert.Equal(t, 0, out)
	}

	{
		out, ok := FindDigitFromBack("naught to see; move along", true)
		assert.False(t, ok)
		assert.Equal(t, 0, out)
	}
}
