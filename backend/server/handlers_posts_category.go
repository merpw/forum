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
	type ResponsePost struct {
		Id            int      `json:"id"`
		Title         string   `json:"title"`
		Content       string   `json:"content"`
		Author        SafeUser `json:"author"`
		Date          string   `json:"date"`
		CommentsCount int      `json:"comments_count"`
		LikesCount    int      `json:"likes_count"`
		Category      string   `json:"category"`
	}

	response := make([]ResponsePost, 0)
	for _, post := range posts {
		postAuthor := srv.DB.GetUserById(post.AuthorId)
		response = append(response, ResponsePost{
			Id:            post.Id,
			Title:         post.Title,
			Content:       post.Content,
			Date:          post.Date,
			Author:        SafeUser{Id: postAuthor.Id, Name: postAuthor.Name},
			CommentsCount: post.CommentsCount,
			LikesCount:    post.LikesCount,
			Category:      post.Category,
		})
	}

	sendObject(w, response)
}
