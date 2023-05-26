package handlers

import (
	"net/http"
	"strconv"
	"strings"
)

// me returns the currently logged in user's information.
//
// (id, name, email, first name, last name, date of birth, gender)
//
//	GET /api/me
func (handlers *Handlers) me(w http.ResponseWriter, r *http.Request) {
	userId := handlers.getUserId(w, r)
	if userId == -1 {
		errorResponse(w, http.StatusUnauthorized)
		return
	}

	user := handlers.DB.GetUserById(userId)

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

	sendObject(w, response)
}

// mePosts returns the posts of the currently logged in user.
//
//	GET /api/me/posts
func (handlers *Handlers) mePosts(w http.ResponseWriter, r *http.Request) {
	userId := handlers.getUserId(w, r)
	if userId == -1 {
		errorResponse(w, http.StatusUnauthorized)
		return
	}
	user := handlers.DB.GetUserById(userId)
	posts := handlers.DB.GetUserPosts(userId)

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

	sendObject(w, response)
}

// mePostsLiked
// returns liked posts of the currently logged in user.
//
//	GET /api/me/posts/liked
func (handlers *Handlers) mePostsLiked(w http.ResponseWriter, r *http.Request) {
	userId := handlers.getUserId(w, r)
	if userId == -1 {
		errorResponse(w, http.StatusUnauthorized)
		return
	}
	posts := handlers.DB.GetUserPostsLiked(userId)

	type Response struct {
		SafePost
		DislikesCount int `json:"dislikesCount"`
	}

	response := make([]Response, 0)
	for _, post := range posts {
		author := handlers.DB.GetUserById(post.AuthorId)
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

	sendObject(w, response)
}

// mePosts Returns the info of the user with the given id. The requester does not need to be logged in.
//
//	GET /api/user/:id
func (handlers *Handlers) userId(w http.ResponseWriter, r *http.Request) {
	userIdStr := strings.TrimPrefix(r.URL.Path, "/api/user/")
	// /api/user/1 -> 1

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		errorResponse(w, http.StatusNotFound)
		return
	}

	user := handlers.DB.GetUserById(userId)
	if user == nil {
		errorResponse(w, http.StatusNotFound)
		return
	}

	sendObject(w, SafeUser{Id: user.Id, Name: user.Name})
}

// mePosts Returns the posts of the user with the given id.
// The requester does not need to be logged in.
//
//	GET /api/user/:id/posts
func (handlers *Handlers) userIdPosts(w http.ResponseWriter, r *http.Request) {
	userIdStr := strings.TrimPrefix(r.URL.Path, "/api/user/")
	userIdStr = strings.TrimSuffix(userIdStr, "/posts")
	// /api/user/1/posts -> 1

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		errorResponse(w, http.StatusNotFound)
		return
	}

	user := handlers.DB.GetUserById(userId)
	if user == nil {
		errorResponse(w, http.StatusNotFound)
		return
	}

	posts := handlers.DB.GetUserPosts(userId)

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

	sendObject(w, response)
}
