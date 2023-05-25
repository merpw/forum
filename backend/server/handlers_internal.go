package server

import (
  "net/http"
  "fmt"
)

func (srv *Server) checkSessionHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		errorResponse(w, http.StatusBadRequest)
		return
	}
	userId := srv.DB.CheckSession(token)
  fmt.Println("userId", userId)
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
