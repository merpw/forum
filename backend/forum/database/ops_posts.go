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
	query, err := db.Query(`
		SELECT DISTINCT posts.* FROM posts 
		LEFT JOIN post_audience ON posts.id = post_audience.post_id 
		LEFT JOIN followers ON post_audience.follow_id = followers.id AND followers.follower_id = ?   
		LEFT JOIN followers AS f2 ON posts.author = f2.user_id AND f2.follower_id = ?
		LEFT JOIN group_members ON group_members.group_id = posts.group_id AND group_members.user_id = ?
		WHERE 
			posts.privacy = 0 AND posts.group_id IS NULL OR
			posts.author = ? OR
			(posts.privacy = 1 AND f2.id IS NOT NULL) OR
			(posts.privacy = 2 AND post_audience.id IS NOT NULL) OR
			(posts.group_id IS NOT NULL AND group_members.id IS NOT NULL)
	`, userId, userId, userId, userId)
	if err != nil {
		log.Panic(err)
	}

	var posts []Post
	for query.Next() {
		var post Post
		err = query.Scan(&post.Id, &post.Title, &post.Content, &post.AuthorId, &post.Date,
			&post.LikesCount, &post.DislikesCount, &post.CommentsCount,
			&post.Categories, &post.Description, &post.GroupId, &post.Privacy)
		if err != nil {
			log.Panic(err)
		}
		posts = append(posts, post)
	}
	query.Close()

	return posts
}

// GetPostById reads post from database by post_id
func (db DB) GetPostById(postId, userId int) *Post {
	row := db.QueryRow(`
    SELECT posts.* FROM posts
		LEFT JOIN post_audience ON post_audience.post_id = posts.id 
		LEFT JOIN followers ON post_audience.follow_id = followers.id AND followers.follower_id = ?   
		LEFT JOIN followers AS f2 ON posts.author = f2.user_id AND f2.follower_id = ?
        LEFT JOIN group_members ON group_members.group_id = posts.group_id AND group_members.user_id = ?
    WHERE posts.id = ? AND (
        posts.privacy = 0 AND posts.group_id IS NULL OR
        posts.author = ? OR
        (posts.privacy = 1 AND f2.id IS NOT NULL) OR
        (posts.privacy = 2 AND post_audience.id IS NOT NULL) OR
        (posts.group_id IS NOT NULL AND group_members.id IS NOT NULL)
    )
`, userId, userId, userId, postId, userId)

	var post Post
	err := row.Scan(&post.Id, &post.Title, &post.Content, &post.AuthorId, &post.Date,
		&post.LikesCount, &post.DislikesCount, &post.CommentsCount,
		&post.Categories, &post.Description, &post.GroupId, &post.Privacy)

	if err != nil {
		return nil
	}

	return &post
}

// AddPost adds post to database, returns id of new post
func (db DB) AddPost(title, content, description string,
	authorId int, categories string, privacy Privacy, groupId sql.NullInt64) int {
	result, err := db.Exec(`INSERT INTO posts 
    	(title, content, author, date, likes_count, dislikes_count, comments_count,
    	 categories, description, group_id, privacy)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		title, content, authorId, time.Now().Format(time.RFC3339), 0, 0, 0, categories, description, groupId, privacy)
	if err != nil {
		log.Panic(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Panic(err)
	}
	return int(id)
}

func (db DB) AddPostAudience(postId, followId int) {
	_, err := db.Exec(`
					INSERT INTO post_audience (post_id, follow_id)
					VALUES (?, ?)`, postId, followId)
	if err != nil {
		log.Panic(err)
	}
}

func (db DB) RemovePostAudience(postId, followId int) {
	_, err := db.Exec("DELETE FROM post_audience WHERE post_id = ? AND follow_id = ?", postId, followId)
	if err != nil {
		log.Panic(err)
	}
}

func (db DB) GetPostAudienceId(postId, followId int) (id *int) {
	err := db.QueryRow("SELECT id FROM post_audience WHERE post_id = ? AND followId = ?").Scan(id)
	if err != nil {
		return nil
	}
	return
}

func (db DB) GetPostFollowStatus(postId, followId int) bool {
	err := db.QueryRow(`
		SELECT id FROM post_audience 
		    WHERE post_id = ? AND follow_id = ?
		`, postId, followId).Err()
	return err != nil
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
			&post.Categories, &post.Description, &post.GroupId, &post.Privacy)
		if err != nil {
			log.Panic(err)
		}
		posts = append(posts, post)
	}
	query.Close()

	return posts
}

func (db DB) GetPostsByUserId(userId, followerId int) []Post {
	query, err := db.Query(`
		SELECT DISTINCT posts.* FROM posts
		LEFT JOIN post_audience ON posts.id = post_audience.post_id AND post_audience.follow_id = ?
		LEFT JOIN followers ON posts.author = followers.user_id AND followers.follower_id = ?
		WHERE posts.author = ? 
		AND (posts.privacy = 0 OR	
			(posts.privacy = 1 AND followers.id IS NOT NULL) OR
			(posts.privacy = 2 AND post_audience.post_id IS NOT NULL)
		)`, followerId, followerId, userId)
	if err != nil {
		log.Panic(err)
	}

	var posts []Post
	for query.Next() {
		var post Post
		err = query.Scan(&post.Id, &post.Title, &post.Content, &post.AuthorId, &post.Date,
			&post.LikesCount, &post.DislikesCount, &post.CommentsCount,
			&post.Categories, &post.Description, &post.GroupId, &post.Privacy)
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
	SELECT * FROM posts WHERE id IN (SELECT post_id FROM post_reactions WHERE reaction = 1 AND author_id = ?) 
		AND (
		    privacy = 0 OR (
		        author IN (SELECT user_id FROM followers WHERE follower_id = ?)
	    )
	)`, userId, userId)

	if err != nil {
		log.Panic(err)
	}

	var posts []Post
	for query.Next() {
		var post Post
		err = query.Scan(&post.Id, &post.Title, &post.Content, &post.AuthorId, &post.Date,
			&post.LikesCount, &post.DislikesCount, &post.CommentsCount,
			&post.Categories, &post.Description, &post.GroupId, &post.Privacy)
		if err != nil {
			log.Panic(err)
		}
		posts = append(posts, post)
	}
	query.Close()

	return posts
}

func (db DB) GetCategoryPosts(category string, userId int) []Post {
	query, err := db.Query(`
		SELECT DISTINCT posts.* FROM posts 
		LEFT JOIN post_audience ON posts.id = post_audience.post_id 
		LEFT JOIN followers ON post_audience.follow_id = followers.id AND followers.follower_id = ?   
		LEFT JOIN followers AS f2 ON posts.author = f2.user_id AND f2.follower_id = ?
		WHERE 
			categories LIKE '%' || ? || '%' AND 
			(
				posts.author = ? OR
				posts.privacy = 0 OR
				(posts.privacy = 1 AND followers.id IS NOT NULL) OR
				(posts.privacy = 2 AND post_audience.post_id IS NOT NULL)
			) 
`, userId, userId, category, userId)
	if err != nil {
		log.Panic(err)
	}

	var posts []Post
	for query.Next() {
		var post Post
		err = query.Scan(&post.Id, &post.Title, &post.Content, &post.AuthorId, &post.Date,
			&post.LikesCount, &post.DislikesCount, &post.CommentsCount,
			&post.Categories, &post.Description, &post.GroupId, &post.Privacy)
		if err != nil {
			log.Panic(err)
		}
		posts = append(posts, post)
	}
	query.Close()

	return posts
}

// GetCommentById gets a comment by its id
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
