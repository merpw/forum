package handlers

import (
	"backend/common/server"
	. "backend/forum/database"
	"encoding/json"
	"net/http"
	"strings"
)

// POST /api/groups/create -> group_id
func (h *Handlers) groupsCreate(w http.ResponseWriter, r *http.Request) {
	userId := h.getUserId(w, r)

	requestBody := struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Invite      []int  `json:"invite"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Body is not valid", http.StatusBadRequest)
		return
	}

	requestBody.Title = strings.TrimSpace(requestBody.Title)

	if len(requestBody.Title) < MinTitleLength {
		http.Error(w, "Title is too short", http.StatusBadRequest)
		return
	}

	if len(requestBody.Title) > MaxTitleLength {
		http.Error(w, "Title is too long", http.StatusBadRequest)
		return
	}

	requestBody.Description = strings.TrimSpace(requestBody.Description)

	if len(requestBody.Description) < MinDescriptionLength {
		http.Error(w, "Description is too short", http.StatusBadRequest)
		return
	}

	if len(requestBody.Description) > MaxDescriptionLength {
		http.Error(w, "Description is too long", http.StatusBadRequest)
		return
	}

	for _, id := range requestBody.Invite {
		if id == userId {
			server.ErrorResponse(w, http.StatusBadRequest)
			return
		}

		if h.DB.GetUserById(id) == nil {
			server.ErrorResponse(w, http.StatusNotFound)
			return
		}
	}
	groupId := h.DB.AddGroup(userId, requestBody.Title, requestBody.Description)
	h.DB.AddMembership(int(groupId), userId)

	for _, id := range requestBody.Invite {
		h.DB.AddInvitation(GroupInvite, userId, id)
	}

	server.SendObject(w, groupId)
}
