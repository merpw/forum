package database

import (
	"log"
	"time"
)

// GetPosts reads all posts from database
//
// panics if error occurs
func (db DB) GetPosts() []Post {
	query, err := db.Query("SELECT * FROM posts")
	if err != nil {
		log.Panic(err)
	}
	defer query.Close()

	var posts []Post
	for query.Next() {
		var post Post
		err = query.Scan(&post.Id, &post.Title, &post.Content, &post.AuthorId, &post.Date, &post.LikesCount, &post.DislikesCount, &post.CommentsCount)
		if err != nil {
			log.Panic(err)
		}
		post.Author = db.GetUserById(post.AuthorId)
		posts = append(posts, post)
	}
	return posts
}

func (db DB) GetPostById(id int) *Post {
	query, err := db.Query("SELECT * FROM posts WHERE id = ?", id)
	if err != nil {
		log.Panic(err)
	}
	defer query.Close()

	var post Post
	if !query.Next() {
		return nil
	}
	err = query.Scan(&post.Id, &post.Title, &post.Content, &post.AuthorId, &post.Date, &post.LikesCount, &post.DislikesCount, &post.CommentsCount)
	if err != nil {
		log.Panic(err)
	}
	post.Author = db.GetUserById(post.AuthorId)
	return &post
}

// AddPost adds post to database, returns id of new post
func (db DB) AddPost(title, content string, authorId int) int {
	result, err := db.Exec(`INSERT INTO posts (title, content, author, date, likes_count, dislikes_count, comments_count) 
								  VALUES (?, ?, ?, ?, ?, ?, ?)`, title, content, authorId, time.Now().Format(time.RFC3339), 0, 0, 0)
	if err != nil {
		log.Panic(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Panic(err)
	}
	return int(id)
}
