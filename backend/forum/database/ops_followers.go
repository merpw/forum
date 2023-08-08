package database

import (
	"log"
	"time"
)

func (db DB) GetUserFollowers(userId int) (followerIds []int) {
	query, err := db.Query("SELECT follower_id FROM followers WHERE user_id = ? ORDER BY timestamp", userId)
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

func (db DB) GetFollowStatus(followerId, userId int) *FollowStatus {
	row := db.QueryRow(`
    SELECT CASE 
    WHEN (
       SELECT 1 FROM followers WHERE user_id = ? AND follower_id = ?) THEN 1
    ELSE (
        SELECT CASE 
        WHEN (
        	SELECT 1 FROM invitations WHERE from_user_id = ? AND to_user_id = ?) THEN 2
        ELSE 0
    	END
    )
    END 
    AS follow_status
    `, userId, followerId, followerId, userId)

	var followStatus = new(FollowStatus)
	err := row.Scan(followStatus)
	if err != nil {
		log.Panic(err)
	}
	return followStatus
}

func (db DB) RemoveFollower(followerId, userId int) FollowStatus {
	_, err := db.Exec("DELETE FROM followers WHERE user_id = ? AND follower_id = ?", userId, followerId)
	if err != nil {
		log.Panic(err)
	}
	return NotFollowing
}

func (db DB) AddFollower(followerId, userId int) FollowStatus {
	_, err := db.Exec(`INSERT INTO followers 
    	(user_id, follower_id, timestamp)
		VALUES (?, ?, ?)`,
		userId, followerId, time.Now().Format(time.RFC3339))
	if err != nil {
		log.Panic(err)
	}

	return Following
}
