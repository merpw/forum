package handlers

import (
	"backend/common/server"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// GET /api/invitations
func (h *Handlers) invitations(w http.ResponseWriter, r *http.Request) {
	userId := h.getUserId(w, r)

	server.SendObject(w, h.DB.GetAllInvitations(userId))
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
		Id           int    `json:"id"`
		Type         int    `json:"type"`
		AssociatedId int    `json:"associated_id"`
		UserId       int    `json:"user_id"`
		TimeStamp    string `json:"timestamp"`
	}{
		Id:           invitation.Id,
		Type:         invitation.Type,
		AssociatedId: invitation.AssociatedId,
		UserId:       invitation.UserId,
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
		server.ErrorResponse(w, http.StatusNotFound)
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

	if requestBody.Response {
		h.DB.Follow(invitation.AssociatedId, invitation.UserId)
	}

	h.DB.RespondToInvitation(invitation.Id)
}
