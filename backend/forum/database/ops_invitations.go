package database

import (
	"log"
	"time"
)

// GetAllInvitations returns slice of all users ids
func (db DB) GetAllInvitations(toUserId int) (invitationIds []int) {
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
	query, err := db.Query("SELECT * FROM invitations WHERE id = ?", id)
	if err != nil {
		log.Panic(err)
	}

	if !query.Next() {
		return nil
	}

	var invitation Invitation

	err = query.Scan(&invitation.Id, &invitation.Type,
		&invitation.FromUserId, &invitation.ToUserId, &invitation.TimeStamp)
	if err != nil {
		log.Panic(err)
	}

	query.Close()

	return &invitation
}

// RespondToInvitation deletes invitation row in invitations table
func (db DB) RespondToInvitation(id int) {
	_, err := db.Exec("DELETE FROM invitations WHERE id = ?", id)
	if err != nil {
		log.Panic(err)
	}
}

func (db DB) CheckIfInvitationExists(fromUserId, toUserId int) bool {
	query, err := db.Query("SELECT * FROM invitations WHERE from_user_id = ? AND to_user_id = ?",
		fromUserId, toUserId)
	if err != nil {
		log.Panic(err)
	}
	var invitationExists bool

	if !query.Next() {
		invitationExists = false
	} else {
		invitationExists = true
	}

	query.Close()

	return invitationExists

}

func (db DB) RequestToFollow(fromUserId, toUserId int) FollowStatus {
	_, err := db.Exec(`INSERT INTO invitations 
    	(type, from_user_id, to_user_id, timestamp)
		VALUES (?, ?, ?, ?)`,
		0, fromUserId, toUserId, time.Now().Format(time.RFC3339))
	if err != nil {
		log.Panic(err)
	}
	return RequestToFollow
}

func (db DB) RevokeInvitation(fromUserId, toUserId int) FollowStatus {
	_, err := db.Exec("DELETE FROM invitations WHERE from_user_id = ? AND to_user_id = ?", fromUserId, toUserId)
	if err != nil {
		log.Panic(err)
	}

	return NotFollowing
}
