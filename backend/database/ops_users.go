package database

import (
	"log"
)

// GetUserById returns user with specified id
//
// returns nil if user not found
func (db DB) GetUserById(id int) *User {
	query, err := db.Query("SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		log.Panic(err)
	}
	var user User
	if !query.Next() {
		return nil
	}
	err = query.Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		log.Panic(err)
	}
	return &user
}

// AddUser adds user to database, returns id of new user
func (db DB) AddUser(user User) int64 {
	result, err := db.Exec(`INSERT INTO users (name, email, password) VALUES (?, ?, ?)`, user.Name, user.Email, user.Password)
	if err != nil {
		log.Panic(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Panic(err)
	}
	return id
}
