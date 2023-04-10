package migrations

import (
	"database/sql"
	"fmt"
	"log"
)

type Migration struct {
	Up   func(db *sql.DB) error
	Down func(db *sql.DB) error
}

var Migrations = []Migration{
	v001,
	v002,
}

// Migrate migrates the database to the specified revision
//
// If the database is empty, it will create the initial schema
func Migrate(db *sql.DB, toRevision int) error {
	err := Check(db)
	if err != nil {
		return fmt.Errorf("database precheck failed: %w", err)
	}

	dbRevision, err := GetVersion(db)
	if err != nil {
		return err
	}

	if dbRevision == toRevision {
		return nil
	}

	if toRevision > len(Migrations) || toRevision < 1 {
		return fmt.Errorf("invalid revision %d, availible revisions 1..%d", toRevision, len(Migrations))
	}

	if dbRevision == 0 {
		log.Printf("Empty/unknown database, initial schema will be created")
	}

	if dbRevision > toRevision {
		log.Printf("Migrating database down from revision %d to %d", dbRevision, toRevision)
		for i := dbRevision; i > toRevision; i-- {
			err := Migrations[i-1].Down(db)
			if err != nil {
				return err
			}
			err = SetVersion(db, i-1)
			if err != nil {
				return err
			}
		}
	} else {
		log.Printf("Migrating database up from revision %d to %d", dbRevision, toRevision)
		for i := dbRevision; i < toRevision; i++ {
			err := Migrations[i].Up(db)
			if err != nil {
				return err
			}
			err = SetVersion(db, i+1)
			if err != nil {
				return err
			}
		}
	}

	err = Check(db)
	if err != nil {
		return fmt.Errorf("database postcheck failed: %w", err)
	}

	return nil
}
