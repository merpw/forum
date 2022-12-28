package server

import (
	"fmt"
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
	sendObject(w, "Me")
}

func (srv *Server) apiMePostsHandler(w http.ResponseWriter, r *http.Request) {
	sendObject(w, "Me posts")
}

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
	userId, _ := strconv.Atoi(strings.Split(r.URL.Path, "/")[3])
	// todo database stuff for "posts by user's Id" + Error handling during managing data
	sendObject(w, fmt.Sprintf("show posts by user with Id %v", userId))
}
