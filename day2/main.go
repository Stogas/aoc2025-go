// started        at: 2025-12-02 10:48:28+02:00 // went on a bathroom break for about 6 minutes tho
// finished part1 at: 2025-12-02 11:24:34+02:00
// finished part2 at: ---

package main

import (
	_ "embed"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

//go:embed test.txt
var testInput string

type idRange struct {
	firstID int
	lastID  int
}

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

	var sumOfInvalidNumbers int
	for _, numberRange := range parsed {
		for _, invalidNumber := range processRange(numberRange) {
			fmt.Printf("invalid ID: %d\n", invalidNumber)
			sumOfInvalidNumbers += invalidNumber
		}
	}

	return sumOfInvalidNumbers
}

// part2 solves part 2 of the day's challenge.
// Keep this function signature intact for unit tests to work seamlessly.
func part2(input string) int {
	parsed := parseInput(input)
	fmt.Println(parsed)

	return 0
}

func parseInput(input string) []idRange {
	var parsedInput []idRange
	for line := range strings.SplitSeq(input, ",") {
		parsedInput = append(parsedInput, stringToIDRange(line))
	}
	return parsedInput
}

func stringToIDRange(input string) idRange {
	firstString, lastString, found := strings.Cut(input, "-")
	if !found {
		panic("stringToIdRange: failed to find - in ID range")
	}

	firstID, err := strconv.Atoi(firstString)
	if err != nil {
		panic(fmt.Sprintf("stringToInt: failed to convert %q to int: %v", input, err))
	}
	lastID, err := strconv.Atoi(lastString)
	if err != nil {
		panic(fmt.Sprintf("stringToInt: failed to convert %q to int: %v", input, err))
	}

	return idRange{
		firstID: firstID,
		lastID:  lastID,
	}
}

func isValid(input int) bool {
	inputString := strconv.Itoa(input)
	length := len(inputString)

	if length%2 != 0 { // odd length, cannot be symmetrical
		return true
	}

	// fmt.Printf("first half: %v; second half: %v\n", inputString[:length/2-1], inputString[length/2:])

	if inputString[:length/2] == inputString[length/2:] {
		return false
	}

	return true
}

func processRange(input idRange) []int {
	var result []int
	for number := input.firstID; number <= input.lastID; number++ {
		// fmt.Printf("processing number %d\n", number)
		if !isValid(number) {
			result = append(result, number)
		}
	}
	return result
}
