package handlers

import (
	"backend/common/server"
	"net/http"
	"strconv"
	"strings"
)

// userAll returns all userIds in alphabetical order.
//
// GET /api/users
func (h *Handlers) userAll(w http.ResponseWriter, r *http.Request) {
	users := h.DB.GetAllUsers()

	userIds := []int{}
	for _, user := range users {
		userIds = append(userIds, user.Id)
	}

	server.SendObject(w, userIds)
}

// userId Returns the info of the user with the given id. The requester does not need to be logged in.
//
//	GET /api/users/:id
func (h *Handlers) userId(w http.ResponseWriter, r *http.Request) {
	userIdStr := strings.TrimPrefix(r.URL.Path, "/api/users/")
	// /api/user/1 -> 1

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

	server.SendObject(w, SafeUser{Id: user.Id, Name: user.Name})
}

// userIdPosts Returns the posts of the user with the given id.
// The requester does not need to be logged in.
//
//	GET /api/users/:id/posts
func (h *Handlers) userIdPosts(w http.ResponseWriter, r *http.Request) {
	userIdStr := strings.TrimPrefix(r.URL.Path, "/api/users/")
	userIdStr = strings.TrimSuffix(userIdStr, "/posts")
	// /api/user/1/posts -> 1

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
			Id:            post.Id,
			Title:         post.Title,
			Description:   post.Description,
			Author:        SafeUser{user.Id, user.Name},
			Date:          post.Date,
			CommentsCount: post.CommentsCount,
			LikesCount:    post.LikesCount,
			DislikesCount: post.DislikesCount,
			Categories:    post.Categories,
		})
	}

	server.SendObject(w, response)
}
