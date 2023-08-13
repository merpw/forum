package database

import (
	"log"
	"time"
)

func (db DB) AddEvent(groupId, createdBy int, title, description, timeAndDate string) int {
	result, err := db.Exec(` 
			INSERT INTO events
			(title, description, time_and_date, timestamp, creator_id, group_id)
			VALUES (?, ?, ?, ?, ?, ?)`,
		createdBy, groupId, timeAndDate, time.Now().Format(time.RFC3339), title, description)
	if err != nil {
		log.Panic(err)
	}

	id, newErr := result.LastInsertId()
	if newErr != nil {
		log.Panic(newErr)
	}

	db.AddUserToEvent(int(id), createdBy)

	_, invitationErr := db.Exec(`
    INSERT INTO invitations
        (type, from_user_id, to_user_id, timestamp, associated_id)
        SELECT ?, ?, user_id, ?, ?
        FROM group_members
        WHERE group_id = ? AND user_id IS NOT ?`,
		2, createdBy, time.Now().Format(time.RFC3339), id, groupId, createdBy)

	if invitationErr != nil {
		log.Panic(err)
	}
	return int(id)
}

func (db DB) DeleteAllEventInvites(groupId, userId int) {
	_, err := db.Exec(`
    DELETE * FROM invitations WHERE TYPE = 2 AND associated_id = ? AND to_user_id = ?`,
		2, userId, groupId)

	if err != nil {
		log.Panic(err)
	}
}

func (db DB) AddUserToEvent(eventId, userId int) {
	_, err := db.Exec(`
			INSERT INTO event_members
			(event_id, user_id, timestamp) VALUES (?, ?, ?)`, eventId, userId, time.Now().Format(time.RFC3339))
	if err != nil {
		log.Panic(err)
	}
}

func (db DB) GetEventById(eventId int) *EventData {
	var e = EventData{}
	err := db.QueryRow(
		"SELECT * FROM events WHERE id = ?", eventId).
		Scan(&e.Id, &e.GroupId, &e.CreatedBy, &e.Title, &e.Description, &e.TimeAndDate, &e.Timestamp)
	if err != nil {
		return nil
	}

	return &e
}

func (db DB) GetEventUserIds(eventId int) []int {
	query, err := db.Query("SELECT user_id FROM event_members WHERE event_id = ?", eventId)
	if err != nil {
		log.Panic(err)
	}

	var ids []int

	for query.Next() {
		var id int

		err = query.Scan(&id)
		if err != nil {
			log.Panic(err)
		}

		ids = append(ids, id)
	}

	query.Close()

	return ids
}

func (db DB) GetEventIdsByGroupId(groupId int) []int {

	query, err := db.Query("SELECT id FROM events WHERE group_id = ?", groupId)
	if err != nil {
		log.Panic(err)
	}

	var ids []int

	for query.Next() {
		var id int

		err = query.Scan(&id)
		if err != nil {
			log.Panic(err)
		}

		ids = append(ids, id)
	}

	query.Close()

	return ids

}

func (db DB) GetEventStatusById(eventId, userId int) bool {
	return db.QueryRow("SELECT FROM event_members WHERE event_id = ? AND user_id = ?", eventId, userId) != nil
}

func (db DB) DeleteEventMember(eventId, userId int) {
	_, err := db.Exec("DELETE FROM event_members WHERE event_id = ? AND user_id = ?", eventId, userId)
	if err != nil {
		log.Panic(err)
	}
}
