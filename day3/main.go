// started        at: 2025-12-03 11:27:17+02:00
// finished part1 at: 2025-12-03 11:57:22+02:00
// finished part2 at: 2025-12-03 12:32:58+02:00
// part1: 30m 5s, part2: 35m 36s

package main

import (
	_ "embed"
	"flag"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

//go:embed test.txt
var testInput string

// init() runs before main() and formats/validates the embedded inputs.
func init() {
	input = strings.TrimRight(input, "\n")
	if len(input) == 0 {
		panic("empty input.txt file")
	}
	testInput = strings.TrimRight(testInput, "\n")
	if len(testInput) == 0 {
		panic("empty test.txt file")
	}
}

func main() {
	var part int
	var test bool
	flag.IntVar(&part, "part", 1, "part 1 or 2")
	flag.BoolVar(&test, "test", false, "run with test.txt inputs?")
	flag.Parse()
	fmt.Println("Running part", part, ", test inputs = ", test)

	if test {
		input = testInput
	}

	var answer int
	switch part {
	case 1:
		answer = part1(input)
	case 2:
		answer = part2(input)
	default:
		panic("invalid challenge part, must be 1 or 2")
	}
	fmt.Println("Output:", answer)
}

// part1 solves part 1 of the day's challenge.
// Keep this function signature intact for unit tests to work seamlessly.
func part1(input string) int {
	parsed := parseInput(input)
	fmt.Println(parsed)

	finalSum := 0

	for _, row := range parsed {
		digit, idx := findHighestDigit(row)
		substringToTheRight := row[idx+1:]
		fmt.Printf("Highest num is %d at index %d, the substring to the right is %s\n", digit, idx, substringToTheRight)

		if len(substringToTheRight) == 0 { // the highest num is last
			digit, idx = findSecondHighestDigit(row, digit)
			substringToTheRight = row[idx+1:]
		}

		secondDigit, _ := findHighestDigit(substringToTheRight)

		finalSum += digit*10 + secondDigit //nolint:revive // add-constant - this is fine
	}

	return finalSum
}

func findHighestDigit(input string) (digit, index int) {
	highestDigit := slices.Max(rowToInts(input))
	idx := strings.Index(input, strconv.Itoa(highestDigit))

	return highestDigit, idx
}

func findSecondHighestDigit(input string, highestDigit int) (digit, index int) {
	inputWithoutHighest := strings.ReplaceAll(input, strconv.Itoa(highestDigit), "")
	secondHighestDigit := slices.Max(rowToInts(inputWithoutHighest))
	idx := strings.Index(input, strconv.Itoa(secondHighestDigit))

	return secondHighestDigit, idx
}

// part2 solves part 2 of the day's challenge.
// Keep this function signature intact for unit tests to work seamlessly.
func part2(input string) int {
	parsed := parseInput(input)
	fmt.Println(parsed)

	finalSum := 0

	for _, row := range parsed {
		batteries := find12Batteries(row)

		batterySum := 0
		for _, digit := range batteries {
			batterySum = batterySum*10 + digit
		}

		finalSum += batterySum
	}

	return finalSum
}

func find12Batteries(row string) []int {
	var digits []int
	for neededDigits := 11; neededDigits >= 0; neededDigits-- {
		digit, idx := findHighestDigitWithSufficientSpaceForNMoreDigits(row, neededDigits)
		digits = append(digits, digit)
		rowToTheRight := row[idx+1:]

		// fmt.Printf("In row %s, found highest digit %d at index %d, reducing the row to %s\n", row, digit, idx, rowToTheRight)
		row = rowToTheRight
	}
	return digits
}

func findHighestDigitWithSufficientSpaceForNMoreDigits(input string, digitsToTheRight int) (digit, index int) {
	for digitMustBeLowerThan := 10; digitMustBeLowerThan > 1; digitMustBeLowerThan-- {
		// fmt.Printf("Looking for the highest digit in %s that allows for %d digits to the right while being lower than digit %d... ", input, digitsToTheRight, digitMustBeLowerThan)
		d, idx := findHighestDigitLowerThan(input, digitMustBeLowerThan)
		if len(input[idx+1:]) >= digitsToTheRight {
			// fmt.Printf("Found digit %d at index %d\n", d, idx)
			return d, idx
		}
		// fmt.Println("Didn't find it, gonna look for a lower digit")
	}
	panic("This shouldn't happen")
}

func findHighestDigitLowerThan(input string, lowerThanThis int) (digit, index int) {
	var builder strings.Builder
	for _, digit := range input {
		if int(digit-'0') < lowerThanThis {
			builder.WriteByte(byte(digit))
		}
	}
	newRow := builder.String()

	highestDigit := slices.Max(rowToInts(newRow))
	idx := strings.Index(input, strconv.Itoa(highestDigit))

	return highestDigit, idx
}

func parseInput(input string) []string {
	return strings.Split(input, "\n")
}

func rowToInts(input string) []int {
	output := make([]int, len(input))
	for i, digit := range input {
		output[i] = int(digit - '0') // converts to int without unicode/ascii shenanigans
	}
	return output
}
