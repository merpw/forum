package database

import (
	"log"
	"time"
)

func (db DB) AddEvent(groupId, createdBy int, title, description, timeAndDate string) int {
	result, err := db.Exec(` 
			INSERT INTO events
			(group_id, creator_id, title, description, time_and_date)
			VALUES (?, ?, ?, ?, ?)`, groupId, createdBy, title, description, timeAndDate)
	if err != nil {
		log.Panic(err)
	}

	id, newErr := result.LastInsertId()
	if newErr != nil {
		log.Panic(newErr)
	}

	_, invitationErr := db.Exec(`
    INSERT INTO invitations
        (type, from_user_id, to_user_id, timestamp, associated_id)
        SELECT ?, ?, user_id, ?, ?
        FROM group_members
        WHERE group_id = ? AND user_id IS NOT ?`,
		3, createdBy, time.Now().Format(time.RFC3339), id, groupId, createdBy)

	if invitationErr != nil {
		log.Panic(err)
	}
	return int(id)
}

func (db DB) DeleteAllEventInvites(groupId, userId int) {
	_, err := db.Exec(`
    DELETE FROM invitations WHERE TYPE = 3 AND
        associated_id IN (SELECT id FROM events WHERE group_id = ?) AND to_user_id = ?`, groupId, userId)

	if err != nil {
		log.Panic(err)
	}
}

func (db DB) AddEventMember(eventId, userId int) {
	_, err := db.Exec(`
  	INSERT OR IGNORE INTO event_members
    	(event_id, user_id) VALUES (?, ?)`, eventId, userId)

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

func (db DB) GetEventMembers(eventId int) []int {
	query, err := db.Query("SELECT user_id FROM event_members WHERE event_id = ?", eventId)
	if err != nil {
		log.Panic(err)
	}

	var ids = make([]int, 0)

	for query.Next() {
		var id int

		err := query.Scan(&id)
		if err != nil {
			log.Panic(err)
		}

		ids = append(ids, id)
	}

	query.Close()

	return ids
}

func (db DB) GetEventIdsByGroupId(groupId int) []int {

	query, err := db.Query("SELECT id FROM events WHERE group_id = ? ORDER BY id DESC", groupId)
	if err != nil {
		log.Panic(err)
	}

	var ids = make([]int, 0)

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

func (db DB) GetEventStatus(eventId, userId int) *int {
	row := db.QueryRow(`    
		SELECT CASE 
		WHEN (
			SELECT 1 FROM invitations WHERE to_user_id = :userId AND associated_id = :eventId) THEN 0
		ELSE (
			SELECT CASE 
			WHEN (
				SELECT 1 FROM event_members WHERE user_id = :userId AND event_id = :eventId) THEN 1
			ELSE 2
			END
		)
		END 
		AS event_status`, eventId, userId)

	var status = new(int)
	err := row.Scan(status)
	if err != nil {
		log.Panic(err)
	}

	return status
}

func (db DB) DeleteEventMember(eventId, userId int) {
	_, err := db.Exec("DELETE FROM event_members WHERE event_id = ? AND user_id = ?", eventId, userId)
	if err != nil {
		log.Panic(err)
	}
}
