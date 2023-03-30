#!/bin/sh

# make folder cover_reports if it doesn't exist
mkdir -p ./cover_reports

# Run unit tests for forum/server/test package with coverage analysis enabled
go test "forum/server/test" -cover -coverpkg="../../..." -coverprofile="./cover_reports/profile.txt"

# Generate an HTML coverage report using coverage profile generated in previous step
go tool cover -html="./cover_reports/profile.txt" -o "./cover_reports/coverage.html"

# Open the coverage report in a browser window
if [[ $(uname) == "Darwin" ]]; then
    open ./cover_reports/coverage.html
fi

if [[ $(uname) == "Linux" ]]; then
    xdg-open ./cover_reports/coverage.html
fi

if [[ $(uname -r) == *microsoft* ]]; then
    start msedge ./cover_reports/coverage.html
fi
