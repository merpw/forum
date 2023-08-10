package handlers

import (
	"backend/common/server"
	. "backend/forum/database"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handlers) groups(w http.ResponseWriter, r *http.Request) {
	groupIds := h.DB.GetGroupIdsByMembers()

	server.SendObject(w, groupIds)
}

// GET /api/groups/id -> Group
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
		MemberStatus *InviteStatus `json:"member_Status,omitempty"`
		Members      int           `json:"members"`
	}{
		Id:           group.Id,
		Title:        group.Title,
		Description:  group.Description,
		MemberStatus: h.DB.GetGroupMemberStatus(group.Id, userId),
		Members:      h.DB.GetGroupMembersByGroupId(group.Id),
	}

	server.SendObject(w, responseBody)
}

// GET /api/groups/id/posts -> []int{...postIds}
func (h *Handlers) groupsIdPosts(w http.ResponseWriter, r *http.Request) {
	groupIdStr := strings.TrimPrefix(r.URL.Path, "/api/groups/")
	groupIdStr = strings.TrimSuffix(groupIdStr, "/posts")
	// /api/groups/1/posts -> 1

	fmt.Println("id", groupIdStr)
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

// POST /api/groups/id/join -> InviteStatus
func (h *Handlers) groupsIdJoin(w http.ResponseWriter, r *http.Request) {
	groupIdStr := strings.TrimPrefix(r.URL.Path, "/api/groups/")
	groupIdStr = strings.TrimSuffix(groupIdStr, "/join")
	// /api/groups/1/join -> 1

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

// POST /api/groups/id/invite -> InviteStatus
func (h *Handlers) groupsIdInvite(w http.ResponseWriter, r *http.Request) {
	groupIdStr := strings.TrimPrefix(r.URL.Path, "/api/groups/")
	groupIdStr = strings.TrimSuffix(groupIdStr, "/invite")
	// /api/groups/1/invite -> 1

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
		return
	}

	switch *h.DB.GetGroupMemberStatus(groupId, user.Id) {
	case Inactive:
		server.SendObject(w, h.DB.AddInvitation(GroupInvite, groupId, user.Id))
	case Pending:
		server.SendObject(w, h.DB.DeleteInvitationByUserId(GroupInvite, groupId, user.Id))
	case Accepted:
		server.ErrorResponse(w, http.StatusBadRequest)
	}
}

// POST /api/groups/id/leave -> no response
func (h *Handlers) groupsIdLeave(w http.ResponseWriter, r *http.Request) {
	groupIdStr := strings.TrimPrefix(r.URL.Path, "/api/groups/")
	groupIdStr = strings.TrimSuffix(groupIdStr, "/leave")
	// /api/groups/1/leave -> 1

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
		server.ErrorResponse(w, http.StatusBadRequest)
	case Pending:
		server.ErrorResponse(w, http.StatusBadRequest)
	case Accepted:
		server.SendObject(w, h.DB.DeleteGroupMembership(groupId, userId))
	}
}
