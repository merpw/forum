package main

import (
	"backend/migrate/cli"
	"backend/migrate/example"
)

func main() {
	cli.Main(example.Migrations)
}
