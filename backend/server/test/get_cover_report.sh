#!/bin/sh

mkdir -p ./cover_reports

go test "forum/server/test" -cover -coverpkg="../../..." -coverprofile="./cover_reports/profile.txt"

go tool cover -html="./cover_reports/profile.txt" -o "./cover_reports/coverage.html"

if [[ $(uname) == "Darwin" ]]; then
    open ./cover_reports/coverage.html
    exit 0
fi

if [[ $(uname) == "Linux" ]]; then
    xdg-open ./cover_reports/coverage.html
    exit 0
fi

if [[ $(uname -r) == *microsoft* ]]; then
    start msedge ./cover_reports/coverage.html
    exit 0
fi

echo "Could not open coverage report in browser. Please open ./cover_reports/coverage.html manually."