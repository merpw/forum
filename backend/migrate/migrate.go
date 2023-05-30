package migrate

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// Migration defines operations that should be done to move Up or Down between revisions
type Migration struct {
	Up   func(db *sql.DB) error
	Down func(db *sql.DB) error
}

type Migrations []Migration

// Latest is latest available revision
func (migrations Migrations) Latest() int {
	return len(migrations)
}

// Migrate migrates the database to the specified revision
//
// If the database is empty, it will create the initial schema
func (migrations Migrations) Migrate(db *sql.DB, toRevision int) error {
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

	if toRevision > migrations.Latest() || toRevision < 1 {
		return fmt.Errorf("invalid revision %d, availible revisions 1..%d", toRevision, migrations.Latest())
	}

	if dbRevision > migrations.Latest() || dbRevision < 0 {
		return fmt.Errorf("database revision %d is not supported, supported revisions 1..%d",
			dbRevision, migrations.Latest())
	}

	if dbRevision == 0 {
		log.Printf("Empty/unknown database, initial schema will be created")
	}

	if dbRevision > toRevision {
		log.Printf(`WARNING: migrating database DOWN from %d to %d
Down migrations may remove columns and tables, so some data can be lost. Make sure you have a backup.
If you still want to continue, type YES, otherwise press Ctrl+C to abort.
`, dbRevision, toRevision)
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
		if strings.TrimSpace(text) != "YES" {
			return fmt.Errorf("aborted by user")
		}
		for i := dbRevision; i > toRevision; i-- {
			err := migrations[i-1].Down(db)
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
			err := migrations[i].Up(db)
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
