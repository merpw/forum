package server

import (
	"net/http"
	"strconv"
	"strings"
)

func (srv *Server) postsIdReactionHandler(w http.ResponseWriter, r *http.Request) {
	userId := srv.getUserId(w, r)
	if userId == -1 {
		errorResponse(w, http.StatusUnauthorized)
		return
	}

	postIdStr := strings.TrimPrefix(r.URL.Path, "/api/posts/")
	postIdStr = strings.TrimSuffix(postIdStr, "/reaction")
	// /api/posts/1/reaction -> 1

	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		errorResponse(w, http.StatusNotFound)
		return
	}

	post := srv.DB.GetPostById(postId)
	if post == nil {
		errorResponse(w, http.StatusNotFound)
		return
	}

	reaction := srv.DB.GetPostReaction(postId, userId)
	if userId == post.AuthorId {
		sendObject(w, struct {
			Reaction      int `json:"reaction"`
			LikesCount    int `json:"likes_count"`
			DislikesCount int `json:"dislikes_count"`
		}{
			Reaction:      reaction,
			LikesCount:    post.LikesCount,
			DislikesCount: post.DislikesCount,
		})
	} else {
		sendObject(w, struct {
			Reaction   int `json:"reaction"`
			LikesCount int `json:"likes_count"`
		}{
			Reaction:   reaction,
			LikesCount: post.LikesCount,
		})
	}
}
