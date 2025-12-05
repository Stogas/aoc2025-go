#!/usr/bin/env bash

# Usage: ./bench.sh <N>
# Build and benchmark the program in ./dayN using https://github.com/Stogas/aoc-bench-script

set -euo pipefail

if [ "$#" -ne 1 ]; then
	echo "Usage: $0 <N>  (where N is the day number, e.g. 1)"
	exit 2
fi

DAY="$1"

if ! [[ "$DAY" =~ ^[0-9]+$ ]]; then
	echo "Day must be a single integer"
	exit 2
fi

DAY_DIR="day${DAY}"
if [ ! -d "$DAY_DIR" ]; then
	echo "Directory '$DAY_DIR' does not exist"
	exit 1
fi

# location to store the downloaded bench script (ignored by git)
DOWNLOAD_DIR=".aoc-bench-script"
DOWNLOAD_SCRIPT="$DOWNLOAD_DIR/bench.sh"
mkdir -p "$DOWNLOAD_DIR"

if [ -f "$DOWNLOAD_SCRIPT" ]; then
	echo "Bench script already exists at $DOWNLOAD_SCRIPT."
	echo "If you need to update it, delete the file and re-run this script."
else
	echo "Downloading bench script to $DOWNLOAD_SCRIPT..."
	curl -sSL https://raw.githubusercontent.com/Stogas/aoc-bench-script/refs/heads/main/bench.sh -o "$DOWNLOAD_SCRIPT"
	chmod +x "$DOWNLOAD_SCRIPT"
fi
echo

# build the day's binary
mkdir -p bin
BIN_PATH="bin/day${DAY}"

if [ -f "$BIN_PATH" ]; then
	echo "Binary already exists at $BIN_PATH"
	read -p "Would you like to recompile? (y/n) " -n 1 -r
	echo
	if [[ $REPLY =~ ^[Yy]$ ]]; then
		echo "Building ${DAY_DIR}/main.go -> ${BIN_PATH}"
		go build -o "$BIN_PATH" "${DAY_DIR}/main.go"
	else
		echo "Using existing binary at $BIN_PATH"
	fi
else
	echo "Building ${DAY_DIR}/main.go -> ${BIN_PATH}"
	go build -o "$BIN_PATH" "${DAY_DIR}/main.go"
fi

OUTPUT_FILE="${DAY_DIR}/bench.md"
echo "# Day ${DAY} - bench run: $(date -u +"%Y-%m-%dT%H:%M:%SZ")" >> "$OUTPUT_FILE"
echo "Benchmarked using https://github.com/Stogas/aoc-bench-script" > "$OUTPUT_FILE"

for PART in 1 2; do
	# add a readable section header with a blank line before it
	printf '\n%s\n\n' "# Day ${DAY} Part ${PART}" >> "$OUTPUT_FILE"
	echo "Running bench for part ${PART} as \"$DOWNLOAD_SCRIPT\" -q 100 2 \"${BIN_PATH} -part ${PART}\"..."
	# the bench script expects a single string command to run; pass the binary with -part
	# run the downloaded script with iterations=100 and warmups=2
	"$DOWNLOAD_SCRIPT" -q 100 2 "${BIN_PATH} -part ${PART}" >> "$OUTPUT_FILE" 2>&1
done

echo
echo "Benchmark complete. Output written to $OUTPUT_FILE"
