package handlers

import (
	"backend/common/server"
	. "backend/forum/database"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handlers) groups(w http.ResponseWriter, r *http.Request) {
	groupIds := h.DB.GetGroupIdsByMembers()

	server.SendObject(w, groupIds)
}

func (h *Handlers) groupsId(w http.ResponseWriter, r *http.Request) {
	groupIdStr := strings.TrimPrefix(r.URL.Path, "/api/groups/")
	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	group := h.DB.GetGroupById(groupId)
	if group == nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	userId := h.getUserId(w, r)
	responseBody := struct {
		Id           int           `json:"id"`
		Title        string        `json:"title"`
		Description  string        `json:"description"`
		MemberStatus *InviteStatus `json:"member_Status"`
		Members      int           `json:"members"`
	}{
		Id:           group.Id,
		Title:        group.Title,
		Description:  group.Description,
		MemberStatus: h.DB.GetGroupMemberStatus(group.Id, userId),
		Members:      group.Members,
	}

	server.SendObject(w, responseBody)
}

func (h *Handlers) groupsIdPosts(w http.ResponseWriter, r *http.Request) {
	groupIdStr := strings.TrimPrefix(r.URL.Path, "/api/groups/")
	groupIdStr = strings.TrimSuffix(r.URL.Path, "/posts")
	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	group := h.DB.GetGroupById(groupId)
	if group == nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	server.SendObject(w, h.DB.GetGroupPostsById(groupId))
}

// Groups POST endpoints

func (h *Handlers) groupsCreate(w http.ResponseWriter, r *http.Request) {
	userId := h.getUserId(w, r)
	if userId == -1 {
		server.ErrorResponse(w, http.StatusUnauthorized)
		return
	}

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

	for _, id := range requestBody.Invite {
		if id == userId {
			server.ErrorResponse(w, http.StatusBadRequest)
			return
		}
	}
	groupId := h.DB.AddGroup(userId, requestBody.Title, requestBody.Description)
	h.DB.AddMembership(int(groupId), userId)

	for _, id := range requestBody.Invite {
		if h.DB.GetUserById(id) == nil {
			server.ErrorResponse(w, http.StatusNotFound)
			return
		}

		h.DB.AddInvitation(GroupInvite, userId, id)
	}

	server.SendObject(w, groupId)
}

func (h *Handlers) groupsIdJoin(w http.ResponseWriter, r *http.Request) {
	groupIdStr := strings.TrimPrefix(r.URL.Path, "/api/groups/")
	groupIdStr = strings.TrimSuffix(r.URL.Path, "/join")
	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	group := h.DB.GetGroupById(groupId)
	if group == nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	userId := h.getUserId(w, r)

	switch *h.DB.GetGroupMemberStatus(groupId, userId) {
	case Inactive:
		server.SendObject(w, h.DB.AddInvitation(GroupJoin, userId, groupId))
	case Pending:
		server.SendObject(w, h.DB.DeleteInvitationByUserId(GroupJoin, userId, groupId))
	case Accepted:
		server.ErrorResponse(w, http.StatusBadRequest)
	}
}

func (h *Handlers) groupsIdInvite(w http.ResponseWriter, r *http.Request) {
	groupIdStr := strings.TrimPrefix(r.URL.Path, "/api/groups/")
	groupIdStr = strings.TrimSuffix(r.URL.Path, "/join")
	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	group := h.DB.GetGroupById(groupId)
	if group == nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	requestBody := struct {
		UserId int `json:"user_id"`
	}{}

	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Body is not valid", http.StatusBadRequest)
		return
	}

	user := h.DB.GetUserById(requestBody.UserId)
	if user == nil {
		server.ErrorResponse(w, http.StatusNotFound)
	}

	switch *h.DB.GetGroupMemberStatus(groupId, user.Id) {
	case Inactive:
		server.SendObject(w, h.DB.AddInvitation(GroupInvite, groupId, user.Id))
	case Pending:
		server.SendObject(w, h.DB.DeleteInvitationByUserId(GroupInvite, groupId, user.Id))
		// Perhaps ErrorResponse here
	case Accepted:
		server.ErrorResponse(w, http.StatusBadRequest)
	}
}

func (h *Handlers) groupsIdLeave(w http.ResponseWriter, r *http.Request) {
	groupIdStr := strings.TrimPrefix(r.URL.Path, "/api/groups/")
	groupIdStr = strings.TrimSuffix(r.URL.Path, "/leave")
	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	group := h.DB.GetGroupById(groupId)
	if group == nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	userId := h.getUserId(w, r)
	if userId == -1 {
		server.ErrorResponse(w, http.StatusUnauthorized)
		return
	}

	switch *h.DB.GetGroupMemberStatus(groupId, userId) {
	case Inactive:
		server.ErrorResponse(w, http.StatusBadRequest)
	case Pending:
		server.ErrorResponse(w, http.StatusBadRequest)
	case Accepted:
		server.SendObject(w, h.DB.DeleteGroupMembership(groupId, userId))
	}
}
