package handlers

import (
	"net/http"
)

var categories = []string{"facts", "rumors", "other"}

// posts returns a json list of all posts from the database
func (handlers *Handlers) posts(w http.ResponseWriter, _ *http.Request) {
	posts := handlers.DB.GetAllPosts()

	response := make([]SafePost, 0)
	for _, post := range posts {
		postAuthor := handlers.DB.GetUserById(post.AuthorId)
		response = append(response, SafePost{
			Id:            post.Id,
			Title:         post.Title,
			Description:   post.Description,
			Date:          post.Date,
			Author:        SafeUser{Id: postAuthor.Id, Name: postAuthor.Name},
			CommentsCount: post.CommentsCount,
			LikesCount:    post.LikesCount,
			DislikesCount: post.DislikesCount,
			Categories:    post.Categories,
		})
	}

	sendObject(w, response)
}
