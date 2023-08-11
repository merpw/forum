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
	return groupId
}

func (db DB) GetTopGroups() (groupIds []int) {
	query, err := db.Query(`
SELECT id FROM groups 
    INNER JOIN (
    	SELECT group_id, COUNT(*) AS member_count FROM group_members
		GROUP BY group_id
	) ON group_id = id                
ORDER BY member_count DESC
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

	return groupIds
}

func (db DB) GetGroupMemberCount(groupId int) (memberCount int) {
	err := db.QueryRow(`
	SELECT COUNT(*) FROM group_members WHERE group_id = ?
	`, groupId).Scan(&memberCount)
	if err != nil {
		log.Panic(err)
	}
	return memberCount
}

func (db DB) GetGroupById(groupId int) *Group {
	group := &Group{
		Id: groupId,
	}
	err := db.QueryRow(`SELECT title, description FROM groups WHERE id = ?`,
		groupId).Scan(&group.Title, &group.Description)
	if err != nil {
		return nil
	}

	return group
}

func (db DB) GetGroupMemberStatus(groupId, userId int) (memberStatus *InviteStatus) {
	err := db.QueryRow(`SELECT CASE
    WHEN (SELECT 1 FROM group_members WHERE group_id = ? AND user_id = ?) THEN 1
    ELSE (
        SELECT CASE
            WHEN (SELECT 1 FROM invitations WHERE type = 1 AND from_user_id = ? AND to_user_id = ?) THEN 2
            WHEN (SELECT 1 FROM invitations WHERE type = 2 AND from_user_id = ? AND to_user_id = ?) THEN 2
            ELSE 0
        END
    )
		END AS member_status
	`, groupId, userId, groupId, userId, userId, groupId).Scan(&memberStatus)
	if err != nil {
		log.Panic(err)
	}

	return memberStatus
}

func (db DB) GetGroupCreatorId(groupId int) *int {
	var creatorId int
	err := db.QueryRow("SELECT creator_id FROM groups WHERE id = ?", groupId).Scan(&creatorId)
	if err != nil {
		return nil
	}

	return &creatorId
}

func (db DB) GetGroupMembers(groupId int, withPending bool) (members []int) {

	var q string
	if withPending {
		q = `
		SELECT user_id FROM group_members WHERE group_id = :groupId
		                                  
		UNION -- users invited to join -- 
        SELECT to_user_id FROM invitations WHERE type = 1 AND associated_id = :groupId
                                           
		UNION -- users requested to join -- 
		SELECT from_user_id FROM invitations WHERE type = 2 AND associated_id = :groupId
        `
	} else {
		q = "SELECT user_id FROM group_members WHERE group_id = ?"
	}

	query, err := db.Query(q, groupId)
	if err != nil {
		log.Panic(err)
	}

	members = make([]int, 0)
	for query.Next() {
		var memberId int
		err = query.Scan(&memberId)
		if err != nil {
			log.Panic(err)
		}
		members = append(members, memberId)
	}
	query.Close()

	return members
}

func (db DB) GetGroupPostsById(groupId int) (postIds []int) {
	query, err := db.Query("SELECT id FROM posts WHERE group_id = ? ORDER BY id DESC", groupId)
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
	return postIds
}

func (db DB) DeleteGroupMembership(groupId, userId int) {
	_, err := db.Exec("DELETE FROM group_members WHERE group_id = ? AND user_id = ?", groupId, userId)
	if err != nil {
		log.Panic(err)
	}
}

func (db DB) AddMembership(groupId, userId int) {
	_, err := db.Exec(`INSERT INTO group_members 
    	(group_id, user_id, timestamp)
		VALUES (?, ?, ?)`,
		groupId, userId, time.Now().Format(time.RFC3339))
	if err != nil {
		log.Panic(err)
	}
}
