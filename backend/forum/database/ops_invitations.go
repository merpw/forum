package database

import "log"

// GetAllInvitations returns slice of all users ids
func (db DB) GetAllInvitations(id int) (invitationIds []int) {
	query, err := db.Query("SELECT id FROM invitations WHERE user_id = ? ORDER BY timestamp DESC", id)
	if err != nil {
		log.Panic(err)
	}

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
