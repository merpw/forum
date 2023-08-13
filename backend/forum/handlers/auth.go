package handlers

import (
	"backend/common/integrations/auth"
	"backend/common/server"
	"backend/forum/external"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/mail"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	MinUsernameLength = 3
	MaxUsernameLength = 15

	MinFirstNameLength = 1
	MaxFirstNameLength = MaxUsernameLength

	MinLastNameLength = 1
	MaxLastNameLength = MaxUsernameLength

	MinPasswordLength = 8
	MaxPasswordLength = 128

	MinEmailLength = 3
	MaxEmailLength = 254

	MinBioLength = 0
	MaxBioLength = 200
)

var (
	UsernameRegex   = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	IdUsernameRegex = regexp.MustCompile(`^u\d+$`)
	AvatarRegex     = regexp.MustCompile(`^[0-9]\.jpg$`)

	MinLoginLength = int(math.Min(MinUsernameLength, MinEmailLength))
	MaxLoginLength = int(math.Max(MaxUsernameLength, MaxEmailLength))
)

var Genders = []string{"male", "female", "other"}

func (h *Handlers) signup(w http.ResponseWriter, r *http.Request) {
	if h.getUserId(w, r) != -1 {
		http.Error(w, "You are already logged in", http.StatusBadRequest)
		return
	}

	requestBody := struct {
		Username  string `json:"username"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		DoB       string `json:"dob"`
		Gender    string `json:"gender"`
		Bio       string `json:"bio"`
		Avatar    string `json:"avatar"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Body is not valid", http.StatusBadRequest)
		return
	}

	username := strings.TrimSpace(requestBody.Username)
	if len(username) == 0 {
		username = "u" + strconv.Itoa(h.DB.GetLastUserId()+1)
		goto skipUserNameCheck
	}
	if IdUsernameRegex.MatchString(username) {
		http.Error(w, "invalid username format", http.StatusBadRequest)
		return
	}
	if len(username) < MinUsernameLength {
		http.Error(w, "username is too short", http.StatusBadRequest)
		return
	}
	if len(username) > MaxUsernameLength {
		http.Error(w, "username is too long", http.StatusBadRequest)
		return
	}
	if !UsernameRegex.MatchString(username) {
		http.Error(w, "username is not valid, only letters, numbers and underscores are allowed", http.StatusBadRequest)
		return
	}
	if h.DB.IsNameTaken(username) {
		http.Error(w, "username is already in use", http.StatusBadRequest)
		return
	}

skipUserNameCheck:

	firstName := strings.TrimSpace(requestBody.FirstName)
	if len(firstName) < MinFirstNameLength {
		http.Error(w, "first name is too short", http.StatusBadRequest)
		return
	}
	if len(firstName) > MaxFirstNameLength {
		http.Error(w, "first name is too long", http.StatusBadRequest)
		return
	}

	lastName := strings.TrimSpace(requestBody.LastName)
	if len(lastName) < MinLastNameLength {
		http.Error(w, "last name is too short", http.StatusBadRequest)
		return
	}
	if len(lastName) > MaxLastNameLength {
		http.Error(w, "last name is too long", http.StatusBadRequest)
		return
	}

	if requestBody.DoB == "" {
		http.Error(w, "date of birth is invalid", http.StatusBadRequest)
		return
	}

	dob, dobErr := time.Parse("2006-01-02", requestBody.DoB)
	if dobErr != nil {
		http.Error(w, "date of birth is invalid", http.StatusBadRequest)
		return
	}

	now := time.Now()
	mindate := time.Date(1900, time.January, 1, 0, 0, 0, 0, time.UTC)

	if dob.After(now) || dob.Before(mindate) {
		http.Error(w, "date of birth is invalid", http.StatusBadRequest)
		return
	}

	email := strings.TrimSpace(strings.ToLower(requestBody.Email))
	_, err = mail.ParseAddress(email)
	if err != nil || email != strings.TrimSpace(email) {
		http.Error(w, "email is invalid", http.StatusBadRequest)
		return
	}
	if h.DB.IsEmailTaken(email) {
		http.Error(w, "email is already taken", http.StatusBadRequest)
		return
	}

	if !isPresent(Genders, requestBody.Gender) {
		http.Error(w, "Gender is not valid", http.StatusBadRequest)
		return
	}

	if len(requestBody.Password) < MinPasswordLength {
		http.Error(w, "password is too short", http.StatusBadRequest)
		return
	}

	if len(requestBody.Password) > MaxPasswordLength {
		http.Error(w, "password is too long", http.StatusBadRequest)
		return
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), 10)
	if err != nil {
		http.Error(w, "invalid password", http.StatusBadRequest)
		return
	}

	var avatar sql.NullString
	requestBody.Avatar = strings.TrimSpace(requestBody.Avatar)

	if !AvatarRegex.MatchString(requestBody.Avatar) && len(requestBody.Avatar) != 0 {
		http.Error(w, "avatar file is invalid", http.StatusBadRequest)
		return
	}

	if len(requestBody.Avatar) == 0 {
		avatar = sql.NullString{String: "", Valid: false}
	} else {
		avatar = sql.NullString{String: requestBody.Avatar, Valid: true}
	}

	var bio sql.NullString
	requestBody.Bio = strings.TrimSpace(requestBody.Bio)

	if len(requestBody.Bio) > MaxBioLength {
		http.Error(w, "bio is too long", http.StatusBadRequest)
	}

	if len(requestBody.Bio) == MinBioLength {
		bio = sql.NullString{String: "", Valid: false}
	} else {
		bio = sql.NullString{String: requestBody.Bio, Valid: true}
	}

	id := h.DB.AddUser(
		username,
		email,
		string(encryptedPassword),
		sql.NullString{String: firstName, Valid: true},
		sql.NullString{String: lastName, Valid: true},
		sql.NullString{String: requestBody.DoB, Valid: true},
		sql.NullString{String: requestBody.Gender, Valid: true},
		avatar,
		bio,
	)

	external.RevalidateURL(fmt.Sprintf("/user/%d", id))
}

