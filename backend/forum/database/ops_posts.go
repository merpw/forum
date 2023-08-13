package database

import (
	"database/sql"
	"log"
	"time"
)

// GetAllPosts reads all posts from database (reads only userId, not user object)
//
// panics if error occurs
func (db DB) GetAllPosts(userId int) []Post {
	queryStr := `
		SELECT DISTINCT posts.* 
		FROM posts 
		-- Join to identify super-private posts that have the requesting user in their audience
		LEFT JOIN post_audience ON posts.id = post_audience.post_id 
		-- Join to determine authors the requesting user follows
		LEFT JOIN followers AS f2 ON posts.author = f2.user_id AND f2.follower_id = ?
		WHERE 
			-- Posts that are public
			posts.privacy = 0 OR
		
			-- Posts that are private and the requesting user is the author
			posts.author = ? OR
			
			-- Private posts where the requesting user is following the author
			(posts.privacy = 1 AND f2.id IS NOT NULL) OR
			
			-- Super-private posts where the requesting user is in the audience
			(posts.privacy = 2 AND post_audience.id IS NOT NULL)
	`

	query, err := db.Query(queryStr, userId, userId)
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
func (db DB) AddPost(title, content, description string,
	authorId int, categories string, privacy Privacy, groupId sql.NullInt64) int {
	result, err := db.Exec(`INSERT INTO posts 
    	(title, content, author, date, likes_count,
    	 dislikes_count, comments_count, categories, description, privacy, group_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		title, content, authorId, time.Now().Format(time.RFC3339), 0, 0, 0, categories, description, privacy, groupId)
	if err != nil {
		log.Panic(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Panic(err)
	}
	return int(id)
}

func (db DB) GetPostFollowStatus(postId, followId int) bool {
	err := db.QueryRow(`
				SELECT id 
				FROM post_audience WHERE post_id = ? AND follow_id = ? `,
		postId, followId).Err()

	return err != nil
}

func (db DB) AddPostAudience(postId, followerId int) {
	_, err := db.Exec(`
					INSERT INTO post_audience (post_id, follow_id)
					VALUES (?, ?)`, postId, followerId)
	if err != nil {
		log.Panic(err)
	}
}

func (db DB) GetMePosts(userId int) []Post {
	query, err := db.Query("SELECT * FROM posts WHERE author = ?", userId)
	if err != nil {
		log.Panic(err)
	}

	var posts []Post
	for query.Next() {
		var post Post
		err = query.Scan(&post.Id, &post.Title, &post.Content, &post.AuthorId, &post.Date,
			&post.LikesCount, &post.DislikesCount, &post.CommentsCount,
			&post.Categories, &post.Description)
		if err != nil {
			log.Panic(err)
		}
		posts = append(posts, post)
	}
	query.Close()

	return posts
}

func (db DB) GetUserPosts(userId, followerId int) []Post {
	queryStr := `
		SELECT DISTINCT posts.*
		FROM posts
		-- Left Join with post_audience to capture super-private posts for the requesting user
		LEFT JOIN post_audience ON posts.Id = post_audience.post_id AND post_audience.follow_id = ?
		-- Left Join with followers table to check if the requesting user is a follower of the post's author
		LEFT JOIN followers ON posts.author = followers.user_id AND followers.follower_id = ?
		WHERE 
			posts.author = ? AND 
			(
				-- Posts that are public
				posts.privacy = 0 OR
				-- Posts that are private and the requesting user is a follower
				(posts.privacy = 1 AND followers.id IS NOT NULL) OR
				-- Posts that are super-private and the requesting user is in the audience
				(posts.privacy = 2 AND post_audience.post_id IS NOT NULL)
	)
`
	query, err := db.Query(queryStr, followerId, followerId, userId)
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
	query, err := db.Query(`
		SELECT * FROM posts
		WHERE id IN (
			SELECT post_id
			FROM post_reactions
			WHERE author_id = ? AND reaction = 1
		)
		AND (
			privacy = 0
			OR (
				privacy = 1
				AND author IN (
					SELECT user_id
					FROM followers
					WHERE follower_id = ?
				)
			)
		)
	`, userId, userId)

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
