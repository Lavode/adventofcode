// Package nlp provides various ways to process natural language.
package nlp

import "regexp"

var spelledOutDigits = regexp.MustCompile("(one|two|three|four|five|six|seven|eight|nine)")
var numericalDigits = regexp.MustCompile("(1|2|3|4|5|6|7|8|9)")
var digitMap = map[string]int{
	"1":     1,
	"one":   1,
	"2":     2,
	"two":   2,
	"3":     3,
	"three": 3,
	"4":     4,
	"four":  4,
	"5":     5,
	"five":  5,
	"6":     6,
	"six":   6,
	"7":     7,
	"seven": 7,
	"8":     8,
	"eight": 8,
	"9":     9,
	"nine":  9,
}

// FindDigitFromFront finds the first spelled-out digit starting at the
// front of `input`. If `includeNumeral` is set to true, numerals such as
// `1`, `2` up to `9` are matched as well.
//
// The returned integer is in the range [0, 9]. The boolean indicates whether
// one was found at all.
func FindDigitFromFront(input string, includeNumeral bool) (int, bool) {
	for i := 0; i < len(input); i++ {
		substring := string(input[:i+1])

		if out, ok := findDigit(substring, includeNumeral); ok {
			return out, ok
		}

	}

	return 0, false
}

// FindDigitFromBack finds the first spelled-out digit starting at the
// back of `input`. If `includeNumeral` is set to true, numerals such as
// `1`, `2` up to `9` are matched as well.
//
// The returned integer is in the range [0, 9]. The boolean indicates whether
// one was found at all.
func FindDigitFromBack(input string, includeNumeral bool) (int, bool) {
	for i := len(input) - 1; i >= 0; i-- {
		substring := string(input[i:])

		if out, ok := findDigit(substring, includeNumeral); ok {
			return out, ok
		}
	}

	return 0, false
}

func findDigit(input string, includeNumeral bool) (int, bool) {
	if includeNumeral {
		if match := numericalDigits.FindString(input); match != "" {
			return digitMap[match], true
		}
	}

	if match := spelledOutDigits.FindString(input); match != "" {
		return digitMap[match], true
	}

	return 0, false
}
