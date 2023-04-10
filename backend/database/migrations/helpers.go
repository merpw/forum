package migrations

import (
	"database/sql"
	"fmt"
)

// GetVersion returns the current database version (user_version pragma)
func GetVersion(db *sql.DB) (int, error) {
	var version int
	err := db.QueryRow("PRAGMA user_version").Scan(&version)
	if err != nil {
		return 0, err
	}
	return version, nil
}

// SetVersion sets the current database version (user_version pragma)
func SetVersion(db *sql.DB, version int) error {
	_, err := db.Exec(fmt.Sprintf("PRAGMA user_version = %d", version))
	return err
}

// Check checks the database for errors (integrity_check pragma)
func Check(db *sql.DB) error {
	_, err := db.Exec("PRAGMA integrity_check")
	return err
}
