package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// const inputFile = "input.test.txt"

const inputFile = "input.txt"

func main() {
	debugMessages := getInput()
	getRates(debugMessages)

	getLifeSupportRating(debugMessages)
}

func getLifeSupportRating(debugMessages []string) {
	oxygenRatingCandidates := make([]string, len(debugMessages))
	copy(oxygenRatingCandidates, debugMessages)

	co2ScrubberRatingCandidates := make([]string, len(debugMessages))
	copy(co2ScrubberRatingCandidates, debugMessages)

	for i := 0; len(oxygenRatingCandidates) > 1; i += 1 {
		// Statistic are always *on the remaining subset
		stats := getBitStatistics(oxygenRatingCandidates)

		// If no unique match is found, this will overflow and cause a
		// runtime error.

		// In case of equality we'll keep the ones with a 1
		mostCommonBit := '1'
		if stats[i][0] > stats[i][1] {
			mostCommonBit = '0'
		}

		// Now filter the list of candidates. We'll keep the oxygen
		// generator ratings which share the most common value at the
		// current position, and the co2 scrubber ratings whose bit is
		// equal to the least common one.
		// We do the filtering in-place, because we're lazy.
		k := 0
		for _, candidate := range oxygenRatingCandidates {
			// the byte() conversion is safe as we're workin with
			// ASCII files containing only 0s (0x30) and 1s (0x31)
			if candidate[i] == byte(mostCommonBit) {
				// Keep it
				oxygenRatingCandidates[k] = candidate
				k += 1
			}
		}
		// Truncate appropriately
		oxygenRatingCandidates = oxygenRatingCandidates[:k]

	}

	// Aaand the same thing for the CO2 scrubber
	// TODO extract into method
	for i := 0; len(co2ScrubberRatingCandidates) > 1; i += 1 {
		// Statistic are always *on the remaining subset
		stats := getBitStatistics(co2ScrubberRatingCandidates)

		// If no unique match is found, this will overflow and cause a
		// runtime error.

		// In caes of equality we'll keep the ones with a 0, so will set the MCB to 1
		mostCommonBit := '1'
		if stats[i][0] > stats[i][1] {
			mostCommonBit = '0'
		}

		// Now filter the list of candidates. We'll keep the co2 scrubber
		// generator ratings which share the most common value at the
		// current position, and the co2 scrubber ratings whose bit is
		// equal to the least common one.
		// We do the filtering in-place, because we're lazy.
		k := 0
		for _, candidate := range co2ScrubberRatingCandidates {
			// the byte() conversion is safe as we're workin with
			// ASCII files containing only 0s (0x30) and 1s (0x31)
			if candidate[i] != byte(mostCommonBit) {
				// Keep it
				co2ScrubberRatingCandidates[k] = candidate
				k += 1
			}
		}
		// Truncate appropriately
		co2ScrubberRatingCandidates = co2ScrubberRatingCandidates[:k]
	}

	fmt.Println(oxygenRatingCandidates)

	oxygenRating, err := strconv.ParseInt(oxygenRatingCandidates[0], 2, 16)
	if err != nil {
		log.Fatalf("Unable to convert oxygen rate to int: %v\n", err)
	}

	co2ScrubberRating, err := strconv.ParseInt(co2ScrubberRatingCandidates[0], 2, 16)
	if err != nil {
		log.Fatalf("Unable to convert CO2 scrubber rate to int: %v\n", err)
	}

	fmt.Printf("Oxygen rating = %d, CO2 scrubber = %d, product = %d\n", oxygenRating, co2ScrubberRating, oxygenRating*co2ScrubberRating)
}

func getRates(debugMessages []string) {
	// For each 'bit', the most common one of all inputs will form a bit of
	// the gamma rate, the least common one the epsilon rate.
	gammaRateString := ""
	epsilonRateString := ""

	stats := getBitStatistics(debugMessages)
	for i := 0; i < len(debugMessages[0]); i += 1 {
		mostCommonBit := '0'
		if stats[i][1] > stats[i][0] {
			mostCommonBit = '1'
		}

		if mostCommonBit == '1' {
			gammaRateString = gammaRateString + "1"
			epsilonRateString = epsilonRateString + "0"
		} else {
			gammaRateString = gammaRateString + "0"
			epsilonRateString = epsilonRateString + "1"
		}
	}

	// Interpret as a big-endian base-two number
	gammaRate, err := strconv.ParseInt(gammaRateString, 2, 16)
	if err != nil {
		log.Fatalf("Unable to convert gamma rate to int: %v\n", err)
	}
	epsilonRate, err := strconv.ParseInt(epsilonRateString, 2, 16)
	if err != nil {
		log.Fatalf("Unable to convert gamma rate to int: %v\n", err)
	}

	fmt.Printf("Gamma rate = %d\nEpsilon rate = %d\n", gammaRate, epsilonRate)
	fmt.Printf("gamme * epsilon = %d\n", gammaRate*epsilonRate)
}

func getBitStatistics(debugMessages []string) map[int](map[int]int) {
	// First layer of map is the zero-indexed position in the bit strings.
	// Second layer of map will be, for the corresponding position, the
	// number of zero- and one- bits respectively.
	out := make(map[int](map[int]int))

	for i := 0; i < len(debugMessages[0]); i += 1 {
		out[i] = make(map[int]int)
	}

	for _, str := range debugMessages {
		for i, chr := range str {
			if chr == '0' {
				out[i][0] += 1
			} else {
				out[i][1] += 1
			}
		}
	}

	return out
}

func getInputLines() []string {
	content, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Error reading file: %v\n", err)
	}

	lines := strings.Split(string(content), "\n")
	// Drop trailing newline
	lines = lines[:len(lines)-1]

	return lines
}

func getInput() []string {
	return getInputLines()
}
