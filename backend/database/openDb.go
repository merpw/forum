package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// return a pointer to a sql.DB object
func Opendb() *sql.DB {
	db, err := sql.Open("sqlite3", "./database/databsase.db")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// how to use the Opendb function
// db, err := Opendb()
// if err != nil {
// 	log.Fatal(err)
// }
