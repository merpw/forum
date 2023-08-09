package database

import (
	"database/sql"
	"errors"
	"log"
	"time"
)

// GetUserInvitations returns slice of all users ids
func (db DB) GetUserInvitations(toUserId int) (invitationIds []int) {
	query, err := db.Query("SELECT id FROM invitations WHERE to_user_id = ? ORDER BY timestamp DESC", toUserId)
	if err != nil {
		log.Panic(err)
	}

	invitationIds = make([]int, 0)

	for query.Next() {
		var invitationId int
		err = query.Scan(&invitationId)
		if err != nil {
			log.Panic(err)
		}
		invitationIds = append(invitationIds, invitationId)
	}

	query.Close()

	return
}

// GetInvitationById returns slice of all users ids
func (db DB) GetInvitationById(id int) *Invitation {
	row := db.QueryRow("SELECT * FROM invitations WHERE id = ?", id)

	var invitation Invitation
	err := row.Scan(&invitation.Id, &invitation.Type,
		&invitation.FromUserId, &invitation.ToUserId, &invitation.TimeStamp)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		log.Panic(err)
	}

	return &invitation
}

// DeleteInvitation deletes invitation row in invitations table
func (db DB) DeleteInvitationById(id int) {
	_, err := db.Exec("DELETE FROM invitations WHERE id = ?", id)
	if err != nil {
		log.Panic(err)
	}
}

func (db DB) AddFollowInvitation(fromUserId, toUserId int) FollowStatus {
	_, err := db.Exec(`INSERT INTO invitations 
    	(type, from_user_id, to_user_id, timestamp)
		VALUES (?, ?, ?, ?)`,
		0, fromUserId, toUserId, time.Now().Format(time.RFC3339))
	if err != nil {
		log.Panic(err)
	}
	return RequestToFollow
}

func (db DB) DeleteInvitationByUserId(fromUserId, toUserId int) FollowStatus {
	_, err := db.Exec("DELETE FROM invitations WHERE from_user_id = ? AND to_user_id = ?", fromUserId, toUserId)
	if err != nil {
		log.Panic(err)
	}

	return NotFollowing
}

func (db DB) AddGroupInvitation(fromUserId, toUserId int) MemberStatus {
	_, err := db.Exec(`INSERT INTO invitations 
    	(type, from_user_id, to_user_id, timestamp)
		VALUES (?, ?, ?, ?)`,
		1, fromUserId, toUserId, time.Now().Format(time.RFC3339))
	if err != nil {
		log.Panic(err)
	}
	return RequestedMembership
}
