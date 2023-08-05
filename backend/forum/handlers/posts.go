package handlers

import (
	"backend/common/server"
	"net/http"
)

var categories = []string{"facts", "rumors", "other"}

// posts returns a json list of all posts from the database
func (h *Handlers) posts(w http.ResponseWriter, _ *http.Request) {
	posts := h.DB.GetAllPosts()

	response := make([]SafePost, 0)
	for _, post := range posts {
		postAuthor := h.DB.GetUserById(post.AuthorId)
		response = append(response, SafePost{
			Id:          post.Id,
			Title:       post.Title,
			Description: post.Description,
			Date:        post.Date,
			Author: SafeUser{
				Id:       postAuthor.Id,
				Username: postAuthor.Username,
				Avatar:   postAuthor.Avatar.String,
				Bio:      postAuthor.Bio.String,
			},
			CommentsCount: post.CommentsCount,
			LikesCount:    post.LikesCount,
			DislikesCount: post.DislikesCount,
			Categories:    post.Categories,
		})
	}

	server.SendObject(w, response)
}
