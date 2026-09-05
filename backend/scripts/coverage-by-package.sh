#!/usr/bin/env bash
# Runs the backend test suite with coverage and prints an accurate per-package
# breakdown, plus the merged total. Run from backend/: ./scripts/coverage-by-package.sh
#
# Why this exists (don't use `go test -cover -coverpkg=...` output directly):
#
# 1. `go test -cover -coverpkg=./internal/...,./pkg/... ./...` prints one
#    `coverage: X%` line per TEST package (e.g. `backend/test/resolvers`), not
#    per SOURCE package. Because tests here live under `backend/test/*`
#    (black-box, separate from the code they exercise) and -coverpkg spans the
#    whole codebase, each of those lines is that one test binary's coverage of
#    the ENTIRE internal+pkg tree — almost always a tiny, meaningless-looking
#    percentage. It also prints an unlabeled `coverage: 0.0%` line for every
#    package `-coverpkg` instruments but that has no tests of its own, which
#    reads like a failure but isn't. None of this tells you which SOURCE
#    package actually needs more tests.
#
# 2. The natural fix — `go tool cover -func=profile.out` grouped by package —
#    runs into a second trap: `go test -coverprofile=X ./...` with `-coverpkg`
#    does NOT deduplicate across test binaries. Every one of the ~15+ test
#    packages compiles its own coverage-instrumented copy of the whole
#    -coverpkg set and appends its own full block listing to the same profile
#    file — so most (file, line-range) blocks appear 15+ times. Naively
#    summing `numStmt` per block (e.g. a quick per-package awk sum) inflates
#    the statement count by roughly that same factor and produces garbage
#    (~1.6% when the real number is ~23%). `go tool cover -func`'s own total
#    line is correct because it deduplicates internally — but it has no
#    per-package grouping, only per-function.
#
# This script does both correctly: dedupe each (file, line-range) block by its
# profile key before counting it (a block is "covered" if ANY duplicate saw a
# nonzero count), THEN group by source package directory. Verified against
# `go tool cover -func`'s total: this produces the same number.

set -euo pipefail

PROFILE="$(mktemp)"
RAW_OUTPUT="$(mktemp)"
trap 'rm -f "$PROFILE" "$RAW_OUTPUT"' EXIT

if ! go test -coverpkg=./internal/...,./pkg/... -coverprofile="$PROFILE" ./... > "$RAW_OUTPUT" 2>&1; then
	echo "go test failed — full output:" >&2
	cat "$RAW_OUTPUT"
	exit 1
fi

awk '
	NR == 1 { next }  # skip the "mode: set" header line
	{
		key = $1
		stmt = $2 + 0
		cnt = $3 + 0
		if (!(key in seen)) {
			seen[key] = 1
			numstmt[key] = stmt
			split(key, a, ":")
			file[key] = a[1]
		}
		if (cnt > 0) coveredflag[key] = 1
	}
	END {
		for (key in seen) {
			n = split(file[key], parts, "/")
			dir = parts[1]
			for (i = 2; i < n; i++) dir = dir "/" parts[i]
			total[dir] += numstmt[key]
			grandtotal += numstmt[key]
			if (key in coveredflag) {
				covered[dir] += numstmt[key]
				grandcovered += numstmt[key]
			}
		}
		for (d in total) {
			pct = (covered[d] + 0) / total[d] * 100
			printf "%6.1f%%  (%5d/%5d)  %s\n", pct, covered[d] + 0, total[d], d
		}
		printf "\n%6.2f%%  (%5d/%5d)  TOTAL\n", grandcovered / grandtotal * 100, grandcovered, grandtotal
	}
' "$PROFILE" | sed 's#github.com/CodeWarrior-debug/perspectize/backend/##' | sort
