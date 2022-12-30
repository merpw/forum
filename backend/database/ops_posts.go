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
		err = query.Scan(&post.Id, &post.Title, &post.Content, &post.AuthorId, &post.Date, &post.LikesCount, &post.DislikesCount, &post.CommentsCount, &post.Category)
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
	err = query.Scan(&post.Id, &post.Title, &post.Content, &post.AuthorId, &post.Date, &post.LikesCount, &post.DislikesCount, &post.CommentsCount, &post.Category)
	if err != nil {
		log.Panic(err)
	}
	query.Close()

	return &post
}

// AddPost adds post to database, returns id of new post
func (db DB) AddPost(title, content string, authorId int, category string) int {
	result, err := db.Exec(`INSERT INTO posts (title, content, author, date, likes_count, dislikes_count, comments_count, category) 
								  VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, title, content, authorId, time.Now().Format(time.RFC3339), 0, 0, 0, category)
	if err != nil {
		log.Panic(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Panic(err)
	}
	return int(id)
}

func (db DB) GetUserPosts(userId int) []Post {
	query, err := db.Query("SELECT * FROM posts WHERE author = ?", userId)
	if err != nil {
		log.Panic(err)
	}

	var posts []Post
	for query.Next() {
		var post Post
		err = query.Scan(&post.Id, &post.Title, &post.Content, &post.AuthorId, &post.Date, &post.LikesCount, &post.DislikesCount, &post.CommentsCount, &post.Category)
		if err != nil {
			log.Panic(err)
		}
		posts = append(posts, post)
	}
	query.Close()

	return posts
}

func (db DB) GetCategoryPosts(category string) []Post {
	query, err := db.Query("SELECT * FROM posts WHERE category = ?", category)
	if err != nil {
		log.Panic(err)
	}

	var posts []Post
	for query.Next() {
		var post Post
		err = query.Scan(&post.Id, &post.Title, &post.Content, &post.AuthorId, &post.Date, &post.LikesCount, &post.DislikesCount, &post.CommentsCount, &post.Category)
		if err != nil {
			log.Panic(err)
		}
		posts = append(posts, post)
	}
	query.Close()

	return posts
}

// comments
func (db DB) GetCommentById(id int) *Comment {
	query, err := db.Query("SELECT * FROM comments WHERE id = ?", id)
	if err != nil {
		log.Panic(err)
	}

	if !query.Next() {
		return nil
	}
	var comment Comment
	err = query.Scan(&comment.Id, &comment.PostId, &comment.AuthorId, &comment.Content, &comment.Date, &comment.LikesCount, &comment.DislikesCount)
	if err != nil {
		log.Panic(err)
	}
	query.Close()

	return &comment
}

func (db DB) GetPostComments(postId int) []Comment {
	query, err := db.Query("SELECT * FROM comments WHERE post_id = ?", postId)
	if err != nil {
		log.Panic(err)
	}

	var comments []Comment
	for query.Next() {
		var comment Comment
		err = query.Scan(&comment.Id, &comment.PostId, &comment.AuthorId, &comment.Content, &comment.Date, &comment.LikesCount, &comment.DislikesCount)
		if err != nil {
			log.Panic(err)
		}
		comments = append(comments, comment)
	}
	query.Close()

	return comments
}
