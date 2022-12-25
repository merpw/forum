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
	commands := strings.Split(strings.TrimPrefix(url, "/api/posts/"), "/")
	postId, err := strconv.Atoi(commands[0])
	lens := len(commands)

	if err != nil {
		switch lens {
		case 1:
			if commands[0] != "create" {
				errorResponse(w, http.StatusBadRequest)
				return
			}
		default:
			errorResponse(w, http.StatusBadRequest)
			return
		}
	}

	switch lens {
	case 1:
		if err != nil {
			srv.createPostHandler(w, r)
		} else {
			srv.showPostHandler(w, r, postId)
		}
	case 2:
		switch commands[1] {
		case "like":
			srv.likePostHandler(w, r, postId)
		case "dislike":
			srv.dislikePostHandler(w, r, postId)
		case "create":
			srv.createCommentHandler(w, r, postId)
		default:
			errorResponse(w, http.StatusBadRequest)
		}
	case 3:
		comment_id := strings.Split(commands[1], "_")
		if len(comment_id) != 2 {
			errorResponse(w, http.StatusBadRequest)
		}
		cid, err := strconv.Atoi(comment_id[1])
		if err != nil || comment_id[0] != "comment" {
			errorResponse(w, http.StatusBadRequest)
		}

		switch commands[2] {
		case "like":
			srv.likeCommentHandler(w, r, postId, cid)
		case "dislike":
			srv.dislikeCommentHandler(w, r, postId, cid)
		default:
			errorResponse(w, http.StatusBadRequest)
		}
	default:
		errorResponse(w, http.StatusBadRequest)
	}
}

func (srv *Server) createCommentHandler(w http.ResponseWriter, r *http.Request, postId int) {
	sendObject(w, fmt.Sprintf("create post %v comment", postId))
}

func (srv *Server) likeCommentHandler(w http.ResponseWriter, r *http.Request, postId, commentId int) {
	sendObject(w, fmt.Sprintf("like post %v comment %v", postId, commentId))
}

func (srv *Server) dislikeCommentHandler(w http.ResponseWriter, r *http.Request, postId int, commentId int) {
	sendObject(w, fmt.Sprintf("dislike post %v comment %v", postId, commentId))
}

// post
func (srv *Server) createPostHandler(w http.ResponseWriter, r *http.Request) {
	sendObject(w, "create post")
}

func (srv *Server) likePostHandler(w http.ResponseWriter, r *http.Request, postId int) {
	sendObject(w, fmt.Sprintf("like post %v", postId))
}

func (srv *Server) dislikePostHandler(w http.ResponseWriter, r *http.Request, postId int) {
	sendObject(w, fmt.Sprintf("dislike post %v", postId))
}

func (srv *Server) showPostHandler(w http.ResponseWriter, r *http.Request, postId int) {
	sendObject(w, fmt.Sprintf("show post %v", postId))
}
