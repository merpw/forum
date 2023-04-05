package server

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (srv *Server) signupHandler(w http.ResponseWriter, r *http.Request) {
	if srv.getUserId(w, r) != -1 {
		http.Error(w, "You are already logged in", http.StatusBadRequest)
		return
	}

	requestBody := struct {
		Name      string `json:"name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Age       string `json:"age"`
		Gender    string `json:"gender"`
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
	if requestBody.Name != strings.TrimSpace(requestBody.Name) {
		http.Error(w, "Name is not valid", http.StatusBadRequest)
		return
	}
	// here is checking name is already in use,
	// statement above checks the incoming name do not include spaces before and after
	if srv.DB.IsNameTaken(requestBody.Name) {
		http.Error(w, "Name is already in use", http.StatusBadRequest)
		return
	}
	// check if email is a valid email
	_, err = mail.ParseAddress(requestBody.Email)
	if err != nil || requestBody.Email != strings.TrimSpace(requestBody.Email) {
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
	srv.DB.AddUser(
		requestBody.Name,
		requestBody.Email,
		string(encryptedPassword),
		requestBody.FirstName,
		requestBody.LastName,
		requestBody.Age,
		requestBody.Gender)
}

func (srv *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	if srv.getUserId(w, r) != -1 {
		http.Error(w, "You are already logged in", http.StatusBadRequest)
		return
	}

	requestBody := struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Body is not valid", http.StatusBadRequest)
		return
	}

	user := srv.DB.GetUserByLogin(strings.ToLower(requestBody.Login))
	if user == nil {
		http.Error(w, "Invalid login or password", http.StatusBadRequest)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password))
	if err != nil {
		http.Error(w, "Invalid password", http.StatusBadRequest)
		return
	}
	token, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	expire := time.Now().Add(24 * time.Hour)

	srv.DB.AddSession(token.String(), int(expire.Unix()), user.Id)

	http.SetCookie(w, &http.Cookie{
		Name:    "forum-token",
		Value:   token.String(),
		Expires: expire,
		Path:    "/",
		//	TODO: add secure on production
	})
}

func (srv *Server) logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("forum-token")
	if err != nil {
		errorResponse(w, http.StatusUnauthorized)
		return
	}
	srv.DB.RemoveSession(cookie.Value)

	http.SetCookie(w, &http.Cookie{
		Name:    "forum-token",
		Value:   "",
		Expires: time.Now(),
		Path:    "/",
	})
}

func (srv *Server) getUserId(w http.ResponseWriter, r *http.Request) int {
	cookie, err := r.Cookie("forum-token")
	if err != nil {
		return -1
	}
	userId := srv.DB.CheckSession(cookie.Value)

	if userId == -1 {
		http.SetCookie(w, &http.Cookie{
			Name:    "forum-token",
			Value:   "",
			Expires: time.Now(),
			Path:    "/",
		})
		return -1
	}

	return userId
}
