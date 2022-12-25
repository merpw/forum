package server

import (
	"net/http"
)

func (srv *Server) signUpHandler(w http.ResponseWriter, r *http.Request) {
	errorBasicCheckPOST(w, r, "/api/signup")

	// todo database stuff for "signup" + Error handling during managing data

	sendObject(w, "signup")
}

func (srv *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	errorBasicCheckPOST(w, r, "/api/login")

	// todo database stuff for "login" + Error handling during managing data

	sendObject(w, "login")
}

func (srv *Server) logoutHandler(w http.ResponseWriter, r *http.Request) {
	errorBasicCheckPOST(w, r, "/api/logout")

	// todo database stuff for "logout" + Error handling during managing data

	w.WriteHeader(http.StatusUnauthorized)
	sendObject(w, "logout")
}
