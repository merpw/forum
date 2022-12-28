package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// return a pointer to a sql.DB object
//
// how to use:
// db := Opendb()
// Special note. No need to defer close a connection as it is automatically closed by sql driver
func Opendb() *sql.DB {
	db, err := sql.Open("sqlite3", "./database/databsase.db")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
