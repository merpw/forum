package handlers

import (
	"backend/common/server"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	MinCommentLength = 1
	MaxCommentLength = 1000
)

// postsIdCommentIdLike likes a comment on a specific post
func (h *Handlers) postsIdCommentIdLike(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userIdCtxKey).(int)

	postId := r.Context().Value(postIdCtxKey).(int)

	post := h.validatedPost(r, userId, postId)

	if post == nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	commentId, err := strconv.Atoi(strings.Split(r.URL.Path, "/")[5])
	// /api/posts/1/comment/2/like ->2

	if err != nil {
		server.ErrorResponse(w, http.StatusNotFound)
	}

	comment := h.DB.GetCommentById(commentId)
	if comment == nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	reaction := h.DB.GetCommentReaction(commentId, userId)

	switch reaction {
	case 0: // if not reacted, add like
		h.DB.AddCommentReaction(commentId, userId, 1)
		h.DB.UpdateCommentLikesCount(commentId, +1)

		server.SendObject(w, +1)

	case 1: // if already liked, unlike
		h.DB.RemoveCommentReaction(commentId, userId)
		h.DB.UpdateCommentLikesCount(commentId, -1)

		server.SendObject(w, 0)

	case -1: // if disliked, remove dislike and add like
		h.DB.RemoveCommentReaction(commentId, userId)
		h.DB.UpdateCommentDislikeCount(commentId, -1)

		h.DB.AddCommentReaction(commentId, userId, 1)
		h.DB.UpdateCommentLikesCount(commentId, +1)

		server.SendObject(w, 1)
	}
}

// postsIdCommentIdDislike dislikes a specific comment on a specific post
func (h *Handlers) postsIdCommentIdDislike(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userIdCtxKey).(int)

	postId := r.Context().Value(postIdCtxKey).(int)

	post := h.validatedPost(r, userId, postId)

	if post == nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	commentId, err := strconv.Atoi(strings.Split(r.URL.Path, "/")[5])
	if err != nil {
		server.ErrorResponse(w, http.StatusNotFound)
	}

	comment := h.DB.GetCommentById(commentId)
	if comment == nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	reaction := h.DB.GetCommentReaction(commentId, userId)

	switch reaction {
	case 0: // if not reacted, add dislike
		h.DB.AddCommentReaction(commentId, userId, -1)
		h.DB.UpdateCommentDislikeCount(commentId, +1)

		server.SendObject(w, -1)

	case -1: // if already disliked, remove dislike
		h.DB.RemoveCommentReaction(commentId, userId)
		h.DB.UpdateCommentDislikeCount(commentId, -1)

		server.SendObject(w, 0)

	case 1: // if liked, remove like and add dislike
		h.DB.RemoveCommentReaction(commentId, userId)
		h.DB.UpdateCommentLikesCount(commentId, -1)

		h.DB.AddCommentReaction(commentId, userId, -1)
		h.DB.UpdateCommentDislikeCount(commentId, +1)

		server.SendObject(w, -1)
	}
}

// postsIdCommentIdReactions returns SafeReaction with data about the reactions of a specific comment
func (h *Handlers) postsIdCommentIdReaction(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userIdCtxKey).(int)

	postId := r.Context().Value(postIdCtxKey).(int)

	post := h.validatedPost(r, userId, postId)

	if post == nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	commentId, err := strconv.Atoi(strings.Split(r.URL.Path, "/")[5])
	if err != nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	comment := h.DB.GetCommentById(commentId)
	if comment == nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	reaction := h.DB.GetCommentReaction(commentId, userId)
	server.SendObject(w,
		SafeReaction{
			Reaction:      reaction,
			LikesCount:    comment.LikesCount,
			DislikesCount: comment.DislikesCount,
		})
}

// postsIdCommentCreate adds a comment on a specific post
func (h *Handlers) postsIdCommentCreate(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userIdCtxKey).(int)

	postId := r.Context().Value(postIdCtxKey).(int)

	post := h.validatedPost(r, userId, postId)

	if post == nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	requestBody := struct {
		Content string `json:"content"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Body is not valid", http.StatusBadRequest)
		return
	}

	requestBody.Content = strings.TrimSpace(requestBody.Content)

	if len(requestBody.Content) < MinCommentLength {
		http.Error(w, "Content is too short", http.StatusBadRequest)
		return
	}

	if len(requestBody.Content) > MaxCommentLength {
		http.Error(w,
			fmt.Sprintf("Content is too long (max %d characters)", MaxContentLength),
			http.StatusBadRequest)
		return
	}

	id := h.DB.AddComment(requestBody.Content, postId, userId)
	h.DB.UpdatePostsCommentsCount(postId, +1)

	server.SendObject(w, id)
}

// postsIdComments returns all comments on a specific post
func (h *Handlers) postsIdComments(w http.ResponseWriter, r *http.Request) {
	postId := r.Context().Value(postIdCtxKey).(int)

	post := h.validatedPost(r, h.getUserId(w, r), postId)

	if post == nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	// posts := srv.DB.GetUserPosts(usersId)
	comments := h.DB.GetPostComments(postId)

	response := make([]SafeComment, 0)
	for _, comment := range comments {
		user := h.DB.GetUserById(comment.AuthorId)
		response = append(response, SafeComment{
			Id:      comment.Id,
			Content: comment.Content,
			Author: SafeUser{
				Id:       user.Id,
				Username: user.Username,
				Avatar:   user.Avatar.String,
				Bio:      user.Bio.String,
			},
			Date:          comment.Date,
			LikesCount:    comment.LikesCount,
			DislikesCount: comment.DislikesCount,
		})
	}

	server.SendObject(w, response)
}
