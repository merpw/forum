package handlers

import (
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

	requestBody.Username = strings.TrimSpace(requestBody.Username)

	if IdUsernameRegex.MatchString(requestBody.Username) {
		http.Error(w, "Username format is not allowed", http.StatusBadRequest)
		return
	}

	if len(requestBody.Username) == 0 {
		requestBody.Username = "u" + strconv.Itoa(h.DB.GetLastUserId()+1)
		goto skipLengthCheck
	}

	if len(requestBody.Username) < MinUsernameLength {
		http.Error(w, "Username is too short", http.StatusBadRequest)
		return
	}

skipLengthCheck:

	if len(requestBody.Username) > MaxUsernameLength {
		http.Error(w, "Username is too long", http.StatusBadRequest)
		return
	}

	if !AvatarRegex.MatchString(requestBody.Avatar) {
		http.Error(w, "Avatar file string is not valid", http.StatusBadRequest)
		return
	}

	requestBody.Bio = strings.TrimSpace(requestBody.Bio)

	if len(requestBody.Bio) == 0 {
		requestBody.Bio = sql.NullString{String: "", Valid: false}.String
	}

	if len(requestBody.Bio) > MaxBioLength {
		http.Error(w, "Bio is too long", http.StatusBadRequest)
		return
	}

	if !UsernameRegex.MatchString(requestBody.Username) {
		http.Error(w, "Username is not valid, only letters, numbers and underscores are allowed",
			http.StatusBadRequest)
		return
	}

	if h.DB.IsNameTaken(requestBody.Username) {
		http.Error(w, "Username is already in use", http.StatusBadRequest)
		return
	}

	requestBody.FirstName = strings.TrimSpace(requestBody.FirstName)

	if len(requestBody.FirstName) < MinFirstNameLength {
		http.Error(w, "First name is not valid", http.StatusBadRequest)
		return
	}

	if len(requestBody.FirstName) > MaxFirstNameLength {
		http.Error(w, "First name is too long", http.StatusBadRequest)
		return
	}

	requestBody.LastName = strings.TrimSpace(requestBody.LastName)
	if len(requestBody.LastName) < MinLastNameLength {
		http.Error(w, "Last name is not valid", http.StatusBadRequest)
		return
	}

	if len(requestBody.LastName) > MaxLastNameLength {
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

	if !isPresent(Genders, requestBody.Gender) {
		http.Error(w, "Gender is not valid", http.StatusBadRequest)
		return
	}

	requestBody.Email = strings.TrimSpace(strings.ToLower(requestBody.Email))
	_, err = mail.ParseAddress(requestBody.Email)

	if err != nil || requestBody.Email != strings.TrimSpace(requestBody.Email) {
		http.Error(w, "Email is not valid", http.StatusBadRequest)
		return
	}

	if h.DB.IsEmailTaken(requestBody.Email) {
		http.Error(w, "Email is already taken", http.StatusBadRequest)
		return
	}

	if len(requestBody.Password) < MinPasswordLength {
		http.Error(w, "Password is too short", http.StatusBadRequest)
		return
	}

	if len(requestBody.Password) > MaxPasswordLength {
		http.Error(w, "Password is too long", http.StatusBadRequest)
		return
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), 10)
	if err != nil {
		http.Error(w, "Password is not valid", http.StatusBadRequest)
		return
	}

	id := h.DB.AddUser(
		requestBody.Username,
		requestBody.Email,
		string(encryptedPassword),
		sql.NullString{String: requestBody.FirstName, Valid: true},
		sql.NullString{String: requestBody.LastName, Valid: true},
		sql.NullString{String: requestBody.DoB, Valid: true},
		sql.NullString{String: requestBody.Gender, Valid: true},
		sql.NullString{String: requestBody.Bio, Valid: true},
		sql.NullString{String: requestBody.Avatar, Valid: true},
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
		h.revokeSession <- cookie.Value
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
			h.revokeSession <- cookie.Value
		}()
		return -1
	}

	return userId
}
