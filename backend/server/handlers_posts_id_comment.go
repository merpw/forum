package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// postsIdCommentIdLikeHandler likes a comment on a post in the database
func (srv *Server) postsIdCommentIdLikeHandler(w http.ResponseWriter, r *http.Request) {

	userId := srv.getUserId(w, r)
	if userId == -1 {
		errorResponse(w, http.StatusUnauthorized)
		return
	}

	commentId, err := strconv.Atoi(strings.Split(r.URL.Path, "/")[5])
	// /api/posts/1/comment/2/like ->2

	if err != nil {
		errorResponse(w, http.StatusNotFound)
	}

	comment := srv.DB.GetCommentById(commentId)
	if comment == nil {
		errorResponse(w, http.StatusNotFound)
		return
	}

	reaction := srv.DB.GetCommentReaction(commentId, userId)

	switch reaction {
	case 0: // if not reacted, add like
		srv.DB.AddCommentReaction(commentId, userId, 1)
		srv.DB.UpdateCommentLikesCount(commentId, +1)

		sendObject(w, +1)

	case 1: // if already liked, unlike
		srv.DB.RemoveCommentReaction(commentId, userId)
		srv.DB.UpdateCommentLikesCount(commentId, -1)

		sendObject(w, 0)

	case -1: // if disliked, remove dislike and add like
		srv.DB.RemoveCommentReaction(commentId, userId)
		srv.DB.UpdateCommentDislikeCount(commentId, -1)

		srv.DB.AddCommentReaction(commentId, userId, 1)
		srv.DB.UpdateCommentLikesCount(commentId, +1)

		sendObject(w, 1)
	}
}

// postsIdCommentIdDislikeHandler dislikes a comment on a post in the database
func (srv *Server) postsIdCommentIdDislikeHandler(w http.ResponseWriter, r *http.Request) {
	userId := srv.getUserId(w, r)
	if userId == -1 {
		errorResponse(w, http.StatusUnauthorized)
		return
	}

	commentId, err := strconv.Atoi(strings.Split(r.URL.Path, "/")[5])
	if err != nil {
		errorResponse(w, http.StatusNotFound)
	}

	comment := srv.DB.GetCommentById(commentId)
	if comment == nil {
		errorResponse(w, http.StatusNotFound)
		return
	}

	reaction := srv.DB.GetCommentReaction(commentId, userId)

	switch reaction {
	case 0: // if not reacted, add dislike
		srv.DB.AddCommentReaction(commentId, userId, -1)
		srv.DB.UpdateCommentDislikeCount(commentId, +1)

		sendObject(w, -1)

	case -1: // if already disliked, remove dislike
		srv.DB.RemoveCommentReaction(commentId, userId)
		srv.DB.UpdateCommentDislikeCount(commentId, -1)

		sendObject(w, 0)

	case 1: // if liked, remove like and add dislike
		srv.DB.RemoveCommentReaction(commentId, userId)
		srv.DB.UpdateCommentLikesCount(commentId, -1)

		srv.DB.AddCommentReaction(commentId, userId, -1)
		srv.DB.UpdateCommentDislikeCount(commentId, +1)

		sendObject(w, -1)
	}
}

func (srv *Server) postsIdCommentIdReactionHandler(w http.ResponseWriter, r *http.Request) {
	userId := srv.getUserId(w, r)
	if userId == -1 {
		errorResponse(w, http.StatusUnauthorized)
		return
	}

	commentId, err := strconv.Atoi(strings.Split(r.URL.Path, "/")[5])
	if err != nil {
		errorResponse(w, http.StatusNotFound)
		return
	}

	comment := srv.DB.GetCommentById(commentId)
	if comment == nil {
		errorResponse(w, http.StatusNotFound)
		return
	}

	reaction := srv.DB.GetCommentReaction(commentId, userId)
	if userId == comment.AuthorId {
		sendObject(w, struct {
			Reaction      int `json:"reaction"`
			LikesCount    int `json:"likes_count"`
			DislikesCount int `json:"dislikes_count"`
		}{
			Reaction:      reaction,
			LikesCount:    comment.LikesCount,
			DislikesCount: comment.DislikesCount,
		})
	} else {
		sendObject(w, struct {
			Reaction   int `json:"reaction"`
			LikesCount int `json:"likes_count"`
		}{
			Reaction:   reaction,
			LikesCount: comment.LikesCount,
		})
	}
}

// postsIdCommentCreateHandler Add comments on a post in the database
func (srv *Server) postsIdCommentCreateHandler(w http.ResponseWriter, r *http.Request) {
	userId := srv.getUserId(w, r)
	if userId == -1 {
		errorResponse(w, http.StatusUnauthorized)
		return
	}

	postIdStr := strings.TrimPrefix(r.URL.Path, "/api/posts/")
	postIdStr = strings.TrimSuffix(postIdStr, "/comment")
	// /api/posts/1/comment -> 1

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

	id := srv.DB.AddComment(requestBody.Content, postId, userId)
	srv.DB.UpdatePostsCommentsCount(postId, +1)

	sendObject(w, id)
}

// postsIdCommentsHandler returns all comments on a post in the database
func (srv *Server) postsIdCommentsHandler(w http.ResponseWriter, r *http.Request) {
	postId, err := strconv.Atoi(strings.Split(r.URL.Path, "/")[3])
	// /api/posts/1/comments -> 1
	
	if err != nil {
		errorResponse(w, http.StatusNotFound)
		return
	}

	post := srv.DB.GetPostById(postId)
	if post == nil {
		errorResponse(w, http.StatusNotFound)
		return
	}

	type ResponseComment struct {
		Id            int      `json:"id"`
		PostId        int      `json:"post_id"`
		AuthorId      int      `json:"author_id"`
		Content       string   `json:"content"`
		Author        SafeUser `json:"author"`
		Date          string   `json:"date"`
		LikesCount    int      `json:"likes_count"`
		DislikesCount int      `json:"dislikes_count"`
	}

	// posts := srv.DB.GetUserPosts(userId)
	comments := srv.DB.GetPostComments(postId)

	response := make([]ResponseComment, 0)
	for _, comment := range comments {
		user := srv.DB.GetUserById(comment.AuthorId)
		response = append(response, ResponseComment{
			Id:            comment.Id,
			PostId:        postId,
			AuthorId:      comment.AuthorId,
			Content:       comment.Content,
			Author:        SafeUser{user.Id, user.Name},
			Date:          comment.Date,
			LikesCount:    comment.LikesCount,
			DislikesCount: comment.DislikesCount,
		})
	}

	sendObject(w, response)
}
