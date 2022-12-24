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

	// /api/posts/id ->                    [id] ... show one post with comments
	// /api/posts/create ->                [create] ... create post
	// /api/posts/id/like ->               [id, like] ... like post
	// /api/posts/id/dislike ->            [id, dislike] ... dislike post
	// /api/posts/id/create ->             [id, create] ... create comment for the post
	// /api/posts/id/comment_id/like ->    [id, comment_id, like] ... like comment of the post
	// /api/posts/id/comment_id/dislike -> [id, comment_id, dislike] ... dislike comment of the post
	var commands = strings.Split(strings.TrimPrefix(url, "/api/posts/"), "/")
	var postId int
	var err error
	postId, err = strconv.Atoi(commands[0])

	if len(commands) == 1 { // /api/posts/id
		srv.postHandler(w, r, postId)
		return
	}

	switch len(commands) {
	case 1:
		switch commands[0] {
		case "create":
			srv.createHandler(w, r)
		default:
			if err != nil {
				errorResponse(w, 400)
				return
			}
			srv.postHandler(w, r, postId)
			return
		}
	case 2:

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
