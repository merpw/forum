package main

import (
	"database/sql"
	"flag"
	"fmt"
	"forum/database/migrations"
	"log"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
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
			err := migrations.Check(db)
			if err != nil {
				log.Fatal(err)
			}
			version, err := migrations.GetVersion(db)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Integrity check passed successfully, current revision is %d\n", version)
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
			toRevision, err := strconv.Atoi(toRevisionStr)
			if err != nil || toRevision <= 0 {
				fmt.Println("ERROR: revision should be positive integer number")
			}

			err = migrations.Migrate(db, toRevision)
			if err != nil {
				log.Fatalf("migration failed, %s\n", err)
			}
			fmt.Println("Migration finished.")
		},
	},
}

func main() {
	dbFile := flag.String("db", "", "specify database to migrate")

	flag.Usage = func() {
		fmt.Printf(`Forum database CLI

Usage: cli [PARAMETERS] [ACTION]

Availible parameters:
`)
		flag.VisitAll(func(f *flag.Flag) {
			if f.DefValue != "" {
				fmt.Printf("\t-%s, default: %s - %s\n", f.Name, f.DefValue, f.Usage)
			} else {
				fmt.Printf("\t-%s - %s\n", f.Name, f.Usage)
			}
		})

		fmt.Println("\nAvailable actions:")
		for _, action := range actions {
			fmt.Printf("\t%s - %s\n", action.command, action.description)
		}
	}

	flag.Parse()

	if *dbFile == "" {
		log.Fatal("ERROR: Database file is not specified")
	}

	db, err := sql.Open("sqlite3", *dbFile+"?_foreign_keys=true") // enable foreign keys
	if err != nil {
		log.Fatal(err)
	}

	action := strings.ToLower(flag.Arg(0))

	for _, act := range actions {
		if action == act.command {
			act.action(db)
			return
		}
	}

	version, err := migrations.GetVersion(db)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("version", version)

	err = migrations.Migrate(db, 2)
	if err != nil {
		log.Fatal(err)
	}
}
