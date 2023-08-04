package handlers

import (
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
	"strconv"
	"strings"
	"time"
)

const (
	PUBLIC = iota
	PRIVATE
)

type SafeUser struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar,omitempty"`
	Bio      string `json:"bio,omitempty"`
}

type SafePost struct {
	Id            int      `json:"id"`
	Title         string   `json:"title"`
	Content       string   `json:"content"`
	Description   string   `json:"description"`
	Author        SafeUser `json:"author"`
	Date          string   `json:"date"`
	CommentsCount int      `json:"comments_count"`
	LikesCount    int      `json:"likes_count"`
	DislikesCount int      `json:"dislikes_count"`
	Categories    string   `json:"categories"`
}

type SafeComment struct {
	Id            int      `json:"id"`
	Content       string   `json:"content"`
	Author        SafeUser `json:"author"`
	Date          string   `json:"date"`
	LikesCount    int      `json:"likes_count"`
	DislikesCount int      `json:"dislikes_count"`
}

type SafeReaction struct {
	Reaction      int `json:"reaction"`
	LikesCount    int `json:"likes_count"`
	DislikesCount int `json:"dislikes_count"`
}

func (h *Handlers) checkUserName(username string) (string, error) {
	username = strings.TrimSpace(username)
	if len(username) == 0 {
		username = "u" + strconv.Itoa(h.DB.GetLastUserId()+1)
		return username, nil
	}
	if IdUsernameRegex.MatchString(username) {
		return "", errors.New("invalid username format")
	}
	if len(username) < MinUsernameLength {
		return "", errors.New("username is too short")
	}
	if len(username) > MaxUsernameLength {
		return "", errors.New("username is too long")
	}
	if !UsernameRegex.MatchString(username) {
		return "", errors.New("username is not valid, only letters, numbers and underscores are allowed")
	}
	if h.DB.IsNameTaken(username) {
		return "", errors.New("username is too long")
	}
	return username, nil
}

func (h *Handlers) checkFirstName(firstName string) (string, error) {
	firstName = strings.TrimSpace(firstName)
	if len(firstName) < MinFirstNameLength {
		return "", errors.New("first name is not valid")
	}
	if len(firstName) > MaxFirstNameLength {
		return "", errors.New("first name is too long")
	}
	return firstName, nil
}

func (h *Handlers) checkLastName(lastName string) (string, error) {
	lastName = strings.TrimSpace(lastName)
	if len(lastName) < MinLastNameLength {
		return "", errors.New("last name is not valid")
	}
	if len(lastName) > MaxLastNameLength {
		return "", errors.New("last name is too long")
	}
	return lastName, nil
}

func (h *Handlers) checkDoB(date string) (string, error) {
	if date == "" {
		return "", errors.New("date of birth is not valid")
	}

	dob, err := time.Parse("2006-01-02", date)
	if err != nil {
		return "", errors.New("date of birth is not valid")
	}

	now := time.Now()
	mindate := time.Date(1900, time.January, 1, 0, 0, 0, 0, time.UTC)

	if dob.After(now) || dob.Before(mindate) {
		return "", errors.New("date of birth is not valid")
	}

	return date, nil
}

func (h *Handlers) checkEmail(email string) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	_, err := mail.ParseAddress(email)
	if err != nil || email != strings.TrimSpace(email) {
		return "", errors.New("email is not valid")
	}
	if h.DB.IsEmailTaken(email) {
		return "", errors.New("email is already taken")
	}
	return email, nil
}

func (h *Handlers) checkPassword(password string) (string, error) {
	if len(password) < MinPasswordLength {
		return "", errors.New("password is too short")
	}

	if len(password) > MaxPasswordLength {
		return "", errors.New("password is too long")
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", errors.New("password is invalid")
	}

	return string(encryptedPassword), nil
}

func (h *Handlers) checkAvatar(avatar string) (sql.NullString, error) {
	avatar = strings.TrimSpace(avatar)
	if len(avatar) == 0 {
		return sql.NullString{String: "", Valid: false}, nil
	}

	if !AvatarRegex.MatchString(avatar) {
		return sql.NullString{String: "", Valid: false}, errors.New("avatar file string is invalid")
	}

	return sql.NullString{String: avatar, Valid: true}, nil
}

func (h *Handlers) checkBio(bio string) (sql.NullString, error) {
	bio = strings.TrimSpace(bio)
	if len(bio) == MinBioLength {
		return sql.NullString{String: "", Valid: false}, nil
	}

	if len(bio) > MaxBioLength {
		return sql.NullString{String: "", Valid: false}, errors.New("bio is too long")
	}

	return sql.NullString{String: bio, Valid: true}, nil
}

func isPresent(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
