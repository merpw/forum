package handlers

import (
	"backend/common/server"
	. "backend/forum/database"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handlers) groups(w http.ResponseWriter, r *http.Request) {
	groupIds := h.DB.GetTopGroups()

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
		MemberStatus *InviteStatus `json:"member_status,omitempty"`
		MemberCount  int           `json:"member_count"`
		CreatorId    int           `json:"creator_id"`
	}{
		Id:           group.Id,
		Title:        group.Title,
		Description:  group.Description,
		CreatorId:    group.CreatorId,
		MemberStatus: h.DB.GetGroupMemberStatus(group.Id, userId),
		MemberCount:  h.DB.GetGroupMemberCount(group.Id),
	}

	server.SendObject(w, responseBody)
}

// GET /api/groups/id/posts -> []int{...postIds}
func (h *Handlers) groupsIdPosts(w http.ResponseWriter, r *http.Request) {
	groupIdStr := strings.TrimPrefix(r.URL.Path, "/api/groups/")
	groupIdStr = strings.TrimSuffix(groupIdStr, "/posts")
	// /api/groups/1/posts -> 1

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

func (h *Handlers) groupsIdMembers(w http.ResponseWriter, r *http.Request) {
	groupIdStr := strings.TrimPrefix(r.URL.Path, "/api/groups/")
	groupIdStr = strings.TrimSuffix(groupIdStr, "/members")
	// /api/groups/1/members -> 1

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

	userId := r.Context().Value(userIdCtxKey).(int)

	withPending := r.URL.Query().Has("withPending")

	if userId == -1 {
		// internal request
		server.SendObject(w, h.DB.GetGroupMembers(groupId, withPending))
		return
	}

	if *h.DB.GetGroupMemberStatus(groupId, userId) != Accepted {
		server.ErrorResponse(w, http.StatusForbidden)
		return
	}

	server.SendObject(w, h.DB.GetGroupMembers(groupId, withPending))
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

	fromUserId := h.getUserId(w, r)
	toUserId := h.DB.GetGroupCreatorId(groupId)

	switch *h.DB.GetGroupMemberStatus(groupId, fromUserId) {
	case Inactive:
		server.SendObject(w, h.DB.AddInvitation(
			GroupJoin,
			fromUserId,
			*toUserId,
			sql.NullInt64{Int64: int64(groupId), Valid: true}))
	case Pending:
		server.SendObject(w, h.DB.DeleteInvitationByUserId(
			GroupJoin,
			fromUserId,
			*toUserId,
			sql.NullInt64{Int64: int64(groupId), Valid: true}))
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

	toUser := h.DB.GetUserById(requestBody.UserId)
	if toUser == nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	fromUserId := h.getUserId(w, r)

	switch *h.DB.GetGroupMemberStatus(groupId, toUser.Id) {
	case Inactive:
		server.SendObject(w, h.DB.AddInvitation(
			GroupInvite,
			fromUserId,
			toUser.Id,
			sql.NullInt64{Int64: int64(groupId), Valid: true}))
	case Pending:
		fallthrough
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
	if userId == *h.DB.GetGroupCreatorId(groupId) {
		http.Error(w, "Creator can't leave group", http.StatusBadRequest)
		return
	}

	switch *h.DB.GetGroupMemberStatus(groupId, userId) {
	case Inactive:
		server.ErrorResponse(w, http.StatusBadRequest)
	case Pending:
		server.ErrorResponse(w, http.StatusBadRequest)
	case Accepted:
		h.DB.DeleteGroupMembership(groupId, userId)
		server.SendObject(w, Inactive)
	}
}
