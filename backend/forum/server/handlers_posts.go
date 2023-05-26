package server

import (
	"net/http"
)

var categories = []string{"facts", "rumors", "other"}

func (srv *Server) apiPostsMasterHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case reApiPosts.MatchString(r.URL.Path):
		srv.postsHandler(w, r)

	case reApiPostsCategories.MatchString(r.URL.Path):
		srv.postsCategoriesHandler(w, r)

	case reApiPostsCategoriesName.MatchString(r.URL.Path):
		srv.postsCategoriesNameHandler(w, r)

	case reApiPostsId.MatchString(r.URL.Path):
		srv.postsIdHandler(w, r)

	case reApiPostsCreate.MatchString(r.URL.Path):
		srv.postsCreateHandler(w, r)

	case reApiPostsIdLike.MatchString(r.URL.Path):
		srv.postsIdLikeHandler(w, r)

	case reApiPostsIdDislike.MatchString(r.URL.Path):
		srv.postsIdDislikeHandler(w, r)

	case reApiPostsIdReaction.MatchString(r.URL.Path):
		srv.postsIdReactionHandler(w, r)

	case reApiPostsIdCommentIdReaction.MatchString(r.URL.Path):
		srv.postsIdCommentIdReactionHandler(w, r)

	case reApiPostsIdComment.MatchString(r.URL.Path):
		srv.postsIdCommentCreateHandler(w, r)

	case reApiPostsIdComments.MatchString(r.URL.Path):
		srv.postsIdCommentsHandler(w, r)

	case reApiPostsIdCommentIdLike.MatchString(r.URL.Path):
		srv.postsIdCommentIdLikeHandler(w, r)

	case reApiPostsIdCommentIdDislike.MatchString(r.URL.Path):
		srv.postsIdCommentIdDislikeHandler(w, r)
	default:
		http.NotFound(w, r)
	}
}

// postsHandler returns a json list of all posts from the database
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
