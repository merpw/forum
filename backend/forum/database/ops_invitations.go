package database

import (
	"log"
	"time"
)

// GetAllInvitations returns slice of all users ids
func (db DB) GetAllInvitations(id int) (invitationIds []int) {
	query, err := db.Query("SELECT id FROM invitations WHERE user_id = ? ORDER BY timestamp DESC", id)
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
		&invitation.AssociatedId, &invitation.UserId, &invitation.TimeStamp)
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

func (db DB) CheckIfInvitationExists(associatedId, userId int) bool {
	query, err := db.Query("SELECT * FROM invitations WHERE associated_id = ? AND user_id = ?",
		associatedId, userId)
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

func (db DB) RequestToFollow(associatedId, userId int) FollowStatus {
	_, err := db.Exec(`INSERT INTO invitations 
    	(type, associated_id, user_id, timestamp)
		VALUES (?, ?, ?, ?)`,
		0, associatedId, userId, time.Now().Format(time.RFC3339))
	if err != nil {
		log.Panic(err)
	}
	return RequestToFollow
}

func (db DB) RevokeInvitation(associatedId, userId int) FollowStatus {
	_, err := db.Exec("DELETE FROM invitations WHERE associated_id = ? AND user_id = ?", associatedId, userId)
	if err != nil {
		log.Panic(err)
	}

	return NotFollowing
}
