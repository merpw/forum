package server

import "net/http"

func (srv *Server) checkSessionHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		errorResponse(w, http.StatusBadRequest)
		return
	}
	userId := srv.DB.CheckSession(token)
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
