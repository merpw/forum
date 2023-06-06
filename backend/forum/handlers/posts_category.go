package handlers

import (
	"backend/common/server"
	"net/http"
	"strings"
)

// postsCategories returns a json list of all categories from the database
func (h *Handlers) postsCategories(w http.ResponseWriter, _ *http.Request) {
	server.SendObject(w, categories)
}

// postsCategoriesName returns a json list of all posts from the database
func (h *Handlers) postsCategoriesName(w http.ResponseWriter, r *http.Request) {
	categoryName := strings.TrimPrefix(r.URL.Path, "/api/posts/categories/")
	// /api/posts/categories/name -> name

	categoryName = strings.ToLower(categoryName)
	// Name -> name

	isValid := false
	for _, cat := range categories {
		if cat == categoryName {
			isValid = true
			break
		}
	}

	if !isValid {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	posts := h.DB.GetCategoryPosts(categoryName)

	response := make([]SafePost, 0)
	for _, post := range posts {
		postAuthor := h.DB.GetUserById(post.AuthorId)
		response = append(response, SafePost{
			Id:            post.Id,
			Title:         post.Title,
			Description:   post.Description,
			Date:          post.Date,
			Author:        SafeUser{Id: postAuthor.Id, Username: postAuthor.Username},
			CommentsCount: post.CommentsCount,
			LikesCount:    post.LikesCount,
			DislikesCount: post.DislikesCount,
			Categories:    post.Categories,
		})
	}

	server.SendObject(w, response)
}
