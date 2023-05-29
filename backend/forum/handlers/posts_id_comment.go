package handlers

import (
	"backend/common/server"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// postsIdCommentIdLike likes a comment on a post in the database
func (handlers *Handlers) postsIdCommentIdLike(w http.ResponseWriter, r *http.Request) {

	userId := handlers.getUserId(w, r)
	if userId == -1 {
		server.ErrorResponse(w, http.StatusUnauthorized)
		return
	}

	commentId, err := strconv.Atoi(strings.Split(r.URL.Path, "/")[5])
	// /api/posts/1/comment/2/like ->2

	if err != nil {
		server.ErrorResponse(w, http.StatusNotFound)
	}

	comment := handlers.DB.GetCommentById(commentId)
	if comment == nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	reaction := handlers.DB.GetCommentReaction(commentId, userId)

	switch reaction {
	case 0: // if not reacted, add like
		handlers.DB.AddCommentReaction(commentId, userId, 1)
		handlers.DB.UpdateCommentLikesCount(commentId, +1)

		server.SendObject(w, +1)

	case 1: // if already liked, unlike
		handlers.DB.RemoveCommentReaction(commentId, userId)
		handlers.DB.UpdateCommentLikesCount(commentId, -1)

		server.SendObject(w, 0)

	case -1: // if disliked, remove dislike and add like
		handlers.DB.RemoveCommentReaction(commentId, userId)
		handlers.DB.UpdateCommentDislikeCount(commentId, -1)

		handlers.DB.AddCommentReaction(commentId, userId, 1)
		handlers.DB.UpdateCommentLikesCount(commentId, +1)

		server.SendObject(w, 1)
	}
}

// postsIdCommentIdDislike dislikes a comment on a post in the database
func (handlers *Handlers) postsIdCommentIdDislike(w http.ResponseWriter, r *http.Request) {
	userId := handlers.getUserId(w, r)
	if userId == -1 {
		server.ErrorResponse(w, http.StatusUnauthorized)
		return
	}

	commentId, err := strconv.Atoi(strings.Split(r.URL.Path, "/")[5])
	if err != nil {
		server.ErrorResponse(w, http.StatusNotFound)
	}

	comment := handlers.DB.GetCommentById(commentId)
	if comment == nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	reaction := handlers.DB.GetCommentReaction(commentId, userId)

	switch reaction {
	case 0: // if not reacted, add dislike
		handlers.DB.AddCommentReaction(commentId, userId, -1)
		handlers.DB.UpdateCommentDislikeCount(commentId, +1)

		server.SendObject(w, -1)

	case -1: // if already disliked, remove dislike
		handlers.DB.RemoveCommentReaction(commentId, userId)
		handlers.DB.UpdateCommentDislikeCount(commentId, -1)

		server.SendObject(w, 0)

	case 1: // if liked, remove like and add dislike
		handlers.DB.RemoveCommentReaction(commentId, userId)
		handlers.DB.UpdateCommentLikesCount(commentId, -1)

		handlers.DB.AddCommentReaction(commentId, userId, -1)
		handlers.DB.UpdateCommentDislikeCount(commentId, +1)

		server.SendObject(w, -1)
	}
}

func (handlers *Handlers) postsIdCommentIdReaction(w http.ResponseWriter, r *http.Request) {
	userId := handlers.getUserId(w, r)
	if userId == -1 {
		server.ErrorResponse(w, http.StatusUnauthorized)
		return
	}

	commentId, err := strconv.Atoi(strings.Split(r.URL.Path, "/")[5])
	if err != nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	comment := handlers.DB.GetCommentById(commentId)
	if comment == nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	reaction := handlers.DB.GetCommentReaction(commentId, userId)
	server.SendObject(w,
		SafeReaction{
			Reaction:      reaction,
			LikesCount:    comment.LikesCount,
			DislikesCount: comment.DislikesCount,
		})
}

// postsIdCommentCreate Add comments on a post in the database
func (handlers *Handlers) postsIdCommentCreate(w http.ResponseWriter, r *http.Request) {
	userId := handlers.getUserId(w, r)
	if userId == -1 {
		server.ErrorResponse(w, http.StatusUnauthorized)
		return
	}

	postIdStr := strings.TrimPrefix(r.URL.Path, "/api/posts/")
	postIdStr = strings.TrimSuffix(postIdStr, "/comment")
	// /api/posts/1/comment -> 1

	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}
	post := handlers.DB.GetPostById(postId)
	if post == nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	requestBody := struct {
		Content string `json:"content"`
	}{}

	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Body is not valid", http.StatusBadRequest)
		return
	}

	requestBody.Content = strings.TrimSpace(requestBody.Content)

	if len(requestBody.Content) < 1 {
		http.Error(w, "Content is too short", http.StatusBadRequest)
		return
	}

	id := handlers.DB.AddComment(requestBody.Content, postId, userId)
	handlers.DB.UpdatePostsCommentsCount(postId, +1)

	server.SendObject(w, id)
}

// postsIdComments returns all comments on a post in the database
func (handlers *Handlers) postsIdComments(w http.ResponseWriter, r *http.Request) {
	postId, err := strconv.Atoi(strings.Split(r.URL.Path, "/")[3])
	// /api/posts/1/comments -> 1

	if err != nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	post := handlers.DB.GetPostById(postId)
	if post == nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	// posts := srv.DB.GetUserPosts(userId)
	comments := handlers.DB.GetPostComments(postId)

	response := make([]SafeComment, 0)
	for _, comment := range comments {
		user := handlers.DB.GetUserById(comment.AuthorId)
		response = append(response, SafeComment{
			Id:            comment.Id,
			Content:       comment.Content,
			Author:        SafeUser{user.Id, user.Name},
			Date:          comment.Date,
			LikesCount:    comment.LikesCount,
			DislikesCount: comment.DislikesCount,
		})
	}

	server.SendObject(w, response)
}
