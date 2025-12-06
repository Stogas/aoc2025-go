// started        at: 2025-12-06 09:46:12+02:00
// finished part1 at: 2025-12-06 10:09:06+02:00 # well this was easy
// paused part2   at: 2025-12-06 10:36:00+02:00
// resumed part2  at: 2025-12-06 11:34:21+02:00
// finished part2 at: 2025-12-06 11:39:54+02:00
// part1: 22m 54s, part2: 32m 27s (active)

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

type worksheet struct {
	rows               [][]int
	multiplyOperations []bool
}

// type worksheet2 struct {
// 	columns            [][]string
// 	multiplyOperations []bool
// }

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
	ws := parseInput(input)
	fmt.Println(ws)

	grandTotal := 0

	for cIdx := range len(ws.multiplyOperations) {
		var columnResult int
		if ws.multiplyOperations[cIdx] {
			columnResult = 1
		} else {
			columnResult = 0
		}

		for rIdx := range len(ws.rows) {
			if ws.multiplyOperations[cIdx] { // multiply
				columnResult *= ws.rows[rIdx][cIdx]
			} else { // sum
				columnResult += ws.rows[rIdx][cIdx]
			}
		}
		grandTotal += columnResult
	}

	return grandTotal
}

// part2 solves part 2 of the day's challenge.
// Keep this function signature intact for unit tests to work seamlessly.
func part2(input string) int {
	// ws := parseInput2(input)
	// fmt.Println(ws)

	return parseInputAndSolvePart2(input)
}

func parseInput(input string) worksheet {
	rows := strings.Split(input, "\n")

	var parsedInput worksheet
	for rIdx := range rows {
		if rIdx < len(rows)-1 { // numbers rows
			parsedInput.rows = append(parsedInput.rows, rowToInts(rows[rIdx]))
			continue
		}
		// operations row
		parsedInput.multiplyOperations = operationsRowToBools(rows[rIdx])
	}
	return parsedInput
}

func rowToInts(row string) []int {
	var numbersInRow []int
	for numString := range strings.FieldsSeq(row) {
		num, err := strconv.Atoi(numString)
		if err != nil {
			panic(fmt.Sprintf("rowToInts: failed to convert %q to int: %v", numString, err))
		}

		numbersInRow = append(numbersInRow, num)
	}
	return numbersInRow
}

func operationsRowToBools(row string) []bool {
	var multiplyOperations []bool
	for operationString := range strings.FieldsSeq(row) {
		multiplyOperations = append(multiplyOperations, operationString == "*")
	}
	return multiplyOperations
}

func parseInputAndSolvePart2(input string) int {
	rows := strings.Split(input, "\n")
	rowLength := len(rows[0])

	var delimiterPositions []int

	for rowPosition := range rowLength {
		potentialDelimiter := true
		for _, row := range rows {
			if row[rowPosition] != ' ' { // not a delimiter colunm
				potentialDelimiter = false
				continue
			}
		}
		if potentialDelimiter {
			delimiterPositions = append(delimiterPositions, rowPosition)
		}
	}
	// add a fake delimiter after the line ends
	delimiterPositions = append(delimiterPositions, rowLength)

	globalSum := 0

	for dIdx, dPos := range delimiterPositions {
		var previousDelimiterPos int
		if dIdx == 0 {
			previousDelimiterPos = -1
		} else {
			previousDelimiterPos = delimiterPositions[dIdx-1]
		}

		fmt.Printf("Checking for operation in dIDx: %d, dPos: %d\n", dIdx, dPos)

		multiply := rows[len(rows)-1][previousDelimiterPos+1] == '*'
		localTotal := 0
		if multiply {
			localTotal = 1
		}

		for column := dPos - 1; column > previousDelimiterPos; column-- {
			var sb strings.Builder
			for rowIdx := 0; rowIdx < len(rows)-1; rowIdx++ { //nolint:intrange // can't change to integer range, not valid suggestion by linter
				sb.WriteByte(rows[rowIdx][column])
			}
			currentNumberString := sb.String()

			currentNumber, err := strconv.Atoi(strings.TrimSpace(currentNumberString))
			if err != nil {
				panic("parseInputAndSolvePart2: couldn't convert number")
			}

			if multiply {
				localTotal *= currentNumber
			} else {
				localTotal += currentNumber
			}
		}

		globalSum += localTotal
	}

	return globalSum
}
