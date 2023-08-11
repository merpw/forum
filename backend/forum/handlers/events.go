package handlers

import (
	"backend/common/server"
	. "backend/forum/database"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handlers) eventsId(w http.ResponseWriter, r *http.Request) {
	groupIdStr := strings.Split(r.URL.Path, "/")[3]
	// /api/groups/1/events/2 -> 1

	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	eventIdStr := strings.Split(r.URL.Path, "/")[5]
	// /api/groups/1/events/2 -> 2

	eventId, err := strconv.Atoi(eventIdStr)
	if err != nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	if *h.DB.GetGroupMemberStatus(groupId, h.getUserId(w, r)) != Accepted {
		server.ErrorResponse(w, http.StatusForbidden)
		return
	}

	event := h.DB.GetEventById(eventId)
	if event == nil {
		server.ErrorResponse(w, http.StatusNotFound)
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
	}

	server.SendObject(w, responseBody)

}

func (h *Handlers) eventsIdGoing(w http.ResponseWriter, r *http.Request) {
	groupIdStr := strings.Split(r.URL.Path, "/")[3]
	// /api/groups/1/events/2/going -> 1

	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	eventIdStr := strings.Split(r.URL.Path, "/")[5]
	// /api/groups/1/events/2/going -> 2

	eventId, err := strconv.Atoi(eventIdStr)
	if err != nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	if *h.DB.GetGroupMemberStatus(groupId, h.getUserId(w, r)) != Accepted {
		server.ErrorResponse(w, http.StatusForbidden)
		return
	}

	server.SendObject(w, h.DB.GetEventUserIds(eventId))

}

func (h *Handlers) eventsIdLeave(w http.ResponseWriter, r *http.Request) {
	groupIdStr := strings.Split(r.URL.Path, "/")[3]
	// /api/groups/1/events/2/going -> 1

	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	eventIdStr := strings.Split(r.URL.Path, "/")[5]
	// /api/groups/1/events/2/going -> 2

	eventId, err := strconv.Atoi(eventIdStr)
	if err != nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	if *h.DB.GetGroupMemberStatus(groupId, h.getUserId(w, r)) != Accepted {
		server.ErrorResponse(w, http.StatusForbidden)
		return
	}

	userId := h.getUserId(w, r)
	if !h.DB.GetEventStatusById(eventId, userId) {
		server.ErrorResponse(w, http.StatusForbidden)
		return
	}

	h.DB.DeleteEventMember(eventId, userId)

}
