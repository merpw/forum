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
		Author        SafeUser `json:"author"` // TODO: maybe remove
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
			Author:        SafeUser{Id: user.Id, Name: user.Name},
			Date:          post.Date,
			CommentsCount: post.CommentsCount,
			LikesCount:    post.LikesCount,
			DislikesCount: post.DislikesCount,
		})
	}

	sendObject(w, response)
}

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
		Author        SafeUser `json:"author"` // TODO: maybe remove
		Date          string   `json:"date"`
		CommentsCount int      `json:"comments_count"`
		LikesCount    int      `json:"likes_count"`
	}

	posts := srv.DB.GetUserPosts(userId)

	var response []ResponsePost
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
