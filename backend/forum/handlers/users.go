package handlers

import (
	"backend/common/server"
	"net/http"
	"strconv"
	"strings"
)

// usersAll returns all userIds in alphabetical order.
//
// GET /api/users
func (h *Handlers) usersAll(w http.ResponseWriter, r *http.Request) {
	userIds := h.DB.GetAllUserIds()

	server.SendObject(w, userIds)
}

// usersId returns the info of the user with the given id
//
//	GET /api/users/:id
func (h *Handlers) usersId(w http.ResponseWriter, r *http.Request) {
	userIdStr := strings.TrimPrefix(r.URL.Path, "/api/users/")
	// /api/users/1 -> 1

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	user := h.DB.GetUserById(userId)
	if user == nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	server.SendObject(w, SafeUser{
		Id:       user.Id,
		Username: user.Username,
		Avatar:   user.Avatar.String,
		Bio:      user.Bio.String,
	})
}

// usersIdPosts Returns the posts of the user with the given id.
//
//	GET /api/users/:id/posts
func (h *Handlers) usersIdPosts(w http.ResponseWriter, r *http.Request) {
	userIdStr := strings.TrimPrefix(r.URL.Path, "/api/users/")
	userIdStr = strings.TrimSuffix(userIdStr, "/posts")
	// /api/users/1/posts -> 1

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	user := h.DB.GetUserById(userId)
	if user == nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	posts := h.DB.GetUserPosts(userId)

	response := make([]SafePost, 0)
	for _, post := range posts {
		response = append(response, SafePost{
			Id:          post.Id,
			Title:       post.Title,
			Description: post.Description,
			Author: SafeUser{
				Id:       user.Id,
				Username: user.Username,
				Avatar:   user.Avatar.String,
				Bio:      user.Bio.String,
			},
			Date:          post.Date,
			CommentsCount: post.CommentsCount,
			LikesCount:    post.LikesCount,
			DislikesCount: post.DislikesCount,
			Categories:    post.Categories,
		})
	}

	server.SendObject(w, response)
}
