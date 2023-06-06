package handlers

import (
	"backend/common/server"
	"backend/forum/external"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// postsCreate creates a new post in the database
const (
	MinTitleLength = 1
	MaxTitleLength = 100

	MinContentLength = 1
	MaxContentLength = 10000

	MinDescriptionLength = 1
	MaxDescriptionLength = 200
)

// postsCreate creates a new post
func (h *Handlers) postsCreate(w http.ResponseWriter, r *http.Request) {
	userId := h.getUserId(w, r)
	if userId == -1 {
		server.ErrorResponse(w, http.StatusUnauthorized)
		return
	}

	requestBody := struct {
		Title       string   `json:"title"`
		Content     string   `json:"content"`
		Description string   `json:"description"`
		Categories  []string `json:"categories"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Body is not valid", http.StatusBadRequest)
		return
	}

	requestBody.Title = strings.TrimSpace(requestBody.Title)
	requestBody.Content = strings.TrimSpace(requestBody.Content)
	requestBody.Description = strings.TrimSpace(requestBody.Description)

	if len(requestBody.Title) < MinTitleLength {
		http.Error(w, "Title is too short", http.StatusBadRequest)
		return
	}
	if len(requestBody.Title) > MaxTitleLength {
		http.Error(w, fmt.Sprintf("Title is too long, maximum length is %v", MaxTitleLength),
			http.StatusBadRequest)
		return
	}

	if len(requestBody.Content) < MinContentLength {
		http.Error(w, "Content is too short", http.StatusBadRequest)
		return
	}
	if len(requestBody.Content) > MaxContentLength {
		http.Error(w, fmt.Sprintf("Content is too long, maximum length is %v", MaxContentLength),
			http.StatusBadRequest)
		return
	}

	if len(requestBody.Description) < MinDescriptionLength {
		if len(requestBody.Content) < MaxDescriptionLength {
			requestBody.Description = requestBody.Content
		} else {
			requestBody.Description = requestBody.Content[:MaxDescriptionLength]
		}
	}

	if len(requestBody.Description) > MaxDescriptionLength {
		requestBody.Description = requestBody.Description[:MaxDescriptionLength]
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
	if len(requestBody.Categories) == 0 {
		isValid = false
	}

	if !isValid {
		http.Error(w, "Categories are not valid", http.StatusBadRequest)
		return
	}

	id := h.DB.AddPost(requestBody.Title, requestBody.Content, requestBody.Description,
		userId, strings.Join(requestBody.Categories, ","))
	server.SendObject(w, id)

	external.RevalidateURL(fmt.Sprintf("/post/%v", id))

	external.RevalidateURL("/")

	for _, category := range categories {
		external.RevalidateURL(fmt.Sprintf("/category/%v", category))
	}
}
