package database

import (
	"log"
)

// AddPostReaction adds specified reaction to post
//
// reaction=1 is like and reaction=-1 is dislike
func (db DB) AddPostReaction(postId, userId, reaction int) {
	_, err := db.Exec(`INSERT INTO post_reactions 
    	(post_id, author_id, reaction) VALUES (?, ?, ?)`,
		postId, userId, reaction)
	if err != nil {
		log.Panic(err)
	}
}

// UpdatePostLikesCount changes post's likes_count value
func (db DB) UpdatePostLikesCount(postId int, change int) {
	_, err := db.Exec("UPDATE posts SET likes_count = likes_count + ? WHERE id = ?", change, postId)
	if err != nil {
		log.Panic(err)
	}
}

// UpdatePostDislikeCount changes post's dislikes_count value
func (db DB) UpdatePostDislikeCount(postId int, change int) {
	_, err := db.Exec("UPDATE posts SET dislikes_count = posts.dislikes_count + ? WHERE id = ?", change, postId)
	if err != nil {
		log.Panic(err)
	}
}

// RemovePostReaction removes user's post reaction
func (db DB) RemovePostReaction(postId, userId int) {
	_, err := db.Exec("DELETE FROM post_reactions WHERE post_id = ? AND author_id = ?", postId, userId)
	if err != nil {
		log.Panic(err)
	}
}

// GetPostReaction returns user's reaction to post
//
// returns 1 if user liked post, -1 if disliked and 0 if not reacted
func (db DB) GetPostReaction(postId, userId int) int {
	query, err := db.Query("SELECT reaction FROM post_reactions WHERE post_id = ? AND author_id = ?", postId, userId)
	if err != nil {
		log.Panic(err)
	}
	defer query.Close()

	if !query.Next() {
		return 0
	}
	var reaction int
	err = query.Scan(&reaction)
	if err != nil {
		log.Panic(err)
	}
	return reaction
}

// GetCommentReaction returns user's reaction to comment
//
// returns 1 if user liked comment, -1 if disliked and 0 if not reacted
func (db DB) GetCommentReaction(commentId, userId int) int {
	query, err := db.Query(`SELECT reaction FROM comment_reactions 
                WHERE comment_id = ? AND author_id = ?`, commentId, userId)
	if err != nil {
		log.Panic(err)
	}
	defer query.Close()

	if !query.Next() {
		return 0
	}
	var reaction int
	err = query.Scan(&reaction)
	if err != nil {
		log.Panic(err)
	}
	return reaction
}

// AddCommentReaction adds specified reaction to comment
//
// reaction=1 is like and reaction=-1 is dislike
func (db DB) AddCommentReaction(commentId, userId, reaction int) {
	_, err := db.Exec(`INSERT INTO comment_reactions 
    	(comment_id, author_id, reaction) VALUES (?, ?, ?)`,
		commentId, userId, reaction)
	if err != nil {
		log.Panic(err)
	}
}

// UpdateCommentLikesCount changes comment's likes_count value
func (db DB) UpdateCommentLikesCount(commentId int, change int) {
	_, err := db.Exec("UPDATE comments SET likes_count = likes_count + ? WHERE id = ?", change, commentId)
	if err != nil {
		log.Panic(err)
	}
}

// UpdateCommentDislikeCount changes comment's dislikes_count value
func (db DB) UpdateCommentDislikeCount(commentId int, change int) {
	_, err := db.Exec("UPDATE comments SET dislikes_count = comments.dislikes_count + ? WHERE id = ?", change, commentId)
	if err != nil {
		log.Panic(err)
	}
}

// RemoveCommentReaction removes user's comment reaction
func (db DB) RemoveCommentReaction(commentId, userId int) {
	_, err := db.Exec("DELETE FROM comment_reactions WHERE comment_id = ? AND author_id = ?", commentId, userId)
	if err != nil {
		log.Panic(err)
	}
}
