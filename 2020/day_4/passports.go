package main

import "fmt"
import "io/ioutil"
import "regexp"
import "strconv"
import "strings"

const inputFile string = "passports.input"

type Passport struct {
	byr string
	iyr string
	eyr string
	hgt string
	hcl string
	ecl string
	pid string
	cid string
}

func (passport Passport) requiredFieldsPresent() bool {
	return passport.byr != "" &&
		passport.iyr != "" &&
		passport.eyr != "" &&
		passport.hgt != "" &&
		passport.hcl != "" &&
		passport.ecl != "" &&
		passport.pid != ""
}

func (passport Passport) isValid() bool {
	// Birth year: Four digits, numeric, 1920 <= x <= 2002
	byr, err := strconv.Atoi(passport.byr)
	if err != nil || len(passport.byr) != 4 || byr < 1920 || byr > 2002 {
		return false
	}

	// Issue year: Four digits, numeric, 2010 <= x <= 2020
	iyr, err := strconv.Atoi(passport.iyr)
	if err != nil || len(passport.iyr) != 4 || iyr < 2010 || iyr > 2020 {
		return false
	}

	// Expiration year: Four digits, numeric, 2020 <= x <= 2030
	eyr, err := strconv.Atoi(passport.eyr)
	if err != nil || len(passport.eyr) != 4 || eyr < 2020 || eyr > 2030 {
		return false
	}

	// Height: Number, followed by 'cm' or 'in'.
	//   For cm: 150 <= x <= 193
	//   For in: 59 <= x <= 76
	heightMatcher := regexp.MustCompile(`^(\d+)(cm|in)$`)
	match := heightMatcher.FindStringSubmatch(passport.hgt)
	if match == nil {
		return false
	}

	height, err := strconv.Atoi(match[1])
	if err != nil {
		return false
	}

	unit := match[2]
	switch unit {
	case "cm":
		if height < 150 || height > 193 {
			return false
		}
	case "in":
		if height < 59 || height > 76 {
			return false
		}
	default:
		// Should not happen at all, regex above should catch it
		panic(fmt.Sprintf("Invalid unit: %v\n", unit))
	}

	// Hair colour: Six-digit hexadecimal number, prefixed with '#'
	hairColourMatcher := regexp.MustCompile("^#([0-9a-f]{6})$")
	match = hairColourMatcher.FindStringSubmatch(passport.hcl)
	if match == nil {
		return false
	}

	// Eye colour: One of (amb, blu, brn, gry, grn, hzl, oth)
	validEyeColours := map[string]bool{
		"amb": true,
		"blu": true,
		"brn": true,
		"gry": true,
		"grn": true,
		"hzl": true,
		"oth": true,
	}
	_, ok := validEyeColours[passport.ecl]
	if !ok {
		return false
	}

	// Passport ID: Nine-digit number
	passportIDMatcher := regexp.MustCompile("^([0-9]{9})$")
	match = passportIDMatcher.FindStringSubmatch(passport.pid)
	if match == nil {
		return false
	}

	// Country ID is not validated, missing or not

	return true
}

func main() {
	taskOne()
	taskTwo()
}

func taskOne() {
	fmt.Println("== Task one ==")

	validPassports := 0
	passports := loadPassports()

	for _, passport := range passports {
		if passport.requiredFieldsPresent() {
			validPassports += 1
		}
	}

	fmt.Printf("Total passports: %d. Valid: %d\n", len(passports), validPassports)
}

func taskTwo() {
	fmt.Println("== Task two ==")

	validPassports := 0
	passports := loadPassports()

	for _, passport := range passports {
		if passport.isValid() {
			validPassports += 1
		}
	}

	fmt.Printf("Total passports: %d. Valid: %d\n", len(passports), validPassports)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func loadPassports() []Passport {
	var passports []Passport

	data, err := ioutil.ReadFile(inputFile)
	// Don't remove trailing new line, as it indicates the end of the last passport entry
	check(err)

	var passport Passport

	for _, line := range strings.Split(string(data), "\n") {
		if line == "" {
			// Individual passports are separated by empty lines
			passports = append(passports, passport)
			passport = Passport{}
		} else {
			// Non-empty lines contain key:value pairs for current
			// passport, each separated by a space
			for _, pair := range strings.Split(line, " ") {
				kv := strings.Split(pair, ":")

				switch kv[0] {
				case "byr":
					passport.byr = kv[1]
				case "iyr":
					passport.iyr = kv[1]
				case "eyr":
					passport.eyr = kv[1]
				case "hgt":
					passport.hgt = kv[1]
				case "hcl":
					passport.hcl = kv[1]
				case "ecl":
					passport.ecl = kv[1]
				case "pid":
					passport.pid = kv[1]
				case "cid":
					passport.cid = kv[1]
				default:
					panic(fmt.Sprintf("Invalid passport field: %v\n", kv[0]))
				}
			}
		}
	}

	return passports
}
