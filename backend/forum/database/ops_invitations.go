package database

import (
	"database/sql"
	"log"
	"time"
)

// GetUserInvitations returns slice of all users ids
func (db DB) GetUserInvitations(toUserId int) (invitationIds []int) {
	query, err := db.Query("SELECT id FROM invitations WHERE to_user_id = ? ORDER BY id DESC", toUserId)
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
	inv := &Invitation{}
	err := db.QueryRow("SELECT * FROM invitations WHERE id = ?", id).Scan(
		&inv.Id, &inv.Type, &inv.FromUserId, &inv.ToUserId, &inv.TimeStamp, &inv.AssociatedId)
	if err != nil {
		return nil
	}
	return inv
}

// DeleteInvitation deletes invitation row in invitations table
func (db DB) DeleteInvitationById(id int) {
	_, err := db.Exec("DELETE FROM invitations WHERE id = ?", id)
	if err != nil {
		log.Panic(err)
	}
}

func (db DB) DeleteFollowRequest(fromUserId, toUserId int) InviteStatus {
	_, err := db.Exec(`
		DELETE FROM invitations
			WHERE type = 0
			AND from_user_id = ? 
			AND to_user_id = ?`, fromUserId, toUserId)
	if err != nil {
		log.Panic(err)
	}

	return InviteStatusUnset
}

func (db DB) DeleteInvitationByUserId(invType InviteType, fromUserId, toUserId int,
	associatedId sql.NullInt64) InviteStatus {
	_, err := db.Exec(`DELETE FROM invitations
       						WHERE type = ?
       						AND from_user_id = ? 
       						AND to_user_id = ? 
       						AND associated_id = ?`, invType, fromUserId, toUserId, associatedId)
	if err != nil {
		log.Panic(err)
	}

	return InviteStatusUnset
}

func (db DB) AddInvitation(inviteType InviteType, fromUserId, toUserId int, associatedId sql.NullInt64) InviteStatus {
	_, err := db.Exec(`INSERT INTO invitations 
    	(type, from_user_id, to_user_id, associated_id, timestamp)
		VALUES (?, ?, ?, ?, ?)`,
		inviteType, fromUserId, toUserId, associatedId, time.Now().Format(time.RFC3339))
	if err != nil {
		log.Panic(err)
	}
	return InviteStatusPending
}
