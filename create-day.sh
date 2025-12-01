#!/usr/bin/env bash

if [[ ! -n $1 ]]; then
  echo "No day num provided" >&2
	exit 2
fi

dayDir="day${1}"

if [[ ! -d "${dayDir}" ]]; then
	cp -r skeleton "${dayDir}"
	# Portable sed in-place delete for lines 1-2
	if sed --version >/dev/null 2>&1; then
		# GNU sed (Linux)
		sed -i '1,2d' "${dayDir}/main.go"
	else
		# BSD sed (MacOS)
		sed -i '' '1,2d' "${dayDir}/main.go"
	fi
else
	echo "Day directory already exists" >&2
	exit 3
fi
