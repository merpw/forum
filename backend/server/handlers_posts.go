package server

import (
	"net/http"
)

func (srv *Server) apiPostsHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case reApiPosts.MatchString(r.URL.Path):
		srv.postsHandler(w, r)
	case reApiPostsCategories.MatchString(r.URL.Path):
		srv.postsCategoriesHandler(w, r)
	case reApiPostsCategoriesFacts.MatchString(r.URL.Path):
		srv.postsCategoriesFactsHandler(w, r)
	case reApiPostsCategoriesRumors.MatchString(r.URL.Path):
		srv.postsCategoriesRumorsHandler(w, r)
	case reApiPostsId.MatchString(r.URL.Path):
		srv.postsPostsIdHandler(w, r)
	case reApiPostsCreate.MatchString(r.URL.Path):
		srv.postsCreateHandler(w, r)
	case reApiPostsIdLike.MatchString(r.URL.Path):
		srv.postsIdLikeHandler(w, r)
	case reApiPostsIdDislike.MatchString(r.URL.Path):
		srv.postsPostsIdDislikeHandler(w, r)
	case reApiPostsIdComment.MatchString(r.URL.Path):
		srv.postsIdCommentHandler(w, r)
	case reApiPostsIdCommentIdLike.MatchString(r.URL.Path):
		srv.postsIdCommentIdLikeHandler(w, r)
	case reApiPostsIdCommentIdDislike.MatchString(r.URL.Path):
		srv.postsIdCommentIdDislikeHandler(w, r)
		// case reApiSignup.MatchString(r.URL.Path):
		// 	srv.signupHandler(w, r)
		// case reApiLogin.MatchString(r.URL.Path):
		// 	srv.loginHandler(w, r)
		// case reApiLogout.MatchString(r.URL.Path):
		// 	srv.logoutHandler(w, r)
	}
}

func (srv *Server) postsCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	// todo database post fetching
	sendObject(w, "posts categories list")
}
func (srv *Server) postsHandler(w http.ResponseWriter, r *http.Request) {
	// todo database post fetching
	sendObject(w, srv.posts)
}
func (srv *Server) postsCategoriesFactsHandler(w http.ResponseWriter, r *http.Request) {
	// todo database managing etc
	sendObject(w, "postsCategoriesFactsHandler")
}
func (srv *Server) postsCategoriesRumorsHandler(w http.ResponseWriter, r *http.Request) {
	// todo database managing etc
	sendObject(w, "postsCategoriesRumorsHandler")
}
func (srv *Server) postsPostsIdHandler(w http.ResponseWriter, r *http.Request) {
	// todo database managing etc
	sendObject(w, "postsPostsIdHandler")
}
func (srv *Server) postsCreateHandler(w http.ResponseWriter, r *http.Request) {
	// todo database managing etc
	sendObject(w, "postsCreateHandler")
}
func (srv *Server) postsIdLikeHandler(w http.ResponseWriter, r *http.Request) {
	// todo database managing etc
	sendObject(w, "postsIdLikeHandler")
}
func (srv *Server) postsPostsIdDislikeHandler(w http.ResponseWriter, r *http.Request) {
	// todo database managing etc
	sendObject(w, "postsPostsIdDislikeHandler")
}
func (srv *Server) postsIdCommentHandler(w http.ResponseWriter, r *http.Request) {
	// todo database managing etc
	sendObject(w, "postsIdCommentHandler")
}
func (srv *Server) postsIdCommentIdLikeHandler(w http.ResponseWriter, r *http.Request) {
	// todo database managing etc
	sendObject(w, "postsPostsIdCommentIdLikeHandler")
}
func (srv *Server) postsIdCommentIdDislikeHandler(w http.ResponseWriter, r *http.Request) {
	// todo database managing etc
	sendObject(w, "postsIdCommentIdDislikeHandler")
}
