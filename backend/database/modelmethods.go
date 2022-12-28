package database

import (
	"log"
	"time"
)

// GetPosts reads all posts from database
//
// panics if error occurs
func (db DB) GetPosts() []Post {
	q, err := db.Query("SELECT * FROM posts")
	if err != nil {
		log.Panic(err)
	}
	var posts []Post
	for q.Next() {
		var post Post
		err = q.Scan(&post.Id, &post.Title, &post.Content, &post.Author, &post.Date, &post.Likes, &post.Dislikes, &post.UsersReactions, &post.CommentsCount)
		if err != nil {
			log.Panic(err)
		}
		posts = append(posts, post)
	}
	return posts
}

func (db DB) GetPostById(id int) *Post {
	query, err := db.Query("SELECT * FROM posts WHERE id = ?", id)
	if err != nil {
		log.Panic(err)
	}
	var post Post
	if !query.Next() {
		return nil
	}
	err = query.Scan(&post.Id, &post.Title, &post.Content, &post.Author, &post.Date, &post.Likes, &post.Dislikes, &post.UsersReactions, &post.CommentsCount)
	if err != nil {
		log.Panic(err)
	}
	return &post
}

// AddPost adds post to database, returns id of new post
func (db DB) AddPost(post Post) int64 {
	result, err := db.Exec(`INSERT INTO posts (title, content, author, date, likes, dislikes, user_reactions, comments_count) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, post.Title, post.Content, post.Author, time.Now().Format(time.RFC3339), 0, 0, "", 0)
	if err != nil {
		log.Panic(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Panic(err)
	}
	return id
}
