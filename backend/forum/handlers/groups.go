package handlers

import (
	"backend/common/server"
	"backend/forum/database"
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
		Id           int                    `json:"id"`
		Title        string                 `json:"title"`
		Description  string                 `json:"description"`
		MemberStatus *database.MemberStatus `json:"member_Status"`
		Members      int                    `json:"members"`
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
	case database.NotMember:
		server.SendObject(w, h.DB.AddGroupInvitation(userId, groupId))
	case database.RequestedMembership:
		server.SendObject(w, h.DB.DeleteInvitationByUserId(userId, groupId))
	case database.Member:
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
	case database.NotMember:
		server.SendObject(w, h.DB.AddGroupInvitation(groupId, user.Id))
	case database.RequestedMembership:
		server.ErrorResponse(w, http.StatusBadRequest)
	case database.Member:
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
	case database.NotMember:
		server.ErrorResponse(w, http.StatusBadRequest)
	case database.RequestedMembership:
		server.ErrorResponse(w, http.StatusBadRequest)
	case database.Member:
		server.SendObject(w, h.DB.DeleteGroupMembership(groupId, userId))
	}
}
