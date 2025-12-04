// started        at: 2025-12-04 11:27:26+02:00
// finished part1 at: 2025-12-04 11:53:56+02:00
// finished part2 at: 2025-12-04 12:06:07+02:00

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

type cellGrid [][]bool

func (g cellGrid) adjacentRolls(rowIdx, columnIdx int) int {
	maxRowIdx := len(g) - 1
	maxColumnIdx := len(g[0]) - 1

	count := 0

	for cRowIdx := rowIdx - 1; cRowIdx <= rowIdx+1; cRowIdx++ {
		if cRowIdx < 0 || cRowIdx > maxRowIdx {
			continue
		}

		for cColumnIdx := columnIdx - 1; cColumnIdx <= columnIdx+1; cColumnIdx++ {
			if cColumnIdx < 0 || cColumnIdx > maxColumnIdx || (cRowIdx == rowIdx && cColumnIdx == columnIdx) {
				continue
			}

			if g[cRowIdx][cColumnIdx] {
				// fmt.Printf("Found adjacent at %d,%d\n", cRowIdx, cColumnIdx)
				count++
			}
		}
	}

	return count
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

	countOfAccessibleRolls := 0

	for rowIdx, row := range parsed {
		for columnIdx, cell := range row {
			if !cell {
				continue
			}

			// fmt.Printf("Checking roll at %d,%d\n", rowIdx, columnIdx)

			countOfAdjacentRolls := parsed.adjacentRolls(rowIdx, columnIdx)
			if countOfAdjacentRolls < 4 {
				fmt.Printf("Found accessible roll at %d,%d, adjacent count was %d\n", rowIdx, columnIdx, countOfAdjacentRolls)
				countOfAccessibleRolls++
			}
		}
	}

	return countOfAccessibleRolls
}

// part2 solves part 2 of the day's challenge.
// Keep this function signature intact for unit tests to work seamlessly.
func part2(input string) int {
	parsed := parseInput(input)
	fmt.Println(parsed)

	// essentially a stupid bruteforce, but sub-10 milliseconds is fast enough for me to consider this performant (see bench.md)

	var removedRolls int
	for {
		removed := parsed.removeAccessibleRolls()
		if removed == 0 {
			break
		}
		removedRolls += removed
	}

	return removedRolls
}

type coordinates struct {
	x int
	y int
}

func (g cellGrid) removeAccessibleRolls() (removed int) {
	var accessibleRolls []coordinates

	for rowIdx, row := range g {
		for columnIdx, cell := range row {
			if !cell {
				continue
			}

			// fmt.Printf("Checking roll at %d,%d\n", rowIdx, columnIdx)

			countOfAdjacentRolls := g.adjacentRolls(rowIdx, columnIdx)
			if countOfAdjacentRolls < 4 {
				// fmt.Printf("Found accessible roll at %d,%d, adjacent count was %d\n", rowIdx, columnIdx, countOfAdjacentRolls)
				accessibleRolls = append(accessibleRolls, coordinates{x: rowIdx, y: columnIdx})
			}
		}
	}

	for _, cell := range accessibleRolls {
		g[cell.x][cell.y] = false
	}

	return len(accessibleRolls)
}

func parseInput(input string) cellGrid {
	var parsedInput cellGrid
	for line := range strings.SplitSeq(input, "\n") {
		parsedInput = append(parsedInput, stringToBools(line))
	}
	return parsedInput
}

func stringToBools(input string) []bool {
	var output []bool
	for _, r := range input {
		output = append(output, r == '@')
	}
	return output
}
