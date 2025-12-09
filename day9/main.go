// started        at: 2025-12-09 18:21:47+02:00
// pause          at: 2025-12-09 18:54:35+02:00
// resume         at: 2025-12-09 22:50:14+02:00
// finished part1 at: 2025-12-09 23:39:57+02:00 // in reality this was like 5min earlier - after getting the answer I tested out goroutine speedup first
// pause          at: 2025-12-09 23:47:19+02:00 // motherfucker, my part 1 algo is useless for part 2
// finished part2 at: ---
// part1: 1h 16m, part12: ---

package main

import (
	_ "embed"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"sync"
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

// ========== UNIQUE DAY SOLUTION CODE BELOW ==========

type tile struct {
	x int // to the right, starting at 0
	y int // downwards, starting at 0
}

// 1. look for "cut corner" tiles
// 2. when found all four corners, find out the area of each combination
// 3. return the max

// part1 solves part 1 of the day's challenge.
// Keep this function signature intact for unit tests to work seamlessly.
func part1(input string) int {
	tiles, maxX, maxY := parseInput(input)
	fmt.Println(tiles)

	tileMap := make(map[tile]struct{}) // this map is purely used to check for existence
	for _, t := range tiles {
		tileMap[t] = struct{}{}
	}

	maxGridSize := max(maxX, maxY)
	var candidatesUpLeft, candidatesUpRight, candidatesBottomLeft, candidatesBottomRight []tile
	// var candidatesUpLeft []tile

	var wg sync.WaitGroup

	for distanceFromCorner := range maxGridSize { // current search distance from corner
		// ~TODO~: test if running 4 goroutines, one for each corner, is faster
		// result: it does! 26s to 8s!

		if len(candidatesUpLeft) == 0 { // we're still looking for this
			wg.Go(func() {
				for i := 0; i <= distanceFromCorner; i++ { // i is the preference to check xDistance first (0 only checks the first line, 1 checks the second, etc.)
					tileToLookFor := tile{x: distanceFromCorner - i, y: i}
					if _, exists := tileMap[tileToLookFor]; exists { // found a tile
						candidatesUpLeft = append(candidatesUpLeft, tileToLookFor)
					}
				}
			})
		}

		if len(candidatesUpRight) == 0 { // we're still looking for this
			wg.Go(func() {
				for i := 0; i <= distanceFromCorner; i++ { // i is the preference to check xDistance first (0 only checks the first line, 1 checks the second, etc.)
					tileToLookFor := tile{x: maxGridSize - distanceFromCorner + i, y: i}
					if _, exists := tileMap[tileToLookFor]; exists { // found a tile
						candidatesUpRight = append(candidatesUpRight, tileToLookFor)
					}
				}
			})
		}

		if len(candidatesBottomLeft) == 0 { // we're still looking for this
			wg.Go(func() {
				for i := 0; i <= distanceFromCorner; i++ { // i is the preference to check xDistance first (0 only checks the first line, 1 checks the second, etc.)
					tileToLookFor := tile{x: i, y: maxGridSize - distanceFromCorner + i}
					if _, exists := tileMap[tileToLookFor]; exists { // found a tile
						candidatesBottomLeft = append(candidatesBottomLeft, tileToLookFor)
					}
				}
			})
		}

		if len(candidatesBottomRight) == 0 { // we're still looking for this
			wg.Go(func() {
				for i := 0; i <= distanceFromCorner; i++ { // i is the preference to check xDistance first (0 only checks the first line, 1 checks the second, etc.)
					tileToLookFor := tile{x: maxGridSize - i, y: maxGridSize - distanceFromCorner + i}
					if _, exists := tileMap[tileToLookFor]; exists { // found a tile
						candidatesBottomRight = append(candidatesBottomRight, tileToLookFor)
					}
				}
			})
		}

		wg.Wait()

		// if len(candidatesUpLeft) > 0 {
		if len(candidatesUpLeft) > 0 && len(candidatesUpRight) > 0 && len(candidatesBottomLeft) > 0 && len(candidatesBottomRight) > 0 {
			// all corners are found
			break
		}
	}

	fmt.Println()
	fmt.Printf("candidatesUpLeft: %v\n", candidatesUpLeft)
	fmt.Printf("candidatesUpRight: %v\n", candidatesUpRight)
	fmt.Printf("candidatesBottomLeft: %v\n", candidatesBottomLeft)
	fmt.Printf("candidatesBottomRight: %v\n", candidatesBottomRight)

	fmt.Println()
	// pair top-left and bottom-right
	maxAreaA := 0
	for _, tile1 := range candidatesUpLeft {
		for _, tile2 := range candidatesBottomRight {
			xDist := max(tile1.x, tile2.x) - min(tile1.x, tile2.x) + 1
			yDist := max(tile1.y, tile2.y) - min(tile1.y, tile2.y) + 1
			area := xDist * yDist
			if area > maxAreaA {
				fmt.Printf("Between tiles %v and %v, xDist is %d and yDist is %d\n", tile1, tile2, xDist, yDist)
				maxAreaA = area
			}
		}
	}
	fmt.Printf("maxAreaA: %d\n", maxAreaA)

	fmt.Println()
	// pair top-right and bottom-left
	maxAreaB := 0
	for _, tile1 := range candidatesUpRight {
		for _, tile2 := range candidatesBottomLeft {
			xDist := max(tile1.x, tile2.x) - min(tile1.x, tile2.x) + 1
			yDist := max(tile1.y, tile2.y) - min(tile1.y, tile2.y) + 1
			area := xDist * yDist
			if area > maxAreaB {
				fmt.Printf("Between tiles %v and %v, xDist is %d and yDist is %d\n", tile1, tile2, xDist, yDist)
				maxAreaB = area
			}
		}
	}
	fmt.Printf("maxAreaB: %d\n", maxAreaB)

	return max(maxAreaA, maxAreaB)
}

// part2 solves part 2 of the day's challenge.
// Keep this function signature intact for unit tests to work seamlessly.
func part2(input string) int {
	// parsed := parseInput(input)
	// fmt.Println(parsed)

	return 0
}

func parseInput(input string) (tiles []tile, maxX, maxY int) {
	var parsedInput []tile
	for line := range strings.SplitSeq(input, "\n") {
		xString, yString, found := strings.Cut(line, ",")
		if !found {
			panic("parseInput: not coordinates, wtf")
		}

		x := stringToInt(xString)
		y := stringToInt(yString)

		if x > maxX {
			maxX = x
		}
		if y > maxY {
			maxY = y
		}
		parsedInput = append(parsedInput, tile{x: x, y: y})
	}
	return parsedInput, maxX, maxY
}

func stringToInt(input string) int {
	output, err := strconv.Atoi(input)
	if err != nil {
		panic(fmt.Sprintf("stringToInt: failed to convert %q to int: %v", input, err))
	}
	return output
}
