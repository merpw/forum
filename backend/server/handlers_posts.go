package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// postsHandler handles /api/posts/... and calls appropriate handler
func (srv *Server) postsHandler(w http.ResponseWriter, r *http.Request) {

	// trim trailing slash
	url := strings.TrimRight(r.URL.String(), "/")
	if url == "/api/posts" {
		sendObject(w, srv.posts)
		return
	}

	// /api/posts/id -> [id], /api/posts/id/comments -> [id, comments]
	commands := strings.Split(strings.TrimPrefix(url, "/api/posts/"), "/")

	postId, err := strconv.Atoi(commands[0])
	if err != nil {
		errorResponse(w, 404)
		return
	}

	if len(commands) == 1 { // /api/posts/id
		srv.postHandler(w, r, postId)
		return
	}

	switch commands[1] {
	case "comment":
		srv.commentHandler(w, r, postId)
	case "like":
		srv.likeHandler(w, r, postId)
	case "dislike":
		srv.dislikeHandler(w, r, postId)
	default:
		errorResponse(w, 404)
	}
}

func (srv *Server) commentHandler(w http.ResponseWriter, r *http.Request, postId int) {
	sendObject(w, fmt.Sprintf("comment post %v", postId))
}

func (srv *Server) likeHandler(w http.ResponseWriter, r *http.Request, postId int) {
	sendObject(w, fmt.Sprintf("like post %v", postId))
}

func (srv *Server) dislikeHandler(w http.ResponseWriter, r *http.Request, postId int) {
	sendObject(w, fmt.Sprintf("dislike post %v", postId))
}

func (srv *Server) postHandler(w http.ResponseWriter, r *http.Request, postId int) {
	sendObject(w, fmt.Sprintf("post %v", postId))
}

func (srv *Server) createHandler(w http.ResponseWriter, r *http.Request) {
	sendObject(w, "create")
}
