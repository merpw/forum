package handlers

import (
	"backend/common/server"
	"backend/forum/database"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// usersAll returns all userIds in alphabetical order.
//
// GET /api/users
func (h *Handlers) usersAll(w http.ResponseWriter, r *http.Request) {
	userIds := h.DB.GetAllUserIds()

	server.SendObject(w, userIds)
}

// usersId returns the info of the user with the given id
// If profile is Private, send only SafeUser, else, send entire user.
//
//	GET /api/users/:id
func (h *Handlers) usersId(w http.ResponseWriter, r *http.Request) {
	userIdStr := strings.TrimPrefix(r.URL.Path, "/api/users/")
	// /api/users/1 -> 1

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	user := h.DB.GetUserById(userId)
	if user == nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	response := struct {
		SafeUser
		Email     string `json:"email"`
		FirstName string `json:"first_name,omitempty"`
		LastName  string `json:"last_name,omitempty"`
		DoB       string `json:"dob,omitempty"`
		Gender    string `json:"gender,omitempty"`
		Avatar    string `json:"avatar,omitempty"`
		Bio       string `json:"bio,omitempty"`
	}{
		SafeUser: SafeUser{
			Id:           user.Id,
			Username:     user.Username,
			Avatar:       user.Avatar.String,
			Bio:          user.Bio.String,
			FollowStatus: nil,
		},
		Email:     user.Email,
		FirstName: user.FirstName.String,
		LastName:  user.LastName.String,
		DoB:       user.DoB.String,
		Gender:    user.Gender.String,
		Avatar:    user.Avatar.String,
		Bio:       user.Bio.String,
	}

	meId := h.getUserId(w, r)
	if meId != -1 {
		response.SafeUser.FollowStatus = h.DB.GetFollowStatus(meId, userId)
	} else {
		response.SafeUser.FollowStatus = nil
	}

	if response.SafeUser.FollowStatus == nil && user.Privacy == database.Public {
		server.SendObject(w, response)
		return
	}

	if response.SafeUser.FollowStatus == nil && user.Privacy == database.Private {
		server.SendObject(w, response.SafeUser)
		return
	}

	if user.Privacy == database.Public || *response.SafeUser.FollowStatus == database.Following {
		server.SendObject(w, response)
	} else {
		server.SendObject(w, response.SafeUser)
	}

}

func (h *Handlers) usersIdFollow(w http.ResponseWriter, r *http.Request) {
	userIdStr := strings.TrimPrefix(r.URL.Path, "/api/users/")
	userIdStr = strings.TrimSuffix(userIdStr, "/follow")
	// /api/users/1/ -> 1

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	user := h.DB.GetUserById(userId)

	if user == nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	meId := h.getUserId(w, r)
	followStatus := *h.DB.GetFollowStatus(meId, userId)
	fmt.Println("followStatus", followStatus)
	switch followStatus {
	case database.NotFollowing:
		if user.Privacy == database.Private {
			server.SendObject(w, h.DB.RequestToFollow(meId, userId))
			return
		}

		server.SendObject(w, h.DB.Follow(meId, userId))
		return

	case database.Following:
		server.SendObject(w, h.DB.Unfollow(meId, userId))
		return

	case database.RequestToFollow:
		server.SendObject(w, h.DB.RevokeInvitation(meId, userId))
		return
	}
}

// usersIdPosts Returns the posts of the user with the given id.
//
//	GET /api/users/:id/posts
func (h *Handlers) usersIdPosts(w http.ResponseWriter, r *http.Request) {
	userIdStr := strings.TrimPrefix(r.URL.Path, "/api/users/")
	userIdStr = strings.TrimSuffix(userIdStr, "/posts")
	// /api/users/1/posts -> 1

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	user := h.DB.GetUserById(userId)
	if user == nil {
		server.ErrorResponse(w, http.StatusNotFound)
		return
	}

	posts := h.DB.GetUserPosts(userId)

	response := make([]SafePost, 0)
	for _, post := range posts {
		response = append(response, SafePost{
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
		})
	}
	server.SendObject(w, response)
}
