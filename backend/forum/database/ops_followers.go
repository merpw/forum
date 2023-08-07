package database

import (
	"log"
	"time"
)

func (db DB) GetFollowStatus(meId, userId int) *FollowStatus {
	if meId == userId {
		return nil
	}

	query, err := db.Query("SELECT follower_id FROM followers WHERE user_id = ? COLLATE NOCASE", userId)
	if err != nil {
		log.Panic(err)
	}

	var followerId = -2

	if query.Next() {
		err = query.Scan(&followerId)
		if err != nil {
			log.Panic(err)
		}
	}

	query.Close()

	var followStatus = new(FollowStatus)
	if followerId == meId {
		*followStatus = Following
		return followStatus
	}

	if db.CheckIfInvitationExists(meId, userId) {
		*followStatus = RequestToFollow
	} else {
		*followStatus = NotFollowing
	}

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
