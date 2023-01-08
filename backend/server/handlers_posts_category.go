package server

import (
	"net/http"
	"strings"
)

// postsCategoriesHandler returns a json list of all categories from the database
func (srv *Server) postsCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	sendObject(w, categories)
}

func (srv *Server) postsCategoriesNameHandler(w http.ResponseWriter, r *http.Request) {
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

	posts := srv.DB.GetCategoryPosts(categoryName)

	response := make([]SafePost, 0)
	for _, post := range posts {
		postAuthor := srv.DB.GetUserById(post.AuthorId)
		cutPostContentForLists(&post)
		response = append(response, SafePost{
			Id:            post.Id,
			Title:         post.Title,
			Content:       post.Content,
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