func (h *Handlers) login(w http.ResponseWriter, r *http.Request) {
	if h.getUserId(w, r) != -1 {
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

	if len(requestBody.Login) < MinLoginLength && !IdUsernameRegex.MatchString(requestBody.Login) ||
		len(requestBody.Login) > MaxLoginLength ||
		len(requestBody.Password) < MinPasswordLength || len(requestBody.Password) > MaxPasswordLength {
		http.Error(w, "Invalid login or password", http.StatusBadRequest)
		return
	}

	user := h.DB.GetUserByLogin(strings.ToLower(requestBody.Login))
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

	h.DB.AddSession(token.String(), int(expire.Unix()), user.Id)

	http.SetCookie(w, &http.Cookie{
		Name:    "forum-token",
		Value:   token.String(),
		Expires: expire,
		Path:    "/",
		//	TODO: add secure on production
	})
}

func (h *Handlers) logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("forum-token")
	if err != nil {
		server.ErrorResponse(w, http.StatusUnauthorized)
		return
	}
	h.DB.RemoveSession(cookie.Value)

	http.SetCookie(w, &http.Cookie{
		Name:    "forum-token",
		Value:   "",
		Expires: time.Now(),
		Path:    "/",
	})
	go func() {
		event := auth.Event{
			Type: auth.EventTypeTokenRevoked,
			Item: cookie.Value,
		}
		h.event <- event
	}()
}

// getUserId checks if the user's token is valid and returns the user's id
//
// If the token is not valid, it will return -1 and add an empty cookie to the response
func (h *Handlers) getUserId(w http.ResponseWriter, r *http.Request) int {
	cookie, err := r.Cookie("forum-token")
	if err != nil {
		return -1
	}
	userId := h.DB.CheckSession(cookie.Value)

	if userId == -1 {
		http.SetCookie(w, &http.Cookie{
			Name:    "forum-token",
			Value:   "",
			Expires: time.Now(),
			Path:    "/",
		})
		go func() {
			event := auth.Event{
				Type: auth.EventTypeTokenRevoked,
				Item: cookie.Value,
			}
			h.event <- event
		}()
		return -1
	}

	return userId
}
