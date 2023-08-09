package database

import (
	"log"
	"time"
)

func (db DB) AddGroup(creatorId int, title, description string) (groupId int64) {
	result, err := db.Exec(`INSERT INTO groups 
    	(creator_id, title, description, timestamp)
		VALUES (?, ?, ?, ?)`,
		creatorId, title, description, time.Now().Format(time.RFC3339))
	if err != nil {
		log.Panic(err)
	}
	groupId, err = result.LastInsertId()
	if err != nil {
		log.Panic(err)
	}
	return
}

func (db DB) GetGroupIdsByMembers() (groupIds []int) {
	query, err := db.Query(`
	SELECT id, COUNT(id) AS occurrence
	FROM group_members 
	GROUP BY group_id
	ORDER BY occurrence DESC;
	`)
	if err != nil {
		log.Panic(err)
	}

	groupIds = make([]int, 0)
	for query.Next() {
		var groupId int
		err = query.Scan(&groupId)
		if err != nil {
			log.Panic(err)
		}
		groupIds = append(groupIds, groupId)
	}
	query.Close()

	return
}

func (db DB) GetGroupById(groupId int) *Group {
	row := db.QueryRow("SELECT * FROM groups WHERE id = ?", groupId)
	var group Group
	err := row.Scan(&group)
	if err != nil {
		return nil
	}

	return &group
}

func (db DB) GetGroupMemberStatus(groupId, userId int) *MemberStatus {
	row := db.QueryRow(`
	SELECT CASE
	WHEN (
	    SELECT 1 FROM group_members WHERE group_id ? AND user_id = ?) THEN 1
	ELSE (
		SELECT CASE
		WHEN (
			SELECT 1 FROM invitations WHERE from_user_id = ? AND to_user_id = ?) THEN 2
		ELSE 0
		END
		)
	END
	AS member_status
	`, groupId, userId, userId, groupId)

	var memberStatus = new(MemberStatus)
	err := row.Scan(&memberStatus)
	if err != nil {
		log.Panic(err)
	}

	return memberStatus
}

func (db DB) GetGroupPostsById(groupId int) (postIds []int) {
	query, err := db.Query("SELECT id FROM posts WHERE group_id = ? ORDER BY timestamp", groupId)
	if err != nil {
		log.Panic(err)
	}

	postIds = make([]int, 0)
	for query.Next() {
		var postId int
		err = query.Scan(&postId)
		if err != nil {
			log.Panic(err)
		}
		postIds = append(postIds, postId)
	}

	query.Close()
	return
}

func (db DB) DeleteGroupMembership(groupId, userId int) MemberStatus {
	_, err := db.Exec("DELETE FROM group_members WHERE group_id = ? AND userId = ?", groupId, userId)
	if err != nil {
		log.Panic(err)
	}
	return NotMember
}
