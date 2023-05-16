package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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
		DoB       string `json:"dob"`
		Gender    string `json:"gender"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Body is not valid", http.StatusBadRequest)
		return
	}

	requestBody.Name = strings.TrimSpace(requestBody.Name)
	if len(requestBody.Name) < 3 {
		http.Error(w, "Username is too short", http.StatusBadRequest)
		return
	}

	if len(requestBody.Name) > 15 {
		http.Error(w, "Username is too long", http.StatusBadRequest)
		return
	}
	if srv.DB.IsNameTaken(requestBody.Name) {
		http.Error(w, "Username is already in use", http.StatusBadRequest)
		return
	}
	requestBody.FirstName = strings.TrimSpace(requestBody.FirstName)
	if requestBody.FirstName == "" {
		http.Error(w, "First name is not valid", http.StatusBadRequest)
		return
	}
	if len(requestBody.FirstName) > 15 {
		http.Error(w, "First name is too long", http.StatusBadRequest)
		return
	}
	requestBody.LastName = strings.TrimSpace(requestBody.LastName)
	if requestBody.LastName == "" {
		http.Error(w, "Last name is not valid", http.StatusBadRequest)
		return
	}
	if len(requestBody.LastName) > 15 {
		http.Error(w, "Last name is too long", http.StatusBadRequest)
		return
	}

	if requestBody.DoB == "" {
		http.Error(w, "Date of birth is not valid", http.StatusBadRequest)
		return
	}

	dob, err := time.Parse("2006-01-02", requestBody.DoB)
	if err != nil {
		http.Error(w, "Date of birth is not valid", http.StatusBadRequest)
		return
	}
	now := time.Now()
	minDoB := time.Date(1900, time.January, 1, 0, 0, 0, 0, time.UTC)

	if dob.After(now) || dob.Before(minDoB) {
		http.Error(w, "Date of birth is not valid", http.StatusBadRequest)
		return
	}

	if requestBody.Gender != "male" && requestBody.Gender != "female" && requestBody.Gender != "other" {
		http.Error(w, "Gender is not valid", http.StatusBadRequest)
		return
	}

	requestBody.Email = strings.TrimSpace(strings.ToLower(requestBody.Email))
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

	id := srv.DB.AddUser(
		requestBody.Name,
		requestBody.Email,
		string(encryptedPassword),
		sql.NullString{String: requestBody.FirstName, Valid: true},
		sql.NullString{String: requestBody.LastName, Valid: true},
		sql.NullString{String: requestBody.DoB, Valid: true},
		sql.NullString{String: requestBody.Gender, Valid: true},
	)

	err = revalidateURL(fmt.Sprintf("/user/%d", id))
	if err != nil {
		log.Printf("Error while revalidating `/user/%d`: %v", id, err)
	}
}

// loginHandler handles user login requests. It first checks if the user is
// already logged in by calling the `getUserId` method. If the user is already
// logged in, it returns an error message. If not, it decodes the request body
// to get the user's login credentials (username/email and password). It then
// checks if the user exists in the database and if the password is correct
// using the `GetUserByLogin` and `CompareHashAndPassword` methods respectively.
// If the user exists and the password is correct, it generates a new session
// token using the `uuid.NewV4()` method, adds the session to the database using
// the `AddSession` method, and sets a cookie with the session token using
// the `http.SetCookie` method. Finally, it returns a success response.
//
// Path: /api/login
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

// logoutHandler handles user logout requests. It first checks if the user has a
// valid session token by checking the "forum-token" cookie in the request. If the
// cookie is not present or invalid, it returns an unauthorized error response
// using the `errorResponse` function. If the cookie is valid, it removes the
// session from the database using the `RemoveSession` method and sets a new cookie
// with an empty value and an expired time using the `http.SetCookie` method. This
// effectively logs the user out by invalidating their session token.
//
// Path: /api/logout
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

// `func (srv *Server) getUserId(w http.ResponseWriter, r *http.Request) int` is a
// method of the `Server` struct that takes in a `http.ResponseWriter` and a
// `http.Request` as arguments and returns an integer.
//
// It first checks if the user has a valid session token by checking the "forum-token"
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
