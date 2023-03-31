package server

import (
	"encoding/json"
	"net/http"
	"strings"
)

// postsCreateHandler creates a new post in the database
func (srv *Server) postsCreateHandler(w http.ResponseWriter, r *http.Request) {
	userId := srv.GetUserId(w, r)
	if userId == -1 {
		errorResponse(w, http.StatusUnauthorized)
		return
	}

	requestBody := struct {
		Title      string   `json:"title"`
		Content    string   `json:"content"`
		Categories []string `json:"categories"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Body is not valid", http.StatusBadRequest)
		return
	}

	requestBody.Title = strings.TrimSpace(requestBody.Title)
	requestBody.Content = strings.TrimSpace(requestBody.Content)

	if len(requestBody.Title) < 1 {
		http.Error(w, "Title is too short", http.StatusBadRequest)
		return
	}
	if len(requestBody.Content) < 1 {
		http.Error(w, "Content is too short", http.StatusBadRequest)
		return
	}

	if len(requestBody.Title) > 25 {
		http.Error(w, "Title is too long, maximum length is 25", http.StatusBadRequest)
		return
	}

	for i, cat := range requestBody.Categories {
		cat = strings.TrimSpace(cat)
		cat = strings.ToLower(cat)
		requestBody.Categories[i] = cat
	}

	isValid := true
	for _, cat := range requestBody.Categories {
		if !isPresent(categories, cat) {
			isValid = false
			break
		}
	}

	if !isValid {
		http.Error(w, "Categories are not valid", http.StatusBadRequest)
		return
	}

	id := srv.DB.AddPost(requestBody.Title, requestBody.Content, userId, strings.Join(requestBody.Categories, ","))
	sendObject(w, id)
}

func isPresent(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
