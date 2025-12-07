// started        at: 2025-12-07 14:15:00+02:00 # note: this is an estimate, forgot to capture it accurately
// finished part1 at: 2025-12-07 14:33:03+02:00
// finished part2 at: 2025-12-07 14:41:44+02:00 # I was lucky to have a good algorithm in part1

package main

import (
	_ "embed"
	"flag"
	"fmt"
	"strings"
)

//go:embed input.txt
var input string

//go:embed test.txt
var testInput string

type coordinates struct {
	rowIdx    int
	columnIdx int
}

// part1 idea:
// 1. var activeSplitters []splitter
// 2. var beam [len(row)]bool
// 3. process each line, add activeSplitters, mark beams

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
	rows, source := parseInput(input)
	// fmt.Printf("Rows: %v", rows)
	fmt.Printf("Source coords: %d,%d\n", source.rowIdx, source.columnIdx)

	currentBeams := make([]bool, len(rows[0]))
	totalSplitCount := 0

	for _, row := range rows {
		totalSplitCount += createOrSplitBeams(row, currentBeams)
	}

	return totalSplitCount
}

func createOrSplitBeams(row string, currentBeams []bool) (splitCount int) {
	// note that in Go, slices (like []bool) are passed essentially as pointers, so we can modify it here without returning it

	for idx, char := range row {
		if char == 'S' { // there's a beam source here
			currentBeams[idx] = true
		}
		if currentBeams[idx] { // there's a beam descending here
			if char == '^' { // and a splitter that splits it
				if idx-1 >= 0 { // there's space to the left to split a beam there
					currentBeams[idx-1] = true
				}
				if idx+1 < len(row) { // there's space to the right to split a beam there
					currentBeams[idx+1] = true
				}
				currentBeams[idx] = false // we split the beam, it's not in this column anymore
				splitCount++
			}
		}
	}

	return splitCount
}

// part2 solves part 2 of the day's challenge.
// Keep this function signature intact for unit tests to work seamlessly.
func part2(input string) int {
	rows, source := parseInput(input)
	// fmt.Printf("Rows: %v", rows)
	fmt.Printf("Source coords: %d,%d\n", source.rowIdx, source.columnIdx)

	currentTimelines := make([]int, len(rows[0]))

	for _, row := range rows {
		createOrSplitTimelines(row, currentTimelines)
	}

	totalTimelines := 0
	for _, timelineCount := range currentTimelines {
		totalTimelines += timelineCount
	}

	return totalTimelines
}

func createOrSplitTimelines(row string, currentTimelines []int) {
	// note that in Go, slices (like []bool) are passed essentially as pointers, so we can modify it here without returning it

	for idx, char := range row {
		if char == 'S' { // there's a beam source here
			currentTimelines[idx] = 1
		}
		if currentTimelines[idx] > 0 { // there's a timeline descending here
			if char == '^' { // and a splitter that splits it
				if idx-1 >= 0 { // there's space to the left to split timelines there
					currentTimelines[idx-1] += currentTimelines[idx]
				}
				if idx+1 < len(row) { // there's space to the right to split timelines there
					currentTimelines[idx+1] += currentTimelines[idx]
				}
				currentTimelines[idx] = 0 // we split the timelines, it's not in this column anymore
			}
		}
	}
}

func parseInput(input string) (rows []string, source coordinates) {
	for rIdx, line := range strings.Split(input, "\n") {
		rows = append(rows, line)
		if strings.Contains(line, "S") {
			source.rowIdx = rIdx
			source.columnIdx = strings.Index(line, "S")
		}
	}
	return rows, source
}

// func stringToInt(input string) int {
// 	output, err := strconv.Atoi(input)
// 	if err != nil {
// 		panic(fmt.Sprintf("stringToInt: failed to convert %q to int: %v", input, err))
// 	}
// 	return output
// }
