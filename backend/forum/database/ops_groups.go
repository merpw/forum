package database

import "log"

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
