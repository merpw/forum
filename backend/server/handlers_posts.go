package server

import (
	"net/http"
)

var categories = []string{"facts", "rumors", "other"}

func (srv *Server) apiPostsMasterHandler(w http.ResponseWriter, r *http.Request) {
	switch {

	case reApiPosts.MatchString(r.URL.Path):
		srv.postsHandler(w, r) // /api/posts

	case reApiPostsCategories.MatchString(r.URL.Path):
		srv.postsCategoriesHandler(w, r) // /api/posts/categories - list of all categories in forum

	case reApiPostsCategoriesName.MatchString(r.URL.Path):
		srv.postsCategoriesNameHandler(w, r) // /api/posts/categories/{name} - list of all posts in category

	case reApiPostsId.MatchString(r.URL.Path):
		srv.postsIdHandler(w, r) // /api/posts/{id} - get post by id

	case reApiPostsCreate.MatchString(r.URL.Path):
		srv.postsCreateHandler(w, r) // /api/posts/create - create new post

	case reApiPostsIdLike.MatchString(r.URL.Path):
		srv.postsIdLikeHandler(w, r) // /api/posts/{id}/like - like post

	case reApiPostsIdDislike.MatchString(r.URL.Path):
		srv.postsIdDislikeHandler(w, r) // /api/posts/{id}/dislike - dislike post

	case reApiPostsIdReaction.MatchString(r.URL.Path):
		srv.postsIdReactionHandler(w, r) // /api/posts/{id}/reaction - get post reaction of current user

	// /api/posts/{id}/comment/{id}/reaction - get comment reaction of current user
	case reApiPostsIdCommentIdReaction.MatchString(r.URL.Path):
		srv.postsIdCommentIdReactionHandler(w, r)

	case reApiPostsIdComment.MatchString(r.URL.Path):
		srv.postsIdCommentCreateHandler(w, r) // /api/posts/{id}/comment - create new comment

	case reApiPostsIdComments.MatchString(r.URL.Path):
		srv.postsIdCommentsHandler(w, r) // /api/posts/{id}/comments - get all comments of post

	case reApiPostsIdCommentIdLike.MatchString(r.URL.Path):
		srv.postsIdCommentIdLikeHandler(w, r) // /api/posts/{id}/comment/{id}/like - like comment

	case reApiPostsIdCommentIdDislike.MatchString(r.URL.Path):
		srv.postsIdCommentIdDislikeHandler(w, r) // /api/posts/{id}/comment/{id}/dislike - dislike comment
	default:
		http.NotFound(w, r)
	}
}

// postsHandler returns a json list of all posts in the forum.
//
// Example: /api/posts
func (srv *Server) postsHandler(w http.ResponseWriter, _ *http.Request) {
	posts := srv.DB.GetAllPosts()

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
