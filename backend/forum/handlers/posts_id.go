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
func (handlers *Handlers) postsId(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/posts/")
	// /api/posts/1 -> 1

	id, err := strconv.Atoi(idStr)
	if err != nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	// Get the post from the database
	post := handlers.DB.GetPostById(id)
	if post == nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	postAuthor := handlers.DB.GetUserById(post.AuthorId)
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
func (handlers *Handlers) postsIdLike(w http.ResponseWriter, r *http.Request) {

	userId := handlers.getUserId(w, r)
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

	post := handlers.DB.GetPostById(postId)
	if post == nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	reaction := handlers.DB.GetPostReaction(postId, userId)

	switch reaction {
	case 0: // if not reacted, add like
		handlers.DB.AddPostReaction(postId, userId, 1)
		handlers.DB.UpdatePostLikesCount(postId, +1)

		server.SendObject(w, +1)

	case 1: // if already liked, unlike
		handlers.DB.RemovePostReaction(postId, userId)
		handlers.DB.UpdatePostLikesCount(postId, -1)

		server.SendObject(w, 0)

	case -1: // if disliked, remove dislike and add like
		handlers.DB.RemovePostReaction(postId, userId)
		handlers.DB.UpdatePostDislikeCount(postId, -1)

		handlers.DB.AddPostReaction(postId, userId, 1)
		handlers.DB.UpdatePostLikesCount(postId, +1)

		server.SendObject(w, 1)
	}
}

// postsPostsIdDislikeHandler dislikes a post in the database
func (handlers *Handlers) postsIdDislike(w http.ResponseWriter, r *http.Request) {
	userId := handlers.getUserId(w, r)
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

	post := handlers.DB.GetPostById(postId)
	if post == nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	reaction := handlers.DB.GetPostReaction(postId, userId)

	switch reaction {
	case 0: // if not reacted, add dislike
		handlers.DB.AddPostReaction(postId, userId, -1)
		handlers.DB.UpdatePostDislikeCount(postId, +1)

		server.SendObject(w, -1)

	case -1: // if already disliked, remove dislike
		handlers.DB.RemovePostReaction(postId, userId)
		handlers.DB.UpdatePostDislikeCount(postId, -1)

		server.SendObject(w, 0)

	case 1: // if liked, remove like and add dislike
		handlers.DB.RemovePostReaction(postId, userId)
		handlers.DB.UpdatePostLikesCount(postId, -1)

		handlers.DB.AddPostReaction(postId, userId, -1)
		handlers.DB.UpdatePostDislikeCount(postId, +1)

		server.SendObject(w, -1)
	}
}

func (handlers *Handlers) postsIdReaction(w http.ResponseWriter, r *http.Request) {
	userId := handlers.getUserId(w, r)
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

	post := handlers.DB.GetPostById(postId)
	if post == nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	reaction := handlers.DB.GetPostReaction(postId, userId)
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