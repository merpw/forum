# Step by step guide to test coverage in Go

- open backend folder in code with `code .`
- press <kbd>Ctrl</kbd>+<kbd>J</kbd> to open terminal
- run command `go test forum/server/test -cover -coverpkg=./... -coverprofile=forum/server/test/cover_reports/profile.txt` to run all tests and fill `profile.txt`.
- run command `go tool cover -html=forum/server/test/cover_reports/profile.txt -o forum/server/test/cover_reports/coverage.html` to see the coverage report
- open go file you want to see test coverage for
- press <kbd>Ctrl</kbd>+<kbd>Shift</kbd>+<kbd>P</kbd>
- type `Go: Apply Cover Profile`and select apply cover profile
- press <kbd>Enter</kbd>
- past the path of the `profile.txt` file
- press <kbd>Enter</kbd>
- you should see the coverage report for the file you are currently editing

## Bash script to generate test profile coverage report

Inside the backend directory, run the following command in the terminal:

```bash
cd ./server/test
bash get_cover_report.sh
```
