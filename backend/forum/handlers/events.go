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

	if group := h.DB.GetGroupById(groupId); group == nil {
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

	userId := h.getUserId(w, r)

	if code := h.validateEventsId(groupId, eventId, userId); code != http.StatusOK {
		server.ErrorResponse(w, code)
		return
	}

	event := h.DB.GetEventById(eventId)

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

func (h *Handlers) eventsIdUsers(w http.ResponseWriter, r *http.Request) {
	groupIdStr := strings.Split(r.URL.Path, "/")[3]
	// /api/groups/1/events/2/users -> 1

	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	eventIdStr := strings.Split(r.URL.Path, "/")[5]
	// /api/groups/1/events/2/users -> 2

	eventId, err := strconv.Atoi(eventIdStr)
	if err != nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	userId := h.getUserId(w, r)

	if code := h.validateEventsId(groupId, eventId, userId); code != http.StatusOK {
		server.ErrorResponse(w, code)
		return
	}

	server.SendObject(w, h.DB.GetEventMembers(eventId))
}

func (h *Handlers) eventsIdLeave(w http.ResponseWriter, r *http.Request) {
	groupIdStr := strings.Split(r.URL.Path, "/")[3]
	// /api/groups/1/events/2/leave -> 1

	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	eventIdStr := strings.Split(r.URL.Path, "/")[5]
	// /api/groups/1/events/2/leave -> 2

	eventId, err := strconv.Atoi(eventIdStr)
	if err != nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	userId := h.getUserId(w, r)

	if code := h.validateEventsId(groupId, eventId, userId); code != http.StatusOK {
		server.ErrorResponse(w, code)
		return
	}

	h.DB.DeleteEventMember(eventId, userId)
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

	userId := h.getUserId(w, r)

	if code := h.validateEventsId(groupId, eventId, userId); code != http.StatusOK {
		server.ErrorResponse(w, code)
		return
	}

	h.DB.AddEventMember(eventId, userId)
}

func (h *Handlers) validateEventsId(groupId, eventId, userId int) int {
	if group := h.DB.GetGroupById(groupId); group == nil {
		return http.StatusNotFound
	}

	if event := h.DB.GetEventById(eventId); event == nil {
		return http.StatusNotFound
	}

	if *h.DB.GetGroupMemberStatus(groupId, userId) != Accepted {
		return http.StatusForbidden
	}

	return http.StatusOK
}
