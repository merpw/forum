package server

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/mail"
	"strings"
)

func (srv *Server) signupHandler(w http.ResponseWriter, r *http.Request) {
	requestBody := struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Body is not valid", http.StatusBadRequest)
		return
	}
	requestBody.Email = strings.ToLower(requestBody.Email)
	if len(requestBody.Name) < 3 {
		http.Error(w, "Name is too short", http.StatusBadRequest)
		return
	}
	if len(requestBody.Name) > 15 {
		http.Error(w, "Name is too long", http.StatusBadRequest)
		return
	}

	// check if email is a valid email
	_, err = mail.ParseAddress(requestBody.Email)
	if err != nil {
		http.Error(w, "Email is not valid", http.StatusBadRequest)
		return
	}

	if srv.DB.IsEmailTaken(requestBody.Email) {
		http.Error(w, "Email is already taken", http.StatusBadRequest)
		return
	}

	if len(requestBody.Password) < 8 {
		http.Error(w, "Password is too short", http.StatusBadRequest)
		return
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), 10)
	if err != nil {
		http.Error(w, "Password is not valid", http.StatusBadRequest)
		return
	}
	srv.DB.AddUser(requestBody.Name, requestBody.Email, string(encryptedPassword))
}

func (srv *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	// todo database stuff for "login" + Error handling during managing data
	sendObject(w, "login")
}

func (srv *Server) logoutHandler(w http.ResponseWriter, r *http.Request) {
	r.Cookies()
	// todo database stuff for "logout" + Error handling during managing data
	w.WriteHeader(http.StatusUnauthorized)
	sendObject(w, "logout")
}
