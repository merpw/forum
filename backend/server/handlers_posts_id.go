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
	type SafeUser struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}

	type Response struct {
		Id            int      `json:"id"`
		Title         string   `json:"title"`
		Content       string   `json:"content"`
		Author        SafeUser `json:"author"`
		Date          string   `json:"date"`
		CommentsCount int      `json:"comments_count"`
		Comments      []struct {
			Id            int      `json:"id"`
			Content       string   `json:"content"`
			Author        SafeUser `json:"author"`
			Date          string   `json:"date"`
			LikesCount    int      `json:"likes_count"`
			DislikesCount int      `json:"dislikes_count"`
		} `json:"comments"`
		LikesCount int    `json:"likes_count"`
		Category   string `json:"category"`
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/posts/")
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

	response := Response{
		Id:            post.Id,
		Title:         post.Title,
		Content:       post.Content,
		Author:        SafeUser{Id: postAuthor.Id, Name: postAuthor.Name},
		Date:          post.Date,
		CommentsCount: post.CommentsCount,
		LikesCount:    post.LikesCount,
		Category:      post.Category,
	}

	comments := srv.DB.GetPostComments(post.Id)

	response.Comments = make([]struct {
		Id            int      `json:"id"`
		Content       string   `json:"content"`
		Author        SafeUser `json:"author"`
		Date          string   `json:"date"`
		LikesCount    int      `json:"likes_count"`
		DislikesCount int      `json:"dislikes_count"`
	}, len(comments))
	for i, comment := range comments {
		commentAuthor := srv.DB.GetUserById(comment.AuthorId)
		response.Comments[i] = struct {
			Id            int      `json:"id"`
			Content       string   `json:"content"`
			Author        SafeUser `json:"author"`
			Date          string   `json:"date"`
			LikesCount    int      `json:"likes_count"`
			DislikesCount int      `json:"dislikes_count"`
		}{
			Id:            comment.Id,
			Content:       comment.Content,
			Author:        SafeUser{Id: commentAuthor.Id, Name: commentAuthor.Name},
			Date:          comment.Date,
			LikesCount:    comment.LikesCount,
			DislikesCount: comment.DislikesCount,
		}
	}
	// reverse response.Comments
	for i, j := 0, len(response.Comments)-1; i < j; i, j = i+1, j-1 {
		response.Comments[i], response.Comments[j] = response.Comments[j], response.Comments[i]
	}

	sendObject(w, response)
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
