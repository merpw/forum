package database

import (
	"log"
	"time"
)

func (db DB) GetAllFollowersById(id int) (followerIds []int) {
	query, err := db.Query("SELECT follower_id FROM followers WHERE user_id = ? ORDER BY timestamp", id)
	if err != nil {
		log.Panic(err)
	}

	for query.Next() {
		var followerId int
		err = query.Scan(&followerId)
		if err != nil {
			log.Panic(err)
		}
		followerIds = append(followerIds, followerId)
	}

	query.Close()
	return
}

func (db DB) GetFollowStatus(meId, userId int) *FollowStatus {
	if meId == userId {
		return nil
	}

	query, err := db.Query("SELECT id FROM followers WHERE follower_id = ? AND user_id = ? COLLATE NOCASE",
		meId, userId)
	if err != nil {
		log.Panic(err)
	}

	var followStatus = new(FollowStatus)

	if query.Next() {
		if err != nil {
			log.Panic(err)
		}
		*followStatus = Following
		query.Close()
		return followStatus
	}

	query.Close()

	if db.CheckIfInvitationExists(meId, userId) {
		*followStatus = RequestToFollow
		return followStatus
	}

	*followStatus = NotFollowing

	return followStatus
}

func (db DB) Unfollow(followerId, userId int) FollowStatus {
	_, err := db.Exec("DELETE FROM followers WHERE user_id = ? AND follower_id = ?", userId, followerId)
	if err != nil {
		log.Panic(err)
	}
	return NotFollowing
}

func (db DB) Follow(followerId, userId int) FollowStatus {
	_, err := db.Exec(`INSERT INTO followers 
    	(user_id, follower_id, timestamp)
		VALUES (?, ?, ?)`,
		userId, followerId, time.Now().Format(time.RFC3339))
	if err != nil {
		log.Panic(err)
	}

	return Following
}
