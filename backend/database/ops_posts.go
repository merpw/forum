package database

import (
	"log"
	"time"
)

// GetAllPosts reads all posts from database (reads only userId, not user object)
//
// panics if error occurs
func (db DB) GetAllPosts() []Post {
	query, err := db.Query("SELECT * FROM posts")
	if err != nil {
		log.Panic(err)
	}

	var posts []Post
	for query.Next() {
		var post Post
		err = query.Scan(&post.Id, &post.Title, &post.Content, &post.AuthorId, &post.Date, &post.LikesCount, &post.DislikesCount, &post.CommentsCount)
		if err != nil {
			log.Panic(err)
		}
		posts = append(posts, post)
	}
	query.Close()

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
	err = query.Scan(&post.Id, &post.Title, &post.Content, &post.AuthorId, &post.Date, &post.LikesCount, &post.DislikesCount, &post.CommentsCount)
	if err != nil {
		log.Panic(err)
	}
	query.Close()

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

//func (db DB) GetUserPosts(userId int) []Post {
//	query, err := db.Query("SELECT * FROM posts WHERE author = ?", userId)
//	if err != nil {
//		log.Panic(err)
//	}
//
//	var posts []Post
//	for query.Next() {
//		var post Post
//		err = query.Scan(&post.Id, &post.Title, &post.Content, &post.AuthorId, &post.Date, &post.LikesCount, &post.DislikesCount, &post.CommentsCount)
//		if err != nil {
//			log.Panic(err)
//		}
//		posts = append(posts, post)
//	}
//	query.Close()
//
//	return posts
//}
