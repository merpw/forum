package main

import (
	"database/sql"
	"flag"
	"forum/database/migrate"
	"forum/database/migrate/migrations"
	"log"
	"strconv"
)

var actions = []struct {
	command     string
	description string
	action      func(db *sql.DB)
}{
	{
		"stat",
		"show information about database",
		func(db *sql.DB) {
			err := migrate.Check(db)
			if err != nil {
				log.Fatal(err)
			}
			version, err := migrate.GetVersion(db)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Integrity check passed successfully, current revision is %d, latest is %d\n",
				version, len(migrations.Migrations))
		},
	},
	{
		"migrate",
		"[REVISION], migrate to the specific revision",
		func(db *sql.DB) {
			toRevisionStr := flag.Arg(1)
			if toRevisionStr == "" {
				log.Fatal("ERROR: desired revision is not provided")
			}
			var toRevision int
			if toRevisionStr == "latest" {
				toRevision = len(migrations.Migrations)
			} else {
				toRevision, _ = strconv.Atoi(toRevisionStr)
				if toRevision <= 0 {
					log.Fatal("ERROR: revision should be a positive integer number")
				}
			}

			err := migrate.Migrate(db, toRevision)
			if err != nil {
				log.Fatalf("migration failed, %s\n", err)
			}
			log.Println("Migration finished.")
		},
	},
}
