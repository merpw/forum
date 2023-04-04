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
	err = query.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Age, &user.Gender)
	if err != nil {
		log.Panic(err)
	}
	query.Close()

	return &user
}

// GetUserByLogin returns user with specified login
//
// returns nil if user not found
func (db DB) GetUserByLogin(login string) *User {
	query, err := db.Query("SELECT * FROM users WHERE email = ? OR name = ? COLLATE NOCASE", login, login)
	if err != nil {
		log.Panic(err)
	}

	if !query.Next() {
		return nil
	}
	var user User
	err = query.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Age, &user.Gender)
	if err != nil {
		log.Panic(err)
	}
	query.Close()

	return &user
}

// AddUser adds user to database, returns id of new user
func (db DB) AddUser(name, email, password, first_name, last_name, age, gender string) int {
	result, err := db.Exec(`INSERT INTO users (name, email, password, first_name, last_name, age, gender) VALUES (?, ?, ?, ?, ?, ?, ?)`, name, email, password, first_name, last_name, age, gender)
	if err != nil {
		log.Panic(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Panic(err)
	}
	return int(id)
}

// IsEmailTaken checks if email is already in use
func (db DB) IsEmailTaken(email string) bool {
	query, err := db.Query("SELECT 1 FROM users WHERE email = ?", email)
	if err != nil {
		log.Panic(err)
	}
	isTaken := query.Next()
	query.Close()

	return isTaken
}

// IsNameTaken checks if name is already in use.
// Added because authorization supports the name and the email as login
func (db DB) IsNameTaken(name string) bool {
	query, err := db.Query("SELECT 1 FROM users WHERE name = ? COLLATE NOCASE", name)
	if err != nil {
		log.Panic(err)
	}
	isTaken := query.Next()
	query.Close()

	return isTaken
}
