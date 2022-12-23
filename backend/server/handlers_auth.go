package server

import "net/http"

func (srv *Server) signUpHandler(w http.ResponseWriter, r *http.Request) {
	sendObject(w, "signup")
}

func (srv *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	sendObject(w, "login")
}

func (srv *Server) logoutHandler(w http.ResponseWriter, r *http.Request) {
	sendObject(w, "logout")
}
