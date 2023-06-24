package main

import (
	"backend/chat/database/migrations"
	"backend/migrate/cli"
)

func main() {
	cli.Main(migrations.Migrations)
}
