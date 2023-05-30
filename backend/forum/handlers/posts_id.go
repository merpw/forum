package handlers

import (
	"backend/common/server"
	"net/http"
	"strconv"
	"strings"
)

// postsId returns a single post from the database that matches the incoming id of the post in the url
//
// Example: /api/posts/1
func (h *Handlers) postsId(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/posts/")
	// /api/posts/1 -> 1

	id, err := strconv.Atoi(idStr)
	if err != nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	// Get the post from the database
	post := h.DB.GetPostById(id)
	if post == nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	postAuthor := h.DB.GetUserById(post.AuthorId)
	safePost := SafePost{
		Id:            post.Id,
		Title:         post.Title,
		Content:       post.Content,
		Description:   post.Description,
		Author:        SafeUser{Id: postAuthor.Id, Name: postAuthor.Name},
		Date:          post.Date,
		CommentsCount: post.CommentsCount,
		LikesCount:    post.LikesCount,
		DislikesCount: post.DislikesCount,
		Categories:    post.Categories,
	}

	server.SendObject(w, safePost)
}

// postsIdLike likes a post in the database
func (h *Handlers) postsIdLike(w http.ResponseWriter, r *http.Request) {

	userId := h.getUserId(w, r)
	if userId == -1 {
		server.ErrorResponse(w, http.StatusUnauthorized)
		return
	}

	postIdStr := strings.TrimPrefix(r.URL.Path, "/api/posts/")
	postIdStr = strings.TrimSuffix(postIdStr, "/like")
	// /api/posts/1/like -> 1

	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		server.ErrorResponse(w, http.StatusNotFound)
	}

	post := h.DB.GetPostById(postId)
	if post == nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	reaction := h.DB.GetPostReaction(postId, userId)

	switch reaction {
	case 0: // if not reacted, add like
		h.DB.AddPostReaction(postId, userId, 1)
		h.DB.UpdatePostLikesCount(postId, +1)

		server.SendObject(w, +1)

	case 1: // if already liked, unlike
		h.DB.RemovePostReaction(postId, userId)
		h.DB.UpdatePostLikesCount(postId, -1)

		server.SendObject(w, 0)

	case -1: // if disliked, remove dislike and add like
		h.DB.RemovePostReaction(postId, userId)
		h.DB.UpdatePostDislikeCount(postId, -1)

		h.DB.AddPostReaction(postId, userId, 1)
		h.DB.UpdatePostLikesCount(postId, +1)

		server.SendObject(w, 1)
	}
}

// postsPostsIdDislikeHandler dislikes a post in the database
func (h *Handlers) postsIdDislike(w http.ResponseWriter, r *http.Request) {
	userId := h.getUserId(w, r)
	if userId == -1 {
		server.ErrorResponse(w, http.StatusUnauthorized)
		return
	}
	postIdStr := strings.TrimPrefix(r.URL.Path, "/api/posts/")
	postIdStr = strings.TrimSuffix(postIdStr, "/dislike")
	// /api/posts/1/dislike -> 1

	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		server.ErrorResponse(w, http.StatusNotFound)
	}

	post := h.DB.GetPostById(postId)
	if post == nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	reaction := h.DB.GetPostReaction(postId, userId)

	switch reaction {
	case 0: // if not reacted, add dislike
		h.DB.AddPostReaction(postId, userId, -1)
		h.DB.UpdatePostDislikeCount(postId, +1)

		server.SendObject(w, -1)

	case -1: // if already disliked, remove dislike
		h.DB.RemovePostReaction(postId, userId)
		h.DB.UpdatePostDislikeCount(postId, -1)

		server.SendObject(w, 0)

	case 1: // if liked, remove like and add dislike
		h.DB.RemovePostReaction(postId, userId)
		h.DB.UpdatePostLikesCount(postId, -1)

		h.DB.AddPostReaction(postId, userId, -1)
		h.DB.UpdatePostDislikeCount(postId, +1)

		server.SendObject(w, -1)
	}
}

func (h *Handlers) postsIdReaction(w http.ResponseWriter, r *http.Request) {
	userId := h.getUserId(w, r)
	if userId == -1 {
		server.ErrorResponse(w, http.StatusUnauthorized)
		return
	}

	postIdStr := strings.TrimPrefix(r.URL.Path, "/api/posts/")
	postIdStr = strings.TrimSuffix(postIdStr, "/reaction")
	// /api/posts/1/reaction -> 1

	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	post := h.DB.GetPostById(postId)
	if post == nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	reaction := h.DB.GetPostReaction(postId, userId)
	safeReaction := SafeReaction{
		Reaction:      reaction,
		LikesCount:    post.LikesCount,
		DislikesCount: post.DislikesCount,
	}

	if userId != post.AuthorId {
		server.SendObject(w, safeReaction)
		return
	}

	server.SendObject(w, struct {
		SafeReaction
		DislikesCount int `json:"dislikes_count"`
	}{
		SafeReaction:  safeReaction,
		DislikesCount: post.DislikesCount,
	})
}
