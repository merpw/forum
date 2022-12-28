package database

import (
	"log"
	"time"
)

// AddSession adds session to database
func (db DB) AddSession(token string, expire, userId int) {
	_, err := db.Exec(`INSERT INTO sessions (token, expire, user_id) VALUES (?, ?, ?)`, token, expire, userId)
	if err != nil {
		log.Panic(err)
	}
}

// CheckSession checks if session is valid
//
// returns User if session is valid or nil if not
func (db DB) CheckSession(token string) *User {
	query, err := db.Query("SELECT * FROM sessions WHERE token = ?", token)
	if err != nil {
		log.Panic(err)
	}

	if !query.Next() {
		return nil
	}
	var session Session
	err = query.Scan(&session.Id, &session.Token, &session.Expire, &session.UserId)
	if err != nil {
		log.Panic(err)
	}
	query.Close()

	if session.Expire < int(time.Now().Unix()) {
		db.RemoveExpiredSessions()
		return nil
	}
	return db.GetUserById(session.UserId)
}

func (db DB) RemoveExpiredSessions() {
	_, err := db.Exec("DELETE FROM sessions WHERE expire < ?", time.Now().Unix())
	if err != nil {
		log.Panic(err)
	}
}

func (db DB) RemoveSession(token string) {
	_, err := db.Exec("DELETE FROM sessions WHERE token = ?", token)
	if err != nil {
		log.Panic(err)
	}
}
