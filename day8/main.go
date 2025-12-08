// started        at: 2025-12-08 18:37:21+02:00
// paused         at: 2025-12-08 18:58:55+02:00
// resumed        at: 2025-12-08 19:24:30+02:00
// finished part1 at: 2025-12-08 21:45:29+02:00
// finished part2 at: 2025-12-08 22:24:55+02:00
// part1: 2h 42m 33s, part2: 39m 26s
//
//
// although neither of these part solutions are satisfactory,
// both have crazy runtimes (see bench.md).
// I'll have to revisit this and find a better algorithm

package main

import (
	_ "embed"
	"flag"
	"fmt"
	"math"
	"sort"
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
		answer = part1(input, test)
	case 2:
		answer = part2(input, test)
	default:
		panic("invalid challenge part, must be 1 or 2")
	}
	fmt.Println("Output:", answer)
}

// ========== UNIQUE DAY SOLUTION CODE BELOW ==========

type coordinates struct {
	x int
	y int
	z int
}

type junction struct {
	position      coordinates
	circuitNumber int
	connected     bool // is connected to something, i.e. belongs to a circuit with some other junction
}

type connection struct {
	junctionIdx [2]int
	distance    float64
}

type circuit struct {
	number int
	size   int
}

// part1 solves part 1 of the day's challenge.
// Keep this function signature intact for unit tests to work seamlessly.
func part1(input string, test bool) int {
	junctions := parseInput(input)
	fmt.Println(junctions)

	loopCount := 1000
	if test {
		loopCount = 10
	}

	fmt.Println("")
	fmt.Println("Finding all distances sorted increasing")
	allDistances := findAllDistancesSortedIncreasing(junctions, loopCount)

	fmt.Println("Assigning junctions to circuits")
	assignJunctionsToCircuits(junctions, allDistances, loopCount)

	circuitSizesMap := circuitSizes(junctions)

	fmt.Println()
	// change circuits from map to slice of objects
	var circuits []circuit
	for k, v := range circuitSizesMap {
		circuits = append(circuits, circuit{number: k, size: v})
	}
	fmt.Printf("Count of all circuits: %d\n", len(circuits))

	// sort circuits by descending value
	sort.Slice(circuits, func(i, j int) bool {
		return circuits[i].size > circuits[j].size
	})

	fmt.Println()
	fmt.Printf("All circuits:\n%v\n", circuits)

	missingCiruitsTo3 := 3 - len(circuits)
	for range missingCiruitsTo3 {
		circuits = append(circuits, circuit{number: 0, size: 1})
	}

	fmt.Printf("Biggest circuit sizes:\n%v\n%v\n%v\n", circuits[0].size, circuits[1].size, circuits[2].size)

	return circuits[0].size * circuits[1].size * circuits[2].size
}

func findAllDistancesSortedIncreasing(j []junction, loopCount int) (distances []connection) {
	for i := range j { // loop over every junction
		if i >= len(j)-1 { // last element
			break
		}
		for ii := i + 1; ii < len(j); ii++ { // loop over every junction after it (previous ones already have connections calculated)
			c := connection{junctionIdx: [2]int{i, ii}, distance: euclideanDistance(j[i], j[ii])}
			distances = insertSortedDistance(distances, c, loopCount)
		}
	}

	return distances
}

