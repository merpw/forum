package server

import (
	"net/http"
	"strings"
)

// postsCategoriesHandler returns a json list of all categories from the database
//
// Example: /api/posts/categories
func (srv *Server) postsCategoriesHandler(w http.ResponseWriter, _ *http.Request) {
	sendObject(w, categories)
}

// postsCategoriesNameHandler returns a json list of all posts belonging to the specified category
//
// Example: /api/posts/categories/name
func (srv *Server) postsCategoriesNameHandler(w http.ResponseWriter, r *http.Request) {
	categoryName := strings.TrimPrefix(r.URL.Path, "/api/posts/categories/")

	categoryName = strings.ToLower(categoryName)

	if !isPresent(categories, categoryName) {
		errorResponse(w, http.StatusNotFound)
		return
	}

	posts := srv.DB.GetCategoryPosts(categoryName)

	response := make([]SafePost, 0)
	for _, post := range posts {
		postAuthor := srv.DB.GetUserById(post.AuthorId)
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
