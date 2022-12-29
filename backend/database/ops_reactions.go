package database

import "log"

// AddPostReaction adds specified reaction to post
//
// reaction=1 is like and reaction=-1 is dislike
func (db DB) AddPostReaction(postId, userId, reaction int) {
	_, err := db.Exec("INSERT INTO post_reactions (post_id, author_id, reaction) VALUES (?, ?, ?)", postId, userId, reaction)
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