func euclideanDistance(a, b junction) float64 {
	dx := float64(b.position.x - a.position.x)
	dy := float64(b.position.y - a.position.y)
	dz := float64(b.position.z - a.position.z)
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func insertSortedDistance(distances []connection, c connection, loopCount int) []connection {
	if len(distances) == 0 {
		return []connection{c}
	}

	// binary search to find insertion point

	left, right := 0, len(distances)
	for left < right {
		mid := (left + right) / 2
		if distances[mid].distance < c.distance {
			left = mid + 1
		} else {
			right = mid
		}
	}
	insertIndex := left

	// insert at the position
	result := make([]connection, len(distances)+1)
	copy(result, distances[:insertIndex])
	result[insertIndex] = c
	copy(result[insertIndex+1:], distances[insertIndex:])

	// trim result to not exceed wanted shortest connection count
	end := min(loopCount, len(result))
	return result[:end]
}

func assignJunctionsToCircuits(j []junction, shortestConnections []connection, countOfAssignments int) {
	currentCircuit := 1
	for _, c := range shortestConnections {
		// fmt.Printf("\tRemaining assignments: %d\n", countOfAssignments)
		if countOfAssignments <= 0 {
			break
		}
		jidx0 := c.junctionIdx[0]
		jidx1 := c.junctionIdx[1]

		// fmt.Printf("Comparing junctions %v and %v\n", j[jidx0], j[jidx1])
		if j[jidx0].connected != j[jidx1].connected { // Only one of these is connected, connect the other one
			if j[jidx0].connected {
				// fmt.Printf("\tfirst connected! to circuit %d\n", j[jidx0].circuitNumber)
				j[jidx1].circuitNumber = j[jidx0].circuitNumber
				j[jidx1].connected = true
				countOfAssignments--
			} else {
				// fmt.Printf("\tseconds connected! to circuit %d\n", j[jidx1].circuitNumber)
				j[jidx0].circuitNumber = j[jidx1].circuitNumber
				j[jidx0].connected = true
				countOfAssignments--
			}
			continue
		}

		if j[jidx0].connected && j[jidx1].connected { // both are connected
			// fmt.Printf("\tboth connected!\n")
			if j[jidx0].circuitNumber == j[jidx1].circuitNumber { // and both are on the same circuit already, do nothing
				// fmt.Printf("\tto same circuit already: %d\n", j[jidx0].circuitNumber)
				countOfAssignments--
				continue
			}

			// both connected, but different circuits!
			// convert one of them into the other

			// fmt.Printf("I am joining circuit %d (disappears) and %d (merged)\n", j[jidx0].circuitNumber, j[jidx1].circuitNumber)
			joinCircuits(j[jidx0].circuitNumber, j[jidx1].circuitNumber, j)
			countOfAssignments--
			continue
		}

		// neither are connected!
		// fmt.Printf("\tneither connected! connecting to circuit %d\n", currentCircuit)
		j[jidx0].circuitNumber = currentCircuit
		j[jidx1].circuitNumber = currentCircuit
		j[jidx0].connected = true
		j[jidx1].connected = true
		countOfAssignments--
		currentCircuit++
	}
}

func joinCircuits(circuitA, circuitB int, j []junction) {
	// fmt.Printf("Joining circuits %d and %d\n", circuitA, circuitB)
	for i := range j {
		if j[i].circuitNumber == circuitA {
			j[i].circuitNumber = circuitB
		}
	}
}

func circuitSizes(j []junction) map[int]int {
	circuitSizes := make(map[int]int)
	for i := range j {
		if j[i].circuitNumber == 0 {
			continue
		}
		circuitSizes[j[i].circuitNumber]++
	}
	return circuitSizes
}

// part2 solves part 2 of the day's challenge.
// Keep this function signature intact for unit tests to work seamlessly.
func part2(input string, test bool) int {
	junctions := parseInput(input)
	fmt.Println(junctions)

	fmt.Println("")
	fmt.Println("Finding all distances sorted increasing")
	allDistances := findAllDistancesSortedIncreasing2(junctions)

	fmt.Println("Assigning junctions to circuits")
	lastA, lastB := assignJunctionsToCircuits2(junctions, allDistances)

	fmt.Printf("lastA: %v, lastB: %v\n", lastA, lastB)

	return lastA.position.x * lastB.position.x
}

func findAllDistancesSortedIncreasing2(j []junction) (distances []connection) {
	for i := range j { // loop over every junction
		if i >= len(j)-1 { // last element
			break
		}
		for ii := i + 1; ii < len(j); ii++ { // loop over every junction after it (previous ones already have connections calculated)
			c := connection{junctionIdx: [2]int{i, ii}, distance: euclideanDistance(j[i], j[ii])}
			distances = insertSortedDistance(distances, c, math.MaxInt64)
		}
	}

	return distances
}

func assignJunctionsToCircuits2(j []junction, connections []connection) (lastA, lastB junction) {
	uniqueCircuits := len(j)

	currentCircuit := 1
	for _, c := range connections {
		// fmt.Printf("\tAmount of circuits: %d\n", uniqueCircuits)
		jidx0 := c.junctionIdx[0]
		jidx1 := c.junctionIdx[1]
		// if uniqueCircuits <= 1 {
		// 	return j[connections[cIdx-1].junctionIdx[0]], j[connections[cIdx-1].junctionIdx[1]]
		// }

		// fmt.Printf("Comparing junctions %v and %v\n", j[jidx0], j[jidx1])
		if j[jidx0].connected != j[jidx1].connected { // Only one of these is connected, connect the other one
			if j[jidx0].connected {
				// fmt.Printf("\tfirst connected! to circuit %d\n", j[jidx0].circuitNumber)
				j[jidx1].circuitNumber = j[jidx0].circuitNumber
				j[jidx1].connected = true
				uniqueCircuits--
				if uniqueCircuits == 1 {
					return j[jidx0], j[jidx1]
				}
			} else {
				// fmt.Printf("\tseconds connected! to circuit %d\n", j[jidx1].circuitNumber)
				j[jidx0].circuitNumber = j[jidx1].circuitNumber
				j[jidx0].connected = true
				uniqueCircuits--
				if uniqueCircuits == 1 {
					return j[jidx0], j[jidx1]
				}
			}
			continue
		}

		if j[jidx0].connected && j[jidx1].connected { // both are connected
			// fmt.Printf("\tboth connected!\n")
			if j[jidx0].circuitNumber == j[jidx1].circuitNumber { // and both are on the same circuit already, do nothing
				// fmt.Printf("\tto same circuit already: %d\n", j[jidx0].circuitNumber)
				continue
			}

			// both connected, but different circuits!
			// convert one of them into the other

			// fmt.Printf("I am joining circuit %d (disappears) and %d (merged)\n", j[jidx0].circuitNumber, j[jidx1].circuitNumber)
			joinCircuits(j[jidx0].circuitNumber, j[jidx1].circuitNumber, j)
			uniqueCircuits--
			if uniqueCircuits == 1 {
				return j[jidx0], j[jidx1]
			}
			continue
		}

		// neither are connected!
		// fmt.Printf("\tneither connected! connecting to circuit %d\n", currentCircuit)
		j[jidx0].circuitNumber = currentCircuit
		j[jidx1].circuitNumber = currentCircuit
		j[jidx0].connected = true
		j[jidx1].connected = true
		uniqueCircuits--
		if uniqueCircuits == 1 {
			return j[jidx0], j[jidx1]
		}
		currentCircuit++
	}

	return junction{}, junction{}
}

func parseInput(input string) []junction {
	var junctions []junction
	for line := range strings.SplitSeq(input, "\n") {
		junctions = append(junctions, lineToJunction(line))
	}
	return junctions
}

func lineToJunction(input string) junction {
	var j junction

	position := strings.Split(input, ",")
	j.position.x = stringToInt(position[0])
	j.position.y = stringToInt(position[1])
	j.position.z = stringToInt(position[2])

	return j
}

func stringToInt(input string) int {
	output, err := strconv.Atoi(input)
	if err != nil {
		panic(fmt.Sprintf("stringToInt: failed to convert %q to int: %v", input, err))
	}
	return output
}
