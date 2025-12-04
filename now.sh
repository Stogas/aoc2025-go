#!/usr/bin/env bash

# Print current date/time in Europe/Vilnius in RFC 3339 format, portable between macOS and Linux
if [[ "$(uname)" == "Darwin" ]]; then
	# macOS: use gdate if available, else fallback to date with compatible format
	if command -v gdate >/dev/null 2>&1; then
		TZ=Europe/Vilnius gdate --rfc-3339=seconds
	else
		TZ=Europe/Vilnius date '+%Y-%m-%d %H:%M:%S%z' | \
			sed -E 's/([0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2})([+-][0-9]{2})([0-9]{2})/\1\2:\3/'
	fi
else
	# Linux
	TZ=Europe/Vilnius date --rfc-3339=seconds
fi
