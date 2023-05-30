package cli

import (
	"backend/migrate"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	log.SetFlags(0)
	flag.Usage = func() {
		fmt.Printf(`Forum database migrations CLI

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
			fmt.Printf("\t%s %s\n", action.command, action.description)
		}
	}
}

// Main can be used to run specified migrate.Migrations from CLI
func Main(migrations migrate.Migrations) {
	dbFile := flag.String("db", "", "specify database file to work with")
	flag.Parse()

	if *dbFile == "" {
		log.Fatal("ERROR: Database file is not specified")
	}

	action := strings.ToLower(flag.Arg(0))
	if action == "" {
		flag.Usage()
		return
	}
	actionNum := -1
	for i, act := range actions {
		if action == act.command {
			actionNum = i
			break
		}
	}
	if actionNum == -1 {
		log.Fatalf("ERROR: Unknown action %s", action)
	}

	_, err := os.Stat(*dbFile)
	if err != nil {
		if os.IsNotExist(err) {
			if action != "create" {
				log.Fatalf("ERROR: Database file %s does not exist. Run `create` action to create database", *dbFile)
			}
		} else {
			log.Fatal(err)
		}
	} else {
		if action == "create" {
			log.Fatalf("ERROR: %s already exists", *dbFile)
		}
	}

	db, err := sql.Open("sqlite3", *dbFile+"?_foreign_keys=true") // enable foreign keys
	if err != nil {
		log.Fatal(err)
	}

	actions[actionNum].action(db, migrations)
}
