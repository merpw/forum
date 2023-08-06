package handlers

import (
	"backend/common/server"
	"backend/forum/external"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"regexp"
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

	username, usernameError := h.checkUserName(requestBody.Username)
	if usernameError != nil {
		http.Error(w, usernameError.Error(), http.StatusBadRequest)
		return
	}

	firstName, firstNameError := h.checkFirstName(requestBody.FirstName)
	if firstNameError != nil {
		http.Error(w, firstNameError.Error(), http.StatusBadRequest)
		return
	}

	lastName, lastNameError := h.checkLastName(requestBody.LastName)
	if lastNameError != nil {
		http.Error(w, lastNameError.Error(), http.StatusBadRequest)
		return
	}

	dob, dobError := h.checkDoB(requestBody.DoB)
	if dobError != nil {
		http.Error(w, dobError.Error(), http.StatusBadRequest)
		return
	}

	email, emailError := h.checkEmail(requestBody.Email)
	if emailError != nil {
		http.Error(w, emailError.Error(), http.StatusBadRequest)
		return
	}

	if !isPresent(Genders, requestBody.Gender) {
		http.Error(w, "Gender is not valid", http.StatusBadRequest)
		return
	}

	password, passwordError := h.checkPassword(requestBody.Password)
	if passwordError != nil {
		http.Error(w, passwordError.Error(), http.StatusBadRequest)
		return
	}

	avatar, avatarError := h.checkAvatar(requestBody.Avatar)
	if avatarError != nil {
		http.Error(w, avatarError.Error(), http.StatusBadRequest)
		return
	}

	bio, bioError := h.checkBio(requestBody.Bio)
	if bioError != nil {
		http.Error(w, bioError.Error(), http.StatusBadRequest)
		return
	}

	id := h.DB.AddUser(
		username,
		email,
		password,
		sql.NullString{String: firstName, Valid: true},
		sql.NullString{String: lastName, Valid: true},
		sql.NullString{String: dob, Valid: true},
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
