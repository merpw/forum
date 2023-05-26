package handlers

import (
	"net/http"
	"strings"
)

// postsCategories returns a json list of all categories from the database
func (handlers *Handlers) postsCategories(w http.ResponseWriter, _ *http.Request) {
	sendObject(w, categories)
}

func (handlers *Handlers) postsCategoriesName(w http.ResponseWriter, r *http.Request) {
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
		errorResponse(w, http.StatusNotFound)
		return
	}

	posts := handlers.DB.GetCategoryPosts(categoryName)

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
