package database

import (
	"database/sql"
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
		err = query.Scan(&post.Id, &post.Title, &post.Content, &post.AuthorId, &post.Date,
			&post.LikesCount, &post.DislikesCount, &post.CommentsCount,
			&post.Categories, &post.Description, &post.GroupId)
		if err != nil {
			log.Panic(err)
		}
		posts = append(posts, post)
	}
	query.Close()

	return posts
}

// GetPostById reads post from database by post_id, does not require user to be logged in
func (db DB) GetPostById(id int) *Post {
	query, err := db.Query("SELECT * FROM posts WHERE id = ?", id)
	if err != nil {
		log.Panic(err)
	}

	var post Post
	if !query.Next() {
		return nil
	}
	err = query.Scan(&post.Id, &post.Title, &post.Content, &post.AuthorId, &post.Date,
		&post.LikesCount, &post.DislikesCount, &post.CommentsCount, &post.Categories, &post.Description, &post.GroupId)
	if err != nil {
		log.Panic(err)
	}
	query.Close()

	return &post
}

// AddPost adds post to database, returns id of new post
func (db DB) AddPost(title, content, description string, authorId int, categories string, groupId sql.NullInt64) int {
	result, err := db.Exec(`INSERT INTO posts 
    	(title, content, author, date, likes_count, dislikes_count, comments_count, categories, description, group_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		title, content, authorId, time.Now().Format(time.RFC3339), 0, 0, 0, categories, description, groupId)
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
		err = query.Scan(&post.Id, &post.Title, &post.Content, &post.AuthorId, &post.Date,
			&post.LikesCount, &post.DislikesCount, &post.CommentsCount,
			&post.Categories, &post.Description, &post.GroupId)
		if err != nil {
			log.Panic(err)
		}
		posts = append(posts, post)
	}
	query.Close()

	return posts
}

func (db DB) GetUserPostsLiked(userId int) []Post {
	query, err := db.Query(`SELECT * FROM posts WHERE id IN 
	(SELECT post_id FROM post_reactions WHERE author_id = ? AND reaction = 1)`, userId)

	if err != nil {
		log.Panic(err)
	}

	var posts []Post
	for query.Next() {
		var post Post
		err = query.Scan(&post.Id, &post.Title, &post.Content, &post.AuthorId, &post.Date,
			&post.LikesCount, &post.DislikesCount, &post.CommentsCount,
			&post.Categories, &post.Description, &post.GroupId)
		if err != nil {
			log.Panic(err)
		}
		posts = append(posts, post)
	}
	query.Close()

	return posts
}

func (db DB) GetCategoryPosts(category string) []Post {
	query, err := db.Query("SELECT * FROM posts WHERE categories LIKE '%' || ? || '%'", category)
	if err != nil {
		log.Panic(err)
	}

	var posts []Post
	for query.Next() {
		var post Post
		err = query.Scan(&post.Id, &post.Title, &post.Content, &post.AuthorId, &post.Date,
			&post.LikesCount, &post.DislikesCount, &post.CommentsCount,
			&post.Categories, &post.Description, &post.GroupId)
		if err != nil {
			log.Panic(err)
		}
		posts = append(posts, post)
	}
	query.Close()

	return posts
}

// GetCommentById gets Comment Struct Pointer by comment_id from the database, does not require user to be logged in
func (db DB) GetCommentById(id int) *Comment {
	query, err := db.Query("SELECT * FROM comments WHERE id = ?", id)
	if err != nil {
		log.Panic(err)
	}

	if !query.Next() {
		return nil
	}
	var comment Comment
	err = query.Scan(&comment.Id, &comment.PostId, &comment.AuthorId, &comment.Content, &comment.Date,
		&comment.LikesCount, &comment.DislikesCount)
	if err != nil {
		log.Panic(err)
	}
	query.Close()

	return &comment
}

// GetPostComments gets all comments for post using post_id
//
// Example:
//
//	comments := db.GetPostComments(1)
func (db DB) GetPostComments(postId int) []Comment {
	query, err := db.Query("SELECT * FROM comments WHERE post_id = ?", postId)
	if err != nil {
		log.Panic(err)
	}

	var comments []Comment
	for query.Next() {
		var comment Comment
		err = query.Scan(&comment.Id, &comment.PostId, &comment.AuthorId, &comment.Content, &comment.Date,
			&comment.LikesCount, &comment.DislikesCount)
		if err != nil {
			log.Panic(err)
		}
		comments = append(comments, comment)
	}
	query.Close()

	return comments
}
