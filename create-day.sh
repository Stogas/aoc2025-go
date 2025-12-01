#!/usr/bin/env bash

if [[ ! -n $1 ]]; then
  echo "No day num provided" >&2
	exit 2
fi

dayDir="day${1}"

if [[ ! -d "${dayDir}" ]]; then
	cp -r skeleton "${dayDir}"
	sed -i '' '1,2d' "${dayDir}/main.go"
	> "${dayDir}/input.txt"
	> "${dayDir}/test.txt"
else
	echo "Day directory already exists" >&2
	exit 3
fi
