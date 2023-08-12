package handlers

import (
	"backend/common/server"
	. "backend/forum/database"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// GET /api/invitations
func (h *Handlers) invitations(w http.ResponseWriter, r *http.Request) {
	userId := h.getUserId(w, r)

	server.SendObject(w, h.DB.GetUserInvitations(userId))
}

// GET /api/invitations/(/d+)
func (h *Handlers) invitationsId(w http.ResponseWriter, r *http.Request) {
	invitationIdStr := strings.TrimPrefix(r.URL.Path, "/api/invitations/")
	// /api/invitations/1 -> 1

	invitationId, err := strconv.Atoi(invitationIdStr)
	if err != nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	invitation := h.DB.GetInvitationById(invitationId)
	if invitation == nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	respBody := struct {
		Id           int        `json:"id"`
		Type         InviteType `json:"type"`
		FromUserId   int        `json:"from_user_id"`
		ToUserId     int        `json:"to_user_id"`
		AssociatedId int        `json:"associated_id,omitempty"`
		TimeStamp    string     `json:"timestamp"`
	}{
		Id:           invitation.Id,
		Type:         invitation.Type,
		FromUserId:   invitation.FromUserId,
		ToUserId:     invitation.ToUserId,
		AssociatedId: int(invitation.AssociatedId.Int64),
		TimeStamp:    invitation.TimeStamp,
	}
	server.SendObject(w, respBody)
}

// POST /api/invitations/(/d+)/respond
func (h *Handlers) invitationsIdRespond(w http.ResponseWriter, r *http.Request) {
	invitationIdStr := strings.TrimPrefix(r.URL.Path, "/api/invitations/")
	invitationIdStr = strings.TrimSuffix(invitationIdStr, "/respond")
	// /api/invitations/1/respond -> 1

	invitationId, err := strconv.Atoi(invitationIdStr)
	if err != nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	invitation := h.DB.GetInvitationById(invitationId)
	if invitation == nil {
		http.Error(w, "Invitation not found", http.StatusNotFound)
		return
	}

	requestBody := struct {
		Response bool `json:"response"`
	}{}

	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Body is not valid", http.StatusBadRequest)
		return
	}
	groupId := int(invitation.AssociatedId.Int64)

	if requestBody.Response {
		switch invitation.Type {
		case FollowUser:
			h.DB.AddFollower(invitation.FromUserId, invitation.ToUserId)
		case GroupInvite:
			h.DB.AddMembership(groupId, invitation.ToUserId)
		case GroupJoin:
			h.DB.AddMembership(groupId, invitation.FromUserId)
		}

	}

	h.DB.DeleteInvitationById(invitation.Id)
}
