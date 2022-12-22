package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (srv *Server) postsHandler(w http.ResponseWriter, r *http.Request) {
	sendObject(w, srv.posts)
}

func (srv *Server) postHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/posts/" {
		srv.postsHandler(w, r)
		return
	}
	sendObject(w, "post "+r.URL.String())
}

func (srv *Server) createHandler(w http.ResponseWriter, r *http.Request) {
	sendObject(w, "create")
}

func (srv *Server) signUpHandler(w http.ResponseWriter, r *http.Request) {
	sendObject(w, "signup")
}

func (srv *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	sendObject(w, "login")
}

func (srv *Server) logoutHandler(w http.ResponseWriter, r *http.Request) {
	sendObject(w, "logout")
}

// errorResponse responses with specified error code in format "404 Not Found"
func errorResponse(w http.ResponseWriter, code int) {
	http.Error(w, fmt.Sprintf("%v %v", code, http.StatusText(code)), code)
}

// sendObject sends object to http.ResponseWriter
//
// calls errorResponse(500) if error happened
func sendObject(w http.ResponseWriter, object any) {
	w.Header().Set("Content-Type", "application/json")
	objJson, err := json.Marshal(object)
	if err != nil {
		log.Println(err)
		errorResponse(w, 500)
		return
	}
	_, err = w.Write(objJson)
	if err != nil {
		log.Println(err)
		errorResponse(w, 500)
		return
	}
}
