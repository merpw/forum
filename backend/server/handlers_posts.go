package server

import (
	"fmt"
	"log"
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

	// /api/posts/create ->                  [create] ... create post
	// /api/posts/id ->                      [id] ... show one post with comments
	// /api/posts/id/like ->                 [id, like] ... like post
	// /api/posts/id/dislike ->              [id, dislike] ... dislike post
	// /api/posts/id/comments ->             [id, comments] ... show all comments of the post
	// /api/posts/id/create ->               [id, create] ... create comment for the post
	// /api/posts/category/{facts|rumors|created|liked} -> [category, facts|rumors|created|liked] ... show all posts with chosen category. To satisfy task requirements
	// /api/posts/id/comment_id/like ->      [id, comment_id, like] ... like comment of the post
	// /api/posts/id/comment_id/dislike ->   [id, comment_id, dislike] ... dislike comment of the post
	commands := strings.Split(strings.TrimPrefix(url, "/api/posts/"), "/")
	postId, err := strconv.Atoi(commands[0])
	lenOfCommands := len(commands)

	if err != nil {
		switch lenOfCommands {
		case 1:
			if commands[0] != "create" {
				errorResponse(w, http.StatusBadRequest)
				return
			}
		case 2:
			switch commands[1] {
			case "facts", "rumors", "created", "liked":
				switch commands[0] {
				case "category":
				default:
					errorResponse(w, http.StatusBadRequest)
					return
				}
			default:
				errorResponse(w, http.StatusBadRequest)
				return
			}
		default:
			errorResponse(w, http.StatusBadRequest)
			return
		}
	}

	switch lenOfCommands {
	case 1:
		if err != nil {
			srv.createPostHandler(w, r, url)
		} else {
			srv.showPostHandler(w, r, postId, url)
		}
	case 2:
		switch commands[1] {
		case "like":
			srv.likePostHandler(w, r, postId, url)
		case "dislike":
			srv.dislikePostHandler(w, r, postId, url)
		case "create":
			srv.createCommentHandler(w, r, postId, url)
		case "comments":
			srv.showCommentsHandler(w, r, postId, url)
		case "facts", "rumors", "created", "liked":
			srv.categoryPostHandler(w, r, commands[1], url)
		default:
			errorResponse(w, http.StatusBadRequest)
		}
	case 3:
		commentId, err := strconv.Atoi(commands[1])
		if err != nil {
			errorResponse(w, http.StatusBadRequest)
		}

		switch commands[2] {
		case "like":
			srv.likeCommentHandler(w, r, postId, commentId, url)
		case "dislike":
			srv.dislikeCommentHandler(w, r, postId, commentId, url)
		default:
			errorResponse(w, http.StatusBadRequest)
		}
	default:
		errorResponse(w, http.StatusBadRequest)
	}
}

// comment
func (srv *Server) showCommentsHandler(w http.ResponseWriter, r *http.Request, postId int, pathToCheck string) {

	// todo database stuff for "show comments" + Error handling during managing data

	sendObject(w, fmt.Sprintf("show post %v", postId))
}

func (srv *Server) createCommentHandler(w http.ResponseWriter, r *http.Request, postId int, pathToCheck string) {

	// todo database stuff for "create comment" + Error handling during managing data

	sendObject(w, fmt.Sprintf("create post %v comment", postId))
}

func (srv *Server) likeCommentHandler(w http.ResponseWriter, r *http.Request, postId, commentId int, pathToCheck string) {

	// todo database stuff for "like comment" + Error handling during managing data

	sendObject(w, fmt.Sprintf("like post %v comment %v", postId, commentId))
}

func (srv *Server) dislikeCommentHandler(w http.ResponseWriter, r *http.Request, postId int, commentId int, pathToCheck string) {

	// todo database stuff for "dislike comment" + Error handling during managing data

	sendObject(w, fmt.Sprintf("dislike post %v comment %v", postId, commentId))
}

// post
func (srv *Server) createPostHandler(w http.ResponseWriter, r *http.Request, pathToCheck string) {
	cookie, err := r.Cookie("session")
	if err != nil {
		log.Println(err)
		return
	}
	// todo database stuff for "create post" + Error handling during managing data

	sendObject(w, "create post, token: "+cookie.Value)
}

func (srv *Server) likePostHandler(w http.ResponseWriter, r *http.Request, postId int, pathToCheck string) {

	// todo database stuff for "like post" + Error handling during managing data

	sendObject(w, fmt.Sprintf("like post %v", postId))
}

func (srv *Server) dislikePostHandler(w http.ResponseWriter, r *http.Request, postId int, pathToCheck string) {

	// todo database stuff for "dislike post" + Error handling during managing data

	sendObject(w, fmt.Sprintf("dislike post %v", postId))
}

func (srv *Server) showPostHandler(w http.ResponseWriter, r *http.Request, postId int, pathToCheck string) {

	// todo database stuff for "show post" + Error handling during managing data

	sendObject(w, fmt.Sprintf("show post %v", postId))
}

// plan four categories only "facts" and "rumors", also "created" and "liked" posts, should be enough
func (srv *Server) categoryPostHandler(w http.ResponseWriter, r *http.Request, category, pathToCheck string) {

	// todo database stuff for "category post" + Error handling during managing data

	sendObject(w, fmt.Sprintf("show %v posts", category))
}
