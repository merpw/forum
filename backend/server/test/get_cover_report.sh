#!/bin/sh

# Run unit tests for forum/server/test package with coverage analysis enabled
go test "forum/server/test" -cover -coverpkg="../../..." -coverprofile="./cover_reports/profile.txt"

# Generate an HTML coverage report using coverage profile generated in previous step
go tool cover -html="./cover_reports/profile.txt" -o "./cover_reports/coverage.html"

# Open the coverage report in a browser window
open ./cover_reports/coverage.html
