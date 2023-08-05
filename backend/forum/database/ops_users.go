package database

import (
	"database/sql"
	"log"
)

// GetAllUserIds returns slice of all users ids
func (db DB) GetAllUserIds() (userIds []int) {
	query, err := db.Query("SELECT id FROM users ORDER BY username COLLATE NOCASE ASC")
	if err != nil {
		log.Panic(err)
	}

	for query.Next() {
		var userId int
		err = query.Scan(&userId)
		if err != nil {
			log.Panic(err)
		}
		userIds = append(userIds, userId)
	}
	if err = query.Close(); err != nil {
		log.Panic(err)
	}

	return
}

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
	err = query.Scan(
		&user.Id, &user.Username, &user.Email, &user.Password,
		&user.FirstName, &user.LastName, &user.DoB, &user.Gender, &user.Avatar, &user.Bio)
	if err != nil {
		log.Panic(err)
	}
	if err = query.Close(); err != nil {
		log.Panic(err)
	}

	return &user
}

func (db DB) GetLastUserId() int {
	query, err := db.Query("SELECT id FROM users ORDER BY id DESC LIMIT 1")
	if err != nil {
		log.Panic(err)
	}

	if !query.Next() {
		return 0
	}

	var userId int
	err = query.Scan(&userId)
	if err != nil {
		log.Panic(err)
	}

	if err = query.Close(); err != nil {
		log.Panic(err)
	}

	return userId
}

// GetUserByLogin returns user with specified login
//
// returns nil if user not found
func (db DB) GetUserByLogin(login string) *User {
	query, err := db.Query("SELECT * FROM users WHERE email = ? OR username = ? COLLATE NOCASE", login, login)
	if err != nil {
		log.Panic(err)
	}

	if !query.Next() {
		return nil
	}
	var user User
	err = query.Scan(
		&user.Id, &user.Username, &user.Email, &user.Password,
		&user.FirstName, &user.LastName, &user.DoB, &user.Gender, &user.Avatar, &user.Bio)
	if err != nil {
		log.Panic(err)
	}

	if err = query.Close(); err != nil {
		log.Panic(err)
	}

	return &user
}

// AddUser adds user to database, returns id of new user
func (db DB) AddUser(
	username, email, password string,
	firstName, lastName, dob, gender, bio, avatar sql.NullString) int {

	result, err := db.Exec(
		`INSERT INTO users (username, email, password, first_name, last_name, dob, gender, bio, avatar)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		username,
		email,
		password,
		firstName.String,
		lastName.String,
		dob.String,
		gender.String,
		avatar,
		bio,
	)
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
	if err = query.Close(); err != nil {
		log.Panic(err)
	}

	return isTaken
}

// IsNameTaken checks if name is already in use.
// Added because authorization supports the name and the email as login
func (db DB) IsNameTaken(username string) bool {
	query, err := db.Query("SELECT 1 FROM users WHERE username = ? COLLATE NOCASE", username)
	if err != nil {
		log.Panic(err)
	}
	isTaken := query.Next()

	if err = query.Close(); err != nil {
		log.Panic(err)
	}

	return isTaken
}
