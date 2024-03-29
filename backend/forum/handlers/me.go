package handlers

import (
	"backend/common/server"
	. "backend/forum/database"
	"net/http"
)

// me returns the currently logged-in user's information.
func (h *Handlers) me(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userIdCtxKey).(int)

	user := h.DB.GetUserById(userId)
	response := struct {
		SafeUser
		Email     string `json:"email"`
		FirstName string `json:"first_name,omitempty"`
		LastName  string `json:"last_name,omitempty"`
		DoB       string `json:"dob,omitempty"`
		Gender    string `json:"gender,omitempty"`
	}{
		SafeUser: SafeUser{
			Id:        user.Id,
			Username:  user.Username,
			Avatar:    user.Avatar.String,
			Bio:       user.Bio.String,
			Followers: len(h.DB.GetUserFollowers(user.Id)),
			Following: len(h.DB.GetUsersFollowed(user.Id)),
			Privacy:   user.Privacy == Private,
		},
		Email:     user.Email,
		FirstName: user.FirstName.String,
		LastName:  user.LastName.String,
		DoB:       user.DoB.String,
		Gender:    user.Gender.String,
	}

	server.SendObject(w, response)
}

// mePrivacy toggles privacy and sends current privacy status
func (h *Handlers) mePrivacy(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userIdCtxKey).(int)
	user := h.DB.GetUserById(userId)

	if user.Privacy == Private {
		server.SendObject(w, h.DB.UpdateUserPrivacy(Public, userId))
	} else {
		server.SendObject(w, h.DB.UpdateUserPrivacy(Private, userId))
	}
}

// meFollowers sends all followers as IDs in an array
// /api/me/followers -> [followerIds]
func (h *Handlers) meFollowers(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userIdCtxKey).(int)

	server.SendObject(w, h.DB.GetUserFollowers(userId))
}

func (h *Handlers) meFollowing(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userIdCtxKey).(int)

	server.SendObject(w, h.DB.GetUsersFollowed(userId))
}

// mePosts returns the posts of the logged-in user.
func (h *Handlers) mePosts(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userIdCtxKey).(int)

	user := h.DB.GetUserById(userId)
	posts := h.DB.GetMePosts(userId)

	type Response struct {
		SafePost
	}

	response := make([]Response, 0)
	for _, post := range posts {
		response = append(response, Response{
			SafePost: SafePost{
				Id:          post.Id,
				Title:       post.Title,
				Description: post.Description,
				Author: SafeUser{
					Id:       user.Id,
					Username: user.Username,
					Avatar:   user.Avatar.String,
					Bio:      user.Bio.String,
				},
				Date:          post.Date,
				CommentsCount: post.CommentsCount,
				LikesCount:    post.LikesCount,
				DislikesCount: post.DislikesCount,
				Categories:    post.Categories,
			},
		})
	}

	server.SendObject(w, response)
}

// mePostsLiked returns the posts liked by the logged-in user.
func (h *Handlers) mePostsLiked(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userIdCtxKey).(int)

	posts := h.DB.GetUserPostsLiked(userId)

	type Response struct {
		SafePost
		DislikesCount int `json:"dislikesCount"`
	}

	response := make([]Response, 0)
	for _, post := range posts {
		author := h.DB.GetUserById(post.AuthorId)
		response = append(response, Response{
			SafePost: SafePost{
				Id:          post.Id,
				Title:       post.Title,
				Description: post.Description,
				Author: SafeUser{
					Id:       author.Id,
					Username: author.Username,
					Avatar:   author.Avatar.String,
					Bio:      author.Bio.String,
				},
				Date:          post.Date,
				CommentsCount: post.CommentsCount,
				LikesCount:    post.LikesCount,
				DislikesCount: post.DislikesCount,
				Categories:    post.Categories,
			},
		})
	}
	server.SendObject(w, response)
}
