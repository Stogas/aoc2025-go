package main

import (
	"bufio"
	_ "embed"
	"strconv"
	"strings"
	"testing"
)

const (
	part1Key  = "part1"
	part2Key  = "part2"
	skipToken = "?"
	partCount = 2
)

//go:embed test_results.txt
var testResultsRaw string

type expectation struct {
	value int
	skip  bool
}

// loadExpectations parses embedded test_results.txt into expectations per part.
func loadExpectations(t *testing.T) map[string]expectation { //nolint:gocognit // rule triggers here, but this is fine as its not meant for human readability.
	t.Helper()
	exps := map[string]expectation{
		part1Key: {skip: true},
		part2Key: {skip: true},
	}

	scanner := bufio.NewScanner(strings.NewReader(testResultsRaw))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", partCount)
		if len(parts) != partCount {
			t.Fatalf("invalid test_results.txt line %q, expected key=value", line)
		}

		key := strings.ToLower(strings.TrimSpace(parts[0]))
		val := strings.TrimSpace(parts[1])

		if val == "" || val == skipToken || val == "-" {
			exps[key] = expectation{skip: true}
			continue
		}

		parsed, err := strconv.Atoi(val)
		if err != nil {
			t.Fatalf("failed parsing %s=%q: %v", key, val, err)
		}
		exps[key] = expectation{value: parsed}
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("failed reading test_results.txt: %v", err)
	}

	return exps
}

// checkPart runs the given solver on testInput and compares to expectations.
func checkPart(t *testing.T, key string, solver func(string, bool) int) {
	t.Helper()
	exp := loadExpectations(t)[key]
	if exp.skip {
		t.Skipf("%s expectation not set in test_results.txt", key)
	}

	got := solver(testInput, true)
	if got != exp.value {
		t.Fatalf("%s: got %d, want %d", key, got, exp.value)
	}
}

// TestPart1Sample validates part1 against the sample expectations.
func TestPart1Sample(t *testing.T) {
	checkPart(t, part1Key, part1)
}

// TestPart2Sample validates part2 against the sample expectations.
func TestPart2Sample(t *testing.T) {
	checkPart(t, part2Key, part2)
}
