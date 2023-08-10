package handlers

import (
	"backend/common/server"
	. "backend/forum/database"
	"backend/forum/external"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (h *Handlers) eventsCreate(w http.ResponseWriter, r *http.Request) {
	groupIdStr := strings.TrimPrefix(r.URL.Path, "/api/groups/")
	groupIdStr = strings.TrimSuffix(groupIdStr, "/events/create")
	// /api/groups/1/events/create -> 1
	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		server.ErrorResponse(w, http.StatusNotFound)
	}

	if *h.DB.GetGroupMemberStatus(groupId, h.getUserId(w, r)) != Accepted {
		server.ErrorResponse(w, http.StatusForbidden)
	}

	var event SafeEvent
	if err = json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Body is not valid", http.StatusBadRequest)
		return
	}

	event.Title = strings.TrimSpace(event.Title)
	if len(event.Title) < MinTitleLength {
		http.Error(w, "title is too short", http.StatusBadRequest)
		return
	}

	if len(event.Title) > MaxTitleLength {
		http.Error(w, "title is too long", http.StatusBadRequest)
		return
	}

	event.Description = strings.TrimSpace(event.Description)
	if len(event.Description) < MinDescriptionLength {
		http.Error(w, "description is too short", http.StatusBadRequest)
		return
	}

	if len(event.Description) > MaxDescriptionLength {
		http.Error(w, "description is too long", http.StatusBadRequest)
		return
	}

	timeAndDate, timeErr := time.Parse("2006-01-02", event.TimeAndDate)
	if timeErr != nil {
		http.Error(w, "time and date is invalid", http.StatusBadRequest)
		return
	}

	if timeAndDate.Before(time.Now()) {
		http.Error(w, "time and date is invalid", http.StatusBadRequest)
		return
	}

	userId := h.getUserId(w, r)
	id := h.DB.AddEvent(groupId, userId, event.Title, event.Description, event.TimeAndDate)

	h.DB.AddInvitation(Event, userId, userId, sql.NullInt64{Int64: int64(id), Valid: true})
	server.SendObject(w, id)

	external.RevalidateURL(fmt.Sprintf("/groups/%d/events/%d", groupId, id))

	external.RevalidateURL("/")

}
