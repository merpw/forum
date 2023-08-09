package handlers

import (
	"backend/common/server"
	"backend/forum/database"
	"net/http"
	"strconv"
	"strings"
)

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
		server.SendObject(w, h.DB.AddGroupInvitation(groupId, userId))
	case database.RequestedMembership:
		server.SendObject(w, h.DB.DeleteInvitationByUserId(groupId, userId))
	case database.Member:
		server.ErrorResponse(w, http.StatusBadRequest)
	}
}
