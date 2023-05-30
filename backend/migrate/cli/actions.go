package cli

import (
	"backend/migrate"
	"database/sql"
	"flag"
	"log"
	"strconv"
)

var actions = []struct {
	command     string
	description string
	action      func(db *sql.DB, migrations migrate.Migrations)
}{
	{
		"stat",
		"- check database integrity and show current revision",
		func(db *sql.DB, migrations migrate.Migrations) {
			err := migrate.Check(db)
			if err != nil {
				log.Fatal(err)
			}
			version, err := migrate.GetVersion(db)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Integrity check passed successfully, current revision is %d, latest is %d\n",
				version, migrations.Latest())
		},
	},
	{
		"migrate",
		"[REVISION] - migrate to the specific revision or `latest`",
		func(db *sql.DB, migrations migrate.Migrations) {
			toRevisionStr := flag.Arg(1)
			if toRevisionStr == "" {
				log.Fatal("ERROR: desired revision is not provided")
			}
			var toRevision int
			if toRevisionStr == "latest" {
				toRevision = migrations.Latest()
			} else {
				toRevision, _ = strconv.Atoi(toRevisionStr)
				if toRevision <= 0 {
					log.Fatal("ERROR: revision should be a positive integer number")
				}
			}

			err := migrations.Migrate(db, toRevision)
			if err != nil {
				log.Fatalf("ERROR: migration failed, %s\n", err)
			}
			log.Println("Migration finished.")
		},
	},
	{
		"create",
		"- create new database file and migrate it to the latest revision",
		func(db *sql.DB, migrations migrate.Migrations) {
			err := migrations.Migrate(db, migrations.Latest())
			if err != nil {
				log.Fatal(err)
			}
			log.Println("Database created successfully")
		},
	},
}
