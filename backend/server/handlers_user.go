package server

import (
	"net/http"
	"strconv"
	"strings"
)

// apiUserIdHandler handles paths /api/me and /api/me/posts, passing them to the appropriate handler.
func (srv *Server) apiUserMasterHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case reApiUserId.MatchString(r.URL.Path):
		srv.apiUserIdHandler(w, r)
	case reApiUserIdPosts.MatchString(r.URL.Path):
		srv.apiUserIdPostsHandler(w, r)
	}
}

// apiMeHandler returns the currently logged in user's information.
//
// (id, name, email, first name, last name, date of birth, gender)
//
//	GET /api/me
func (srv *Server) apiMeHandler(w http.ResponseWriter, r *http.Request) {
	userId := srv.getUserId(w, r)
	if userId == -1 {
		errorResponse(w, http.StatusUnauthorized)
		return
	}

	user := srv.DB.GetUserById(userId)

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

// apiMePostsHandler returns the posts of the currently logged in user.
//
//	GET /api/me/posts
func (srv *Server) apiMePostsHandler(w http.ResponseWriter, r *http.Request) {
	userId := srv.getUserId(w, r)
	if userId == -1 {
		errorResponse(w, http.StatusUnauthorized)
		return
	}
	user := srv.DB.GetUserById(userId)
	posts := srv.DB.GetUserPosts(userId)

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

// apiMePostsLikedHandler
// returns liked posts of the currently logged in user.
//
// Route: /api/me/posts/liked
func (srv *Server) apiMePostsLikedHandler(w http.ResponseWriter, r *http.Request) {
	userId := srv.getUserId(w, r)
	if userId == -1 {
		errorResponse(w, http.StatusUnauthorized)
		return
	}
	posts := srv.DB.GetUserPostsLiked(userId)

	type Response struct {
		SafePost
	}

	response := make([]Response, 0)
	for _, post := range posts {
		author := srv.DB.GetUserById(post.AuthorId)
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

// apiMePostsHandler Returns the info of the user with the given id. The requester does not need to be logged in.
//
//	GET /api/user/:id
func (srv *Server) apiUserIdHandler(w http.ResponseWriter, r *http.Request) {
	userIdStr := strings.TrimPrefix(r.URL.Path, "/api/user/")

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		errorResponse(w, http.StatusNotFound)
		return
	}

	user := srv.DB.GetUserById(userId)
	if user == nil {
		errorResponse(w, http.StatusNotFound)
		return
	}

	sendObject(w, SafeUser{Id: user.Id, Name: user.Name})
}

// apiMePostsHandler Returns the posts of the user with the given id.
// The requester does not need to be logged in.
//
//	GET /api/user/:id/posts
func (srv *Server) apiUserIdPostsHandler(w http.ResponseWriter, r *http.Request) {
	userIdStr := strings.TrimPrefix(r.URL.Path, "/api/user/")
	userIdStr = strings.TrimSuffix(userIdStr, "/posts")

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		errorResponse(w, http.StatusNotFound)
		return
	}

	user := srv.DB.GetUserById(userId)
	if user == nil {
		errorResponse(w, http.StatusNotFound)
		return
	}

	posts := srv.DB.GetUserPosts(userId)

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
