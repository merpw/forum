package handlers

import "net/http"

func (handlers *Handlers) checkSession(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		errorResponse(w, http.StatusBadRequest)
		return
	}
	userId := handlers.DB.CheckSession(token)
	if userId == -1 {
		sendObject(w, struct {
			Error string `json:"error"`
		}{
			Error: "Invalid token",
		})
		return
	}
	sendObject(w, userId)
}
