// started        at: 2025-12-01 17:29:40+02:00
// finished part1 at: 2025-12-01 18:11:31+02:00
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

const (
	minDialValue           int = 0
	maxDialValue           int = 99
	dialOverflowCorrection int = 100
	startingDialValue      int = 50
)

type command struct {
	clockwise bool
	steps     int
}

type dial struct {
	currentState int
}

func newDial() dial {
	var d dial
	d.currentState = startingDialValue
	return d
}

func (d *dial) rotateDial(c command) {
	var directionMultiplier = -1 // if counter-clockwise
	if c.clockwise {
		directionMultiplier = 1 // if clockwise
	}

	targetUnsafeState := d.currentState + directionMultiplier*c.steps

	d.currentState = fixDialOverflow(targetUnsafeState)
}

func (d dial) isPointingAtZero() bool {
	return d.currentState == 0
}

func fixDialOverflow(input int) int { // designed with recursion in mind
	if input < minDialValue {
		return fixDialOverflow(input + dialOverflowCorrection)
	} else if input > maxDialValue {
		return fixDialOverflow(input - dialOverflowCorrection)
	}

	return input
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
	case 2: //nolint:revive // add-constant rule triggers here, let's just hardcode it
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

	vault := newDial()
	numberOfZeroStates := 0

	for _, cmd := range parsed {
		vault.rotateDial(cmd)
		fmt.Printf("current state: %d\n", vault.currentState)
		if vault.isPointingAtZero() {
			numberOfZeroStates++
		}
	}

	return numberOfZeroStates
}

// part2 solves part 2 of the day's challenge.
// Keep this function signature intact for unit tests to work seamlessly.
func part2(input string) int {
	parsed := parseInput(input)
	fmt.Println(parsed)

	return 0
}

func parseInput(input string) []command {
	var parsedInput []command
	for line := range strings.SplitSeq(input, "\n") {
		parsedInput = append(parsedInput, commandParse(line))
	}
	return parsedInput
}

func commandParse(input string) command {
	// this might seem yucky using runes and converting back and forth,
	// but I'm paranoid that we'll encounter something else other than L or R,
	// so I can't assume the first character will always consume a byte, not more.

	return command{
		clockwise: clockwiseParse([]rune(input)[0]),
		steps:     stringToInt(string([]rune(input)[1:])),
	}
}

func clockwiseParse(input rune) bool {
	switch input {
	case 'L':
		return false
	case 'R':
		return true
	default:
		panic("clockwiseParse had an error, not L nor R")
	}
}

func stringToInt(input string) int {
	output, err := strconv.Atoi(input)
	if err != nil {
		panic(fmt.Sprintf("stringToInt: failed to convert %s to int: %v", input, err))
	}
	return output
}
