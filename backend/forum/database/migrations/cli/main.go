package main

import (
	"backend/forum/database/migrations"
	"backend/migrate/cli"
)

func main() {
	cli.Main(migrations.Migrations)
}
