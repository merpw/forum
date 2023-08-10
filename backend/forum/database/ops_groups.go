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
		var test any
		err = query.Scan(&test, &groupId)
		if err != nil {
			log.Panic(err)
		}
		groupIds = append(groupIds, groupId)
	}
	query.Close()

	return
}

func (db DB) GetGroupMembersByGroupId(groupId int) (members int) {
	err := db.QueryRow(`
						SELECT COUNT(*) 
						FROM group_members 
						WHERE group_id = ?`, groupId).Scan(&members)
	if err != nil {
		log.Panic(err)
	}
	return
}

func (db DB) GetGroupById(groupId int) *Group {
	var t, d string
	err := db.QueryRow(`SELECT title, description FROM groups WHERE id = ?`, groupId).Scan(&t, &d)
	if err != nil {
		return nil
	}

	return &Group{
		Id:          groupId,
		Title:       t,
		Description: d,
	}
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

	return
}

func (db DB) GetGroupCreatorId(groupId int) *int {
	var creatorId int
	err := db.QueryRow("SELECT creator_id FROM groups WHERE id = ?", groupId).Scan(&creatorId)
	if err != nil {
		return nil
	}

	return &creatorId

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
	return
}

func (db DB) DeleteGroupMembership(groupId, userId int) InviteStatus {
	_, err := db.Exec("DELETE FROM group_members WHERE group_id = ? AND user_id = ?", groupId, userId)
	if err != nil {
		log.Panic(err)
	}
	return Inactive
}

func (db DB) AddMembership(groupId, userId int) InviteStatus {
	_, err := db.Exec(`INSERT INTO group_members 
    	(group_id, user_id, timestamp)
		VALUES (?, ?, ?)`,
		groupId, userId, time.Now().Format(time.RFC3339))
	if err != nil {
		log.Panic(err)
	}
	return Accepted
}
