// started        at: 2025-12-05 13:01:07+02:00
// paused         at: 2025-12-05 14:24:15+02:00
// resumed        at: 2025-12-05 16:19:56+02:00
// finished part1 at: 2025-12-05 17:16:14+02:00 # fuckin finally
// finished part2 at: 2025-12-05 17:20:46+02:00 # at least this was instant

package main

import (
	_ "embed"
	"flag"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

//go:embed test.txt
var testInput string

//go:embed test2.txt
var testInput2 string

type freshRange struct {
	from int
	upTo int
}

type database struct {
	ranges      []freshRange
	ingredients []int
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
	var test2 bool
	flag.IntVar(&part, "part", 1, "part 1 or 2")
	flag.BoolVar(&test, "test", false, "run with test.txt inputs?")
	flag.BoolVar(&test2, "test2", false, "run with test2.txt inputs?")
	flag.Parse()
	fmt.Println("Running part", part, ", test inputs = ", test)

	if test {
		input = testInput
	}
	if test2 {
		input = testInput2
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
//
// plan:
// 1. sort ranges
// 2. deduplicate overlapping ranges
// 3. sort ingredients
// 4. check sequentially.
func part1(input string) int {
	db := parseInput(input)
	fmt.Println(db)

	// 1. sort ranges based on the "from" number
	sort.Slice(db.ranges, func(i, j int) bool {
		return db.ranges[i].from < db.ranges[j].from
	})
	fmt.Printf("Sorted ranges DB:\n%v\n\n", db)

	// 2. deduplicate overlapping ranges
	db.ranges = deduplicateRanges2(db.ranges)

	// 3. sort ingredients
	sort.Ints(db.ingredients)

	// print after sorting everything and deduplicating:
	fmt.Printf("Sorted and deduplicated DB:\n%v\n\n", db)

	// 4. check sequentially
	freshItemCount, _ := countFreshItems2(db)
	return freshItemCount
}

// deduplicateRanges - UNSUED, deprecated by deduplicateRanges2()
// damn, this was tough to debug edge cases for. LeetCode helped: https://leetcode.com/problems/merge-intervals/
// func deduplicateRanges(ranges []freshRange) []freshRange {
// 	var deduplicatedRanges []freshRange

// 	currentlyOverlapping := false
// 	overlappingFrom := 0
// 	overlappingTo := 0

// 	for i := 0; i < len(ranges)-1; i++ {
// 		fmt.Printf("Working on range %d-%d...\n", ranges[i].from, ranges[i].upTo)

// 		if ranges[i].upTo >= ranges[i+1].from { // overlaps!
// 			fmt.Printf("\tOverlaps with the next range %d-%d!\n", ranges[i+1].from, ranges[i+1].upTo)

// 			if !currentlyOverlapping { // new overlap range, need to set start marker
// 				overlappingFrom = ranges[i].from
// 				overlappingTo = ranges[i].upTo
// 			} else { // not new range, need to find min/max
// 				overlappingFrom = min(ranges[i].from, overlappingFrom)
// 				overlappingTo = max(ranges[i].upTo, overlappingTo)
// 			}
// 			currentlyOverlapping = true
// 			continue
// 		}

// 		// we don't overlap with the next item
// 		if currentlyOverlapping { // but we finished an overlap range, so add it
// 			fmt.Printf("\tAdding overlapped range as %d-%d\n", overlappingFrom, max(ranges[i].upTo, overlappingTo))
// 			deduplicatedRanges = append(deduplicatedRanges, freshRange{from: min(ranges[i].from, overlappingFrom), upTo: max(ranges[i].upTo, overlappingTo)})
// 			currentlyOverlapping = false
// 			continue
// 		}

// 		fmt.Printf("\tAdding non-overlapped range as %d-%d\n", ranges[i].from, ranges[i].upTo)
// 		deduplicatedRanges = append(deduplicatedRanges, ranges[i])
// 	}

// 	// we didn't process the last element in the ranges slice, so we need to do so if it's still overlapping
// 	if currentlyOverlapping {
// 		deduplicatedRanges = append(deduplicatedRanges, freshRange{from: min(ranges[len(ranges)-1].from, overlappingFrom), upTo: max(ranges[len(ranges)-1].upTo, overlappingTo)})
// 	} else {
// 		deduplicatedRanges = append(deduplicatedRanges, ranges[len(ranges)-1])
// 	}

// 	return deduplicatedRanges
// }

func deduplicateRanges2(ranges []freshRange) []freshRange {
	var deduplicatedRanges []freshRange

	overlappingFrom := ranges[0].from
	overlappingTo := ranges[0].upTo

	for i := 1; i < len(ranges); i++ {
		fmt.Printf("Working on range %d-%d...\n", ranges[i].from, ranges[i].upTo)

		if ranges[i].from <= overlappingTo+1 { // overlaps! (we can include consecutive ranges too, i.e. without gaps - technically can be merged to a single range)
			fmt.Printf("\tOverlaps with the previous range %d-%d!\n", overlappingFrom, overlappingTo)

			overlappingTo = max(ranges[i].upTo, overlappingTo)
			continue
		}

		// we don't overlap with the previous ranges
		fmt.Printf("\tAdding previous range as %d-%d\n", overlappingFrom, overlappingTo)
		deduplicatedRanges = append(deduplicatedRanges, freshRange{from: overlappingFrom, upTo: overlappingTo})

		overlappingFrom = ranges[i].from
		overlappingTo = ranges[i].upTo
	}

	// add the last unadded range
	deduplicatedRanges = append(deduplicatedRanges, freshRange{from: overlappingFrom, upTo: overlappingTo})

	return deduplicatedRanges
}

// countFreshItems2 ASSUMES BOTH RANGES AND INGREDIENTS ARE SORTED.
func countFreshItems2(db database) (freshItemCount int, freshItems []int) {
	fItemCount := 0
	var fItems []int

	// rIdx := scanForFirstRange(db.ranges, db.ingredients[0]) // rIdx is range index that could contai item. Lower index ranges are guaranteed to not have this number (i.e. ranges are lower)

	rangeIdx := 0
	itemIdx := 0

	for rangeIdx < len(db.ranges) && itemIdx < len(db.ingredients) {
		item := db.ingredients[itemIdx]
		currentRange := db.ranges[rangeIdx]

		if item < currentRange.from { // item is lower than this range, item can be discarded
			itemIdx++
			continue
		}

		if item > currentRange.upTo { // item is bigger than this range, range can be discarded
			rangeIdx++
			continue
		}

		fItemCount++
		fItems = append(fItems, item)

		itemIdx++
	}

	return fItemCount, fItems
}

// part2 solves part 2 of the day's challenge.
// Keep this function signature intact for unit tests to work seamlessly.
func part2(input string) int {
	db := parseInput(input)
	fmt.Println(db)

	// 1. sort ranges based on the "from" number
	sort.Slice(db.ranges, func(i, j int) bool {
		return db.ranges[i].from < db.ranges[j].from
	})
	fmt.Printf("Sorted ranges DB:\n%v\n\n", db)

	// 2. deduplicate overlapping ranges
	db.ranges = deduplicateRanges2(db.ranges)

	// print after sorting everything and deduplicating:
	fmt.Printf("Sorted and deduplicated DB:\n%v\n\n", db)

	// 4. check sequentially
	possibleFreshItems := 0
	for _, r := range db.ranges {
		possibleFreshItems += r.upTo - r.from + 1
	}
	return possibleFreshItems
}

func parseInput(input string) database {
	var db database

	processingRanges := true
	for line := range strings.SplitSeq(input, "\n") { // well this is overly complicated, but oh well
		if len(line) == 0 {
			processingRanges = false
			continue
		}

		if processingRanges {
			before, after, found := strings.Cut(line, "-")
			if !found {
				panic("parseInput: range hyphen not found")
			}

			beforeInt, err := strconv.Atoi(before)
			if err != nil {
				panic("parseInput: range not a number")
			}
			afterInt, err := strconv.Atoi(after)
			if err != nil {
				panic("parseInput: range not a number")
			}

			db.ranges = append(db.ranges, freshRange{from: beforeInt, upTo: afterInt})
			continue
		}

		ingredient, err := strconv.Atoi(line)
		if err != nil {
			panic("parseInput: ingredient not a number")
		}

		db.ingredients = append(db.ingredients, ingredient)
	}
	return db
}

// func stringToInt(input string) int {
// 	output, err := strconv.Atoi(input)
// 	if err != nil {
// 		panic(fmt.Sprintf("stringToInt: failed to convert %q to int: %v", input, err))
// 	}
// 	return output
// }
