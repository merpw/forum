package server

import (
	"net/http"
	"strconv"
	"strings"
)

func (srv *Server) apiUserMasterHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case reApiUserId.MatchString(r.URL.Path):
		srv.apiUserIdHandler(w, r)
	case reApiUserIdPosts.MatchString(r.URL.Path):
		srv.apiUserIdPostsHandler(w, r)
	}
}

// apiMeHandler returns current user information (id, name, email) if the current user is logged in.
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
		Id    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}

	sendObject(w, response)
}

// apiMePostsHandler returns the current user's posts to the current user if he is logged in.
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

	type ResponsePost struct {
		Id            int      `json:"id"`
		Title         string   `json:"title"`
		Content       string   `json:"content"`
		Author        SafeUser `json:"author"`
		Date          string   `json:"date"`
		CommentsCount int      `json:"comments_count"`
		LikesCount    int      `json:"likes_count"`
		DislikesCount int      `json:"dislikes_count"`
	}

	response := make([]ResponsePost, 0)
	for _, post := range posts {
		response = append(response, ResponsePost{
			Id:            post.Id,
			Title:         post.Title,
			Content:       post.Content,
			Author:        SafeUser{Id: user.Id, Name: user.Name},
			Date:          post.Date,
			CommentsCount: post.CommentsCount,
			LikesCount:    post.LikesCount,
			DislikesCount: post.DislikesCount,
		})
	}

	sendObject(w, response)
}

// apiMePostsHandler Returns the info of the user with the given id. The requester does not need to be logged in.
//
//	GET /api/user/:id
func (srv *Server) apiUserIdHandler(w http.ResponseWriter, r *http.Request) {
	userIdStr := strings.TrimPrefix(r.URL.Path, "/api/user/")
	// /api/user/1 -> 1

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

	userResponse := struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}{
		Id:   user.Id,
		Name: user.Name,
	}

	sendObject(w, userResponse)
}

// apiMePostsHandler Returns the posts of the user with the given id. The requester does not need to be logged in.
//
//	GET /api/user/:id/posts
func (srv *Server) apiUserIdPostsHandler(w http.ResponseWriter, r *http.Request) {
	userIdStr := strings.TrimPrefix(r.URL.Path, "/api/user/")
	userIdStr = strings.TrimSuffix(userIdStr, "/posts")
	// /api/user/1/posts -> 1

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

	type ResponsePost struct {
		Id            int      `json:"id"`
		Title         string   `json:"title"`
		Content       string   `json:"content"`
		Author        SafeUser `json:"author"`
		Date          string   `json:"date"`
		CommentsCount int      `json:"comments_count"`
		LikesCount    int      `json:"likes_count"`
	}

	posts := srv.DB.GetUserPosts(userId)

	response := make([]ResponsePost, 0)
	for _, post := range posts {
		response = append(response, ResponsePost{
			Id:            post.Id,
			Title:         post.Title,
			Content:       post.Content,
			Author:        SafeUser{user.Id, user.Name},
			Date:          post.Date,
			CommentsCount: post.CommentsCount,
			LikesCount:    post.LikesCount,
		})
	}

	sendObject(w, response)
}
