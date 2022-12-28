package server

import (
	"encoding/json"
	"fmt"
	"forum/database"
	"net/http"
	"time"
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
		srv.postsIdDislikeHandler(w, r)
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

// postsCatergoriesHandler returns a json list of all categories from the database
func (srv *Server) postsCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	// todo database post fetching
	sendObject(w, "posts categories list")
}

// postsHandler returns a json list of all posts from the database
func (srv *Server) postsHandler(w http.ResponseWriter, r *http.Request) {
	// no need of user to be logged in to see all posts
	// from the database, fetch the list of all the posts from the posts table

	// return data in json format, like below:
	// 	[
	//   {
	//     "id": 1,
	//     "title": "Post 1",
	//     "content": "Content One",
	//     "author": { "id": 1, "name": "Max" },
	//     "date": "2022-12-22T19:36:18.166Z",
	//     "likes": 1,
	//     "dislikes": 0,
	// 		"user_reaction": 1, // add only if user is logged in, -1=user disliked, 1=user liked, 0=nothing
	//     "comments_count": 2
	//   },
	//   {
	//     "id": 2,
	//     "title": "Post 2",
	//     "content": "Content Two.",
	//     "author": { "id": 2, "name": "Cat" },
	//     "date": "2022-12-22T19:36:18.166Z",
	//     "likes": 0,
	//     "dislikes": 0,
	// 		"user_reaction": 0, // add only if user is logged in, -1=user disliked, 1=user liked, 0=nothing
	//     "comments_count": 2
	//   }
	// ]

	fmt.Println("I reached in postsHandler") //-debug line

	db := database.OpenDB()

	// Execute the query to retrieve the posts
	rows, err := db.Query("SELECT * FROM posts")
	if err != nil {
		sendObject(w, err)
	}
	defer rows.Close()

	// Create a slice to store the posts
	var posts []database.Post

	// Iterate over the rows and append each post to the slice
	for rows.Next() {
		var post database.Post
		err := rows.Scan(&post.Id, &post.Title, &post.Content, &post.Author, &post.Date, &post.Likes, &post.Dislikes, &post.CommentsCount)
		if err != nil {
			sendObject(w, err)
		}
		posts = append(posts, post)
	}

	// Return the list of posts as a response to the client
	sendObject(w, posts)
	// sendObject(w, srv.posts)
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

// postsPostsIdHandler returns a single post from the database that matches the incoming id of the post in the url
func (srv *Server) postsPostsIdHandler(w http.ResponseWriter, r *http.Request) {
	// todo database managing etc
	sendObject(w, "postsPostsIdHandler")
}

// postsCreateHandler creates a new post in the database
func (srv *Server) postsCreateHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get the post data
	var newPost database.Post
	err := json.NewDecoder(r.Body).Decode(&newPost)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	// Validate the post data
	if newPost.Title == "" || newPost.Content == "" || newPost.Author.Id == 0 {
		http.Error(w, "Invalid post data", http.StatusBadRequest)
		return
	}

	// Connect to the database
	db := database.OpenDB()

	// Insert the new post into the database
	_, err = db.Exec(`INSERT INTO posts (title, content, author, date, likes, dislikes, user_reactions, comments_count) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, newPost.Title, newPost.Content, newPost.Author.Id, time.Now().Format(time.RFC3339), 0, 0, "", 0)
	if err != nil {
		http.Error(w, "Error inserting post into the database", http.StatusInternalServerError)
		return
	}

	// Return success to the client
	// w.WriteHeader(http.StatusOK)
	sendObject(w, http.StatusOK) // compromise for now, negotiate with the client to send the new post id back to the client if possible
	// sendObject(w, "postsCreateHandler")
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
