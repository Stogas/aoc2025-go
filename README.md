# aoc2025-go

[Advent of Code 2025](https://adventofcode.com/2025) in Go, using just the Go v1.25 standard library.

Some best practices are deliberately ignored here. For example - I use `panic()` for any errors, [use `init()` functions](https://github.com/leighmcculloch/gochecknoinits), etc.


Each day's challenge is completely self-contained in its directory, with no efforts made to deduplicate any code (other than having a skeleton template), update older challenges with any improvements made to the template in newer challenges, etc.

#### Benchmark Results

Each day's directory contains benchmark results (`bench.md`) generated using [aoc-bench-script](https://github.com/Stogas/aoc-bench-script). These files show performance measurements for both parts of the day's solution, including system info and timing statistics.

AI is not used for solving daily challenges, and autocomplete suggestions are disabled.

My goal is not to solve the challenge as fast as possible - it's to solve it in a reasonable amount of time while keeping performance good (i.e. no slow bruteforcing - preferably a solution will be printed in less than a second, and optimally faster than 100ms).

This repo also uses [devcontainers](https://blog.stogas.dev/2025/10/18/devcontainers/).

#### Generating a daily challenge directory from the skeleton template

Run `./create-day.sh <N>` where `N` is the day number

#### Running the code for a daily challenge

1. Create `dayN/test.txt`
2. Create `dayN/input.txt`
3. Run: `go run dayN/main.go -part <1 or 2> [-test]`

`-test` will run the `test.txt` inputs instead of `input.txt`

#### Testing a daily challenge

Each day carries a lightweight unit test suite that exercises the sample input (`test.txt`).

1. Populate `dayN/test.txt` with the sample input.
2. Add the expected answers to `dayN/test_results.txt` using the format:

	```
	part1=24000
	part2=45000
	```

	Use `?` (or `-`) to skip a part until you have a solution, e.g. `part2=?`.
3. Run `go test ./dayN` (or simply `go test ./...` from the repo root) to validate the implementation.

`go test` is also wired into `lefthook`, so commits will fail fast if a day's sample answers drift from what is saved in `test_results.txt`.

#### Compiling and running

Almost the same as the previous section, but compiles ahead of time to measure real performance. Can be used for comparing with other solutions for the same daily challenge.

This has the added benefit of actually embedding the input data into the compiled binary, allowing to share the solution program in just a single binary executable.

1. Create `dayN/test.txt`
2. Create `dayN/input.txt`
3. `mkdir bin/`
4. Build: `go build -o bin/dayN dayN/main.go`
5. Run: `bin/dayN -part <1 or 2> [-test]`

`-test` will run the `test.txt` inputs instead of `input.txt`

