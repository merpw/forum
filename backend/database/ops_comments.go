package database

import (
	"log"
	"time"
)

// AddComment adds comment to database, returns id of new comment
func (db DB) AddComment(content string, authorId int) int {
	result, err := db.Exec(`INSERT INTO comments (post_id, author_id, content, date, likes_count, dislikes_count) 
								  VALUES (?, ?, ?, ?, ?)`, content, authorId, time.Now().Format(time.RFC3339), 0, 0)
	if err != nil {
		log.Panic(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Panic(err)
	}
	return int(id)
}
