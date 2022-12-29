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

// LikePost from database, returns number of likes of the post
func (db DB) LikePost(postId, userId int) int {

	// check if user already liked/disliked the post
	query, err := db.Query(`SELECT * FROM likes WHERE parent_id = ? AND author_id = ? AND post_like = 1`, postId, userId)
	if err != nil {
		log.Panic(err)
	}
	var like Like
	var value int
	if query.Next() {
		err := query.Scan(&like.Id, &like.ParentId, &like.AuthorId, &like.PostLike, &like.Value)
		if err != nil {
			log.Panic(err)
		}
		value = like.Value
	}
	switch value {
	case 1:
		// TODO like-1 (was liked before)
		statement, err := db.Prepare(`UPDATE posts SET likes_count = likes_count + ? WHERE id = ?`)
		if err != nil {
			log.Panic(err)
		}
		statement.Exec(-1, like.Id)
		return like.Value - 1
	case -1:
		// TODO dislike-1, then like+1 (because was disliked before)
		return like.Value + 1
	default:
		// TODO make new like record into database
	}

	return 1
}
