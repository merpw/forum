package server

import (
	"encoding/json"
	"forum/database"
	"net/http"
	"strconv"
	"strings"
)

func (srv *Server) apiPostsMasterHandler(w http.ResponseWriter, r *http.Request) {
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
		srv.postsIdHandler(w, r)

	case reApiPostsCreate.MatchString(r.URL.Path):
		srv.postsCreateHandler(w, r)

	case reApiPostsIdLike.MatchString(r.URL.Path):
		srv.postsIdLikeHandler(w, r)

	case reApiPostsIdDislike.MatchString(r.URL.Path):
		srv.postsIdDislikeHandler(w, r)

	case reApiPostsIdComment.MatchString(r.URL.Path):
		srv.postsIdCommentHandler(w, r)

	case reApiPostsIdCommentIdLike.MatchString(r.URL.Path):
		srv.postsIdCommentIdLikeHandler(w, r)

	case reApiPostsIdCommentIdDislike.MatchString(r.URL.Path):
		srv.postsIdCommentIdDislikeHandler(w, r)
	default:
		http.NotFound(w, r)
	}
}

// postsCategoriesHandler returns a json list of all categories from the database
func (srv *Server) postsCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	// todo database post fetching
	sendObject(w, "posts categories list")
}

// postsHandler returns a json list of all posts from the database
func (srv *Server) postsHandler(w http.ResponseWriter, r *http.Request) {
	sendObject(w, srv.DB.GetPosts())
}

// postsCategoriesFactsHandler returns a json list of all posts from the database that match the category "facts"
func (srv *Server) postsCategoriesFactsHandler(w http.ResponseWriter, r *http.Request) {
	// todo database managing etc
	sendObject(w, "postsCategoriesFactsHandler")
}

// postsCategoriesRumorsHandler returns a json list of all posts from the database that match the category "rumors"
func (srv *Server) postsCategoriesRumorsHandler(w http.ResponseWriter, r *http.Request) {
	// todo database managing etc
	sendObject(w, "postsCategoriesRumorsHandler")
}

// postsIdHandler returns a single post from the database that matches the incoming id of the post in the url
func (srv *Server) postsIdHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/posts/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse(w, http.StatusNotFound)
		return
	}

	// Get the post from the database
	post := srv.DB.GetPostById(id)

	if post == nil {
		errorResponse(w, http.StatusNotFound)
		return
	}

	// Send the post back to the client
	sendObject(w, post)
}

// postsCreateHandler creates a new post in the database
func (srv *Server) postsCreateHandler(w http.ResponseWriter, r *http.Request) {
	var post database.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "Body is not valid", http.StatusBadRequest)
		return
	}
	post.Author = 1 // TODO: remove when add auth

	if post.Title == "" || post.Content == "" || post.Author == 0 {
		http.Error(w, "Invalid post data", http.StatusBadRequest)
		return
	}

	id := srv.DB.AddPost(post)
	sendObject(w, id)
}

// postsIdLikeHandler likes a post in the database
func (srv *Server) postsIdLikeHandler(w http.ResponseWriter, r *http.Request) {
	// todo database managing etc
	sendObject(w, "postsIdLikeHandler")
}

// postsPostsIdDislikeHandler dislikes a post in the database
func (srv *Server) postsIdDislikeHandler(w http.ResponseWriter, r *http.Request) {
	// todo database managing etc
	sendObject(w, "postsIdDislikeHandler")
}

// postsIdCommentHandler comments on a post in the database
func (srv *Server) postsIdCommentHandler(w http.ResponseWriter, r *http.Request) {
	// todo database managing etc
	sendObject(w, "postsIdCommentHandler")
}

// postsIdCommentIdLikeHandler likes a comment on a post in the database
func (srv *Server) postsIdCommentIdLikeHandler(w http.ResponseWriter, r *http.Request) {
	// requirement: the user needs to be logged in to do this.
	// step1: check if the user is logged in
	// step2: check if the post exists
	// step3: check if the comment exists
	//
	// step4: check if the user has already his id in the list of user_reactions
	// if user_id not in user_reactions{
	// user_reaction[user_id]=1 // like
	// like++ // for the post like number
	// }else{ // user_id inside database
	// switch user_reaction[user_id] { // old reaction from database
	// case 1: //according to request, because it is "like" handler so like press happens here
	// like-- // decrease like
	// delete(user_reaction[user_id])
	// case -1: // was dislike
	// dislike--
	// like++
	// }
	// }
	// step 5: send the updated number of likes/dislikes on post to the frontend
	// reaction = user_reaction[user_id] // 1 like -1 dislike 0- if no key inside map, then send it to the frontend
	// todo database managing etc
	sendObject(w, "postsIdCommentIdLikeHandler")
}

// postsIdCommentIdDislikeHandler dislikes a comment on a post in the database
func (srv *Server) postsIdCommentIdDislikeHandler(w http.ResponseWriter, r *http.Request) {
	// the user needs to be logged in to do this.
	// todo step1: check if the user is logged in
	// todo step2: check if the post exists
	// check if the comment exists
	//
	// check if the user has already his id in the list of user_reactions
	// if not, add the user id to the list of user_reactions
	// if yes, check the reaction:value against the user_id from the list of user_reactions
	// if the reaction:value is the same as the user_id from the list of user_reactions, make it 0,
	// plus decrease the value of likes/dislikes on post depending upon whether user_reaction_value was 1 or -1, that we made into 0
	// and remove the user_id from the list of user_reactions,
	//
	// else if the reaction:value is not the same as the user_id from the list of user_reactions, make it according
	// to the reaction:value, plus increase the value of likes/dislikes on post depending upon whether user_reaction_value was 1 or -1
	// and add the user_id to the list of user_reactions,
	//
	// send the updated number of likes/dislikes on post to the frontend
	// send the updated user_reaction for the user_id to the frontend.
	// if user_if is now not in the list of user_reactions, send 0 to the frontend
	// todo
	sendObject(w, "postsIdCommentIdDislikeHandler")
}
