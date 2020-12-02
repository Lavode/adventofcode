package main

import "fmt"
import "io/ioutil"
import "regexp"
import "strconv"
import "strings"

const inputFile string = "invalid_passwords.input"

type Password struct {
	SpecFieldStart int
	SpecFieldEnd   int
	Char           string
	Password       string
}

// Validation rules as per part 1:
// Password must contain given character between SpecFieldStart and
// SpecFieldEnd many times.
func (password Password) validate() bool {
	count := strings.Count(password.Password, password.Char)
	return count >= password.SpecFieldStart && count <= password.SpecFieldEnd
}

// Validation rules as per part 2:
// Exactly one of (one-indexed) SpecFieldStart OR SpecFieldEnd characters of
// password must be given character.
func (password Password) validateAlt() bool {
	matchingChars := 0
	// Abusing the fact that passwords are ASCII only, else we'd have to work with runes.
	char := password.Char[0]

	// MinCount and MaxCount are 1-based indices
	if password.Password[password.SpecFieldStart-1] == char {
		matchingChars += 1
	}

	if password.Password[password.SpecFieldEnd-1] == char {
		matchingChars += 1
	}

	return matchingChars == 1
}

func main() {
	taskOne()
	taskTwo()
}

func taskOne() {
	fmt.Println("== Task one ==")

	validCount := 0
	invalidCount := 0

	passwords := getPasswords()
	for _, password := range passwords {
		if password.validate() {
			validCount += 1
		} else {
			invalidCount += 1
		}
	}

	fmt.Printf(
		"Checked %d passwords.\n\tValid: %d\n\tInvalid: %d\n\t",
		validCount+invalidCount,
		validCount,
		invalidCount,
	)
}

func taskTwo() {
	fmt.Println("== Task two ==")

	validCount := 0
	invalidCount := 0

	passwords := getPasswords()
	for _, password := range passwords {
		if password.validateAlt() {
			validCount += 1
		} else {
			invalidCount += 1
		}
	}

	fmt.Printf(
		"Checked %d passwords.\n\tValid: %d\n\tInvalid: %d\n\t",
		validCount+invalidCount,
		validCount,
		invalidCount,
	)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getPasswords() []Password {
	data, err := ioutil.ReadFile(inputFile)
	// Remove trailing newline
	data = data[:len(data)-1]
	check(err)

	var passwords []Password

	matcher := regexp.MustCompile(`^(\d+)-(\d+) ([a-z]): ([a-z]+)$`)
	// ioutil.ReadFile returns a byte slice, strings.Split expects a string
	for _, s := range strings.Split(string(data), "\n") {
		match := matcher.FindStringSubmatch(s)
		if match == nil {
			panic(fmt.Sprintf("Line did not match pattern: %v\n", s))
		}

		minCount, err := strconv.Atoi(match[1])
		check(err)

		maxCount, err := strconv.Atoi(match[2])
		check(err)

		char := match[3]

		password := match[4]

		passwords = append(
			passwords,
			Password{
				SpecFieldStart: minCount,
				SpecFieldEnd:   maxCount,
				Char:           char,
				Password:       password,
			},
		)
	}

	return passwords
}
