package handlers

import (
	"backend/common/server"
	. "backend/forum/database"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handlers) eventsId(w http.ResponseWriter, r *http.Request) {

	eventIdStr := strings.Split(r.URL.Path, "/")[3]
	// /api/groups/1/events/2 -> 2

	eventId, err := strconv.Atoi(eventIdStr)
	if err != nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	event := h.DB.GetEventById(eventId)

	if event == nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	userId := h.getUserId(w, r)

	if *h.DB.GetGroupMemberStatus(event.GroupId, userId) != InviteStatusAccepted {
		server.ErrorResponse(w, http.StatusForbidden)
		return
	}

	var responseBody = SafeEvent{
		Id:          event.Id,
		GroupId:     event.GroupId,
		CreatedBy:   event.CreatedBy,
		Title:       event.Title,
		Description: event.Description,
		TimeAndDate: event.TimeAndDate,
		Timestamp:   event.Timestamp,
		Status:      h.DB.GetEventStatus(eventId, userId),
	}

	server.SendObject(w, responseBody)
}

func (h *Handlers) eventsIdMembers(w http.ResponseWriter, r *http.Request) {
	eventIdStr := strings.Split(r.URL.Path, "/")[3]
	// /api/events/2/members -> 2

	eventId, err := strconv.Atoi(eventIdStr)
	if err != nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	event := h.DB.GetEventById(eventId)

	if event == nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	userId := h.getUserId(w, r)

	if *h.DB.GetGroupMemberStatus(event.GroupId, userId) != InviteStatusAccepted {
		server.ErrorResponse(w, http.StatusForbidden)
		return
	}

	server.SendObject(w, h.DB.GetEventMembers(eventId))
}

func (h *Handlers) eventsIdLeave(w http.ResponseWriter, r *http.Request) {
	eventIdStr := strings.Split(r.URL.Path, "/")[3]
	// /api/events/2/leave -> 2

	eventId, err := strconv.Atoi(eventIdStr)
	if err != nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	event := h.DB.GetEventById(eventId)

	if event == nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	userId := h.getUserId(w, r)

	if *h.DB.GetGroupMemberStatus(event.GroupId, userId) != InviteStatusAccepted {
		server.ErrorResponse(w, http.StatusForbidden)
		return
	}

	h.DB.DeleteEventMember(eventId, userId)
}

func (h *Handlers) eventsIdGoing(w http.ResponseWriter, r *http.Request) {
	eventIdStr := strings.Split(r.URL.Path, "/")[3]
	// /api/events/2/going -> 2

	eventId, err := strconv.Atoi(eventIdStr)
	if err != nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	event := h.DB.GetEventById(eventId)
	if event == nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	userId := h.getUserId(w, r)

	if *h.DB.GetGroupMemberStatus(event.GroupId, userId) != InviteStatusAccepted {
		server.ErrorResponse(w, http.StatusForbidden)
		return
	}

	h.DB.AddEventMember(eventId, userId)
}
