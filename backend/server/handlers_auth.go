package server

import (
	"fmt"
	"log"
	"net/http"
)

func (srv *Server) signUpHandler(w http.ResponseWriter, r *http.Request) {

	defer func() {
		if err := recover(); err != nil {
			log.Printf("ERROR %d. %v\n", http.StatusInternalServerError, err)
			errorResponse(w, http.StatusInternalServerError) // 500 ERROR
		}
	}()

	if r.URL.Path != "/api/signup" {
		log.Printf("ERROR %d. r.URL.Path = %s != \"/api/signup\"\n", http.StatusNotFound, r.URL.Path)
		errorResponse(w, http.StatusNotFound) // 404 ERROR
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("ERROR %d. ParseForm() err: %v\n", http.StatusBadRequest, err)
		errorResponse(w, http.StatusBadRequest) // 400 ERROR
		return
	}

	if r.Method != http.MethodPost { // not POST method case
		log.Printf("ERROR %d. %v\n", http.StatusMethodNotAllowed, fmt.Sprintf("request method %s is inappropriate for the URL %s", r.Method, r.URL.Path))
		errorResponse(w, http.StatusMethodNotAllowed)
		return
	}

	// todo database stuff for signup + Error handling during managing data

	sendObject(w, "signup")
}

func (srv *Server) loginHandler(w http.ResponseWriter, r *http.Request) {

	defer func() {
		if err := recover(); err != nil {
			log.Printf("ERROR %d. %v\n", http.StatusInternalServerError, err)
			errorResponse(w, http.StatusInternalServerError) // 500 ERROR
		}
	}()

	if r.URL.Path != "/api/login" {
		log.Printf("ERROR %d. r.URL.Path = %s != \"/api/login\"\n", http.StatusNotFound, r.URL.Path)
		errorResponse(w, http.StatusNotFound) // 404 ERROR
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("ERROR %d. ParseForm() err: %v\n", http.StatusBadRequest, err)
		errorResponse(w, http.StatusBadRequest) // 400 ERROR
		return
	}

	if r.Method != http.MethodPost { // not POST method case
		log.Printf("ERROR %d. %v\n", http.StatusMethodNotAllowed, fmt.Sprintf("request method %s is inappropriate for the URL %s", r.Method, r.URL.Path))
		errorResponse(w, http.StatusMethodNotAllowed)
		return
	}

	// todo database stuff for login + Error handling during managing data

	sendObject(w, "login")
}

func (srv *Server) logoutHandler(w http.ResponseWriter, r *http.Request) {

	defer func() {
		if err := recover(); err != nil {
			log.Printf("ERROR %d. %v\n", http.StatusInternalServerError, err)
			errorResponse(w, http.StatusInternalServerError) // 500 ERROR
		}
	}()

	if r.URL.Path != "/api/logout" {
		log.Printf("ERROR %d. r.URL.Path = %s != \"/api/logout\"\n", http.StatusNotFound, r.URL.Path)
		errorResponse(w, http.StatusNotFound) // 404 ERROR
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("ERROR %d. ParseForm() err: %v\n", http.StatusBadRequest, err)
		errorResponse(w, http.StatusBadRequest) // 400 ERROR
		return
	}

	if r.Method != http.MethodPost { // not POST method case
		log.Printf("ERROR %d. %v\n", http.StatusMethodNotAllowed, fmt.Sprintf("request method %s is inappropriate for the URL %s", r.Method, r.URL.Path))
		errorResponse(w, http.StatusMethodNotAllowed)
		return
	}

	// todo database stuff for logout + Error handling during managing data

	w.WriteHeader(http.StatusUnauthorized)
	sendObject(w, "logout")
}
