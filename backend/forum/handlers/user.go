package handlers

import (
	"backend/common/server"
	"net/http"
	"strconv"
	"strings"
)

// me returns the currently logged in user's information.
//
// (id, name, email, first name, last name, date of birth, gender)
//
//	GET /api/me
func (h *Handlers) me(w http.ResponseWriter, r *http.Request) {
	userId := h.getUserId(w, r)
	if userId == -1 {
		server.ErrorResponse(w, http.StatusUnauthorized)
		return
	}

	user := h.DB.GetUserById(userId)

	response := struct {
		SafeUser
		Email     string `json:"email"`
		FirstName string `json:"first_name,omitempty"`
		LastName  string `json:"last_name,omitempty"`
		DoB       string `json:"dob,omitempty"`
		Gender    string `json:"gender,omitempty"`
	}{
		SafeUser:  SafeUser{Id: user.Id, Name: user.Name},
		Email:     user.Email,
		FirstName: user.FirstName.String,
		LastName:  user.LastName.String,
		DoB:       user.DoB.String,
		Gender:    user.Gender.String,
	}

	server.SendObject(w, response)
}

// mePosts returns the posts of the currently logged in user.
//
//	GET /api/me/posts
func (h *Handlers) mePosts(w http.ResponseWriter, r *http.Request) {
	userId := h.getUserId(w, r)
	if userId == -1 {
		server.ErrorResponse(w, http.StatusUnauthorized)
		return
	}
	user := h.DB.GetUserById(userId)
	posts := h.DB.GetUserPosts(userId)

	type Response struct {
		SafePost
	}

	response := make([]Response, 0)
	for _, post := range posts {
		response = append(response, Response{
			SafePost: SafePost{
				Id:            post.Id,
				Title:         post.Title,
				Description:   post.Description,
				Author:        SafeUser{Id: user.Id, Name: user.Name},
				Date:          post.Date,
				CommentsCount: post.CommentsCount,
				LikesCount:    post.LikesCount,
				DislikesCount: post.DislikesCount,
				Categories:    post.Categories,
			},
		})
	}

	server.SendObject(w, response)
}

// mePostsLiked
// returns liked posts of the currently logged in user.
//
//	GET /api/me/posts/liked
func (h *Handlers) mePostsLiked(w http.ResponseWriter, r *http.Request) {
	userId := h.getUserId(w, r)
	if userId == -1 {
		server.ErrorResponse(w, http.StatusUnauthorized)
		return
	}
	posts := h.DB.GetUserPostsLiked(userId)

	type Response struct {
		SafePost
		DislikesCount int `json:"dislikesCount"`
	}

	response := make([]Response, 0)
	for _, post := range posts {
		author := h.DB.GetUserById(post.AuthorId)
		response = append(response, Response{
			SafePost: SafePost{
				Id:            post.Id,
				Title:         post.Title,
				Description:   post.Description,
				Author:        SafeUser{Id: author.Id, Name: author.Name},
				Date:          post.Date,
				CommentsCount: post.CommentsCount,
				LikesCount:    post.LikesCount,
				DislikesCount: post.DislikesCount,
				Categories:    post.Categories,
			},
		})
	}

	server.SendObject(w, response)
}

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
