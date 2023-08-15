package handlers

import (
	"backend/common/server"
	. "backend/forum/database"
	"backend/forum/external"
	"database/sql"
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
	userId := r.Context().Value(userIdCtxKey).(int)

	requestBody := struct {
		Title       string   `json:"title"`
		Content     string   `json:"content"`
		Description string   `json:"description"`
		Categories  []string `json:"categories"`
		Privacy     int      `json:"privacy"`
		Audience    []int    `json:"audience"`
		GroupId     *int64   `json:"group_id"`
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

	var groupId sql.NullInt64

	if requestBody.GroupId == nil {
		groupId = sql.NullInt64{Valid: false}
	} else {
		if group := h.DB.GetGroupById(int(*requestBody.GroupId)); group == nil {
			server.ErrorResponse(w, http.StatusBadRequest)
			return
		}
		groupId = sql.NullInt64{Int64: *requestBody.GroupId, Valid: true}
		requestBody.Privacy = int(Public)
	}

	if requestBody.Privacy == int(SuperPrivate) && len(requestBody.Audience) == 0 {
		http.Error(w, "No audience specified for super private post", http.StatusBadRequest)
		return
	}

	if requestBody.Privacy < 0 || requestBody.Privacy > 2 {
		http.Error(w, "Privacy is not valid", http.StatusBadRequest)
		return
	}

	for _, follower := range requestBody.Audience {
		if requestBody.Privacy != int(SuperPrivate) {
			http.Error(w, "Unexpected audience, post is not super private", http.StatusBadRequest)
			return
		}
		if follower == userId {
			http.Error(w, "Can't add yourself as an audience", http.StatusBadRequest)
			return
		}

		user := h.DB.GetUserById(follower)
		if user == nil {
			http.Error(w, "One of the audience users does not exist", http.StatusBadRequest)
			return
		}

		if *h.DB.GetFollowStatus(follower, userId) != Accepted {
			http.Error(w, fmt.Sprintf("Audience is not valid, %s is not your follower", user.Username),
				http.StatusBadRequest)
			return
		}
	}

	id := h.DB.AddPost(
		requestBody.Title,
		requestBody.Content,
		requestBody.Description,
		userId,
		strings.Join(requestBody.Categories, ","),
		Privacy(requestBody.Privacy),
		groupId)

	if requestBody.Privacy == int(SuperPrivate) {
		for _, follower := range requestBody.Audience {
			h.DB.AddPostAudience(id, *h.DB.GetFollowId(follower, userId))
		}
	}

	server.SendObject(w, id)

	external.RevalidateURL(fmt.Sprintf("/post/%v", id))

	external.RevalidateURL("/")

	for _, category := range categories {
		external.RevalidateURL(fmt.Sprintf("/category/%v", category))
	}
}
