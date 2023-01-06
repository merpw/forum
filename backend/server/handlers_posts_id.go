package server

import (
	"net/http"
	"strconv"
	"strings"
)

// postsIdHandler returns a single post from the database that matches the incoming id of the post in the url
//
// Example: /api/posts/1
func (srv *Server) postsIdHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/posts/")
	// /api/posts/1 -> 1

	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse(w, http.StatusNotFound)
		return
	}

	// Get the post from the database
	post := srv.DB.GetPostById(id)
	if post == nil {
		errorResponse(w, http.StatusNotFound)
		return
	}

	postAuthor := srv.DB.GetUserById(post.AuthorId)
	safePost := SafePost{
		Id:            post.Id,
		Title:         post.Title,
		Content:       post.Content,
		Author:        SafeUser{Id: postAuthor.Id, Name: postAuthor.Name},
		Date:          post.Date,
		CommentsCount: post.CommentsCount,
		LikesCount:    post.LikesCount,
		Categories:    post.Categories,
	}

	sendObject(w, safePost)
}

// postsIdLikeHandler likes a post in the database
func (srv *Server) postsIdLikeHandler(w http.ResponseWriter, r *http.Request) {

	userId := srv.getUserId(w, r)
	if userId == -1 {
		errorResponse(w, http.StatusUnauthorized)
		return
	}

	postIdStr := strings.TrimPrefix(r.URL.Path, "/api/posts/")
	postIdStr = strings.TrimSuffix(postIdStr, "/like")
	// /api/posts/1/like -> 1

	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		errorResponse(w, http.StatusNotFound)
	}

	post := srv.DB.GetPostById(postId)
	if post == nil {
		errorResponse(w, http.StatusNotFound)
		return
	}

	reaction := srv.DB.GetPostReaction(postId, userId)

	switch reaction {
	case 0: // if not reacted, add like
		srv.DB.AddPostReaction(postId, userId, 1)
		srv.DB.UpdatePostLikesCount(postId, +1)

		sendObject(w, +1)

	case 1: // if already liked, unlike
		srv.DB.RemovePostReaction(postId, userId)
		srv.DB.UpdatePostLikesCount(postId, -1)

		sendObject(w, 0)

	case -1: // if disliked, remove dislike and add like
		srv.DB.RemovePostReaction(postId, userId)
		srv.DB.UpdatePostDislikeCount(postId, -1)

		srv.DB.AddPostReaction(postId, userId, 1)
		srv.DB.UpdatePostLikesCount(postId, +1)

		sendObject(w, 1)
	}
}

// postsPostsIdDislikeHandler dislikes a post in the database
func (srv *Server) postsIdDislikeHandler(w http.ResponseWriter, r *http.Request) {
	userId := srv.getUserId(w, r)
	if userId == -1 {
		errorResponse(w, http.StatusUnauthorized)
		return
	}
	postIdStr := strings.TrimPrefix(r.URL.Path, "/api/posts/")
	postIdStr = strings.TrimSuffix(postIdStr, "/dislike")
	// /api/posts/1/dislike -> 1

	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		errorResponse(w, http.StatusNotFound)
	}

	post := srv.DB.GetPostById(postId)
	if post == nil {
		errorResponse(w, http.StatusNotFound)
		return
	}

	reaction := srv.DB.GetPostReaction(postId, userId)

	switch reaction {
	case 0: // if not reacted, add dislike
		srv.DB.AddPostReaction(postId, userId, -1)
		srv.DB.UpdatePostDislikeCount(postId, +1)

		sendObject(w, -1)

	case -1: // if already disliked, remove dislike
		srv.DB.RemovePostReaction(postId, userId)
		srv.DB.UpdatePostDislikeCount(postId, -1)

		sendObject(w, 0)

	case 1: // if liked, remove like and add dislike
		srv.DB.RemovePostReaction(postId, userId)
		srv.DB.UpdatePostLikesCount(postId, -1)

		srv.DB.AddPostReaction(postId, userId, -1)
		srv.DB.UpdatePostDislikeCount(postId, +1)

		sendObject(w, -1)
	}
}

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
	safeReaction := SafeReaction{
		Reaction:   reaction,
		LikesCount: post.LikesCount,
	}

	if userId != post.AuthorId {
		sendObject(w, safeReaction)
		return
	}

	sendObject(w, struct {
		SafeReaction
		DislikesCount int `json:"dislikes_count"`
	}{
		SafeReaction:  safeReaction,
		DislikesCount: post.DislikesCount,
	})
}
