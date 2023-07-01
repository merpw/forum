package handlers

import (
	"backend/common/server"
	"net/http"
	"strconv"
	"strings"
)

// postsId returns the post with the given id
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
		Author:        SafeUser{Id: postAuthor.Id, Username: postAuthor.Username},
		Date:          post.Date,
		CommentsCount: post.CommentsCount,
		LikesCount:    post.LikesCount,
		DislikesCount: post.DislikesCount,
		Categories:    post.Categories,
	}

	server.SendObject(w, safePost)
}

// postsIdLike handles the like of the post with the given id
//
// returns current reaction
func (h *Handlers) postsIdLike(w http.ResponseWriter, r *http.Request) {

	userId := r.Context().Value(userIdCtxKey).(int)

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

// postsIdDislike handles the dislike of a post.
//
// returns current reaction
func (h *Handlers) postsIdDislike(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userIdCtxKey).(int)
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

// postsIdReaction returns the current reaction of the user to the post
func (h *Handlers) postsIdReaction(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userIdCtxKey).(int)

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
