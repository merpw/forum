package handlers

import (
	"backend/common/server"
	"net/http"
)

func (handlers *Handlers) checkSession(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		server.ErrorResponse(w, http.StatusBadRequest)
		return
	}
	userId := handlers.DB.CheckSession(token)
	if userId == -1 {
		server.SendObject(w, struct {
			Error string `json:"error"`
		}{
			Error: "Invalid token",
		})
		return
	}
	server.SendObject(w, userId)
}
