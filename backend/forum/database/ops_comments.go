package database

import (
	"log"
	"time"
)

// AddComment adds comment to database, returns id of new comment
func (db DB) AddComment(content string, postId, authorId int) int {
	result, err := db.Exec(`INSERT INTO comments (content, post_id, author_id, date, likes_count, dislikes_count) 
								  VALUES (?, ?, ?, ?, ?, ?)`, content, postId, authorId, time.Now().Format(time.RFC3339), 0, 0)
	if err != nil {
		log.Panic(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Panic(err)
	}
	return int(id)
}

func (db DB) UpdatePostsCommentsCount(postId int, change int) {
	_, err := db.Exec("UPDATE posts SET comments_count = comments_count + ? WHERE id = ?", change, postId)
	if err != nil {
		log.Panic(err)
	}
}
