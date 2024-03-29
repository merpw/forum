package handlers

import (
	. "backend/forum/database"
)

type SafeUser struct {
	Id           int           `json:"id"`
	Username     string        `json:"username"`
	Avatar       string        `json:"avatar,omitempty"`
	Bio          string        `json:"bio,omitempty"`
	FollowStatus *InviteStatus `json:"follow_status,omitempty"`
	Followers    int           `json:"followers_count"`
	Following    int           `json:"following_count"`
	Privacy      bool          `json:"privacy"`
}

type SafePost struct {
	Id            int      `json:"id"`
	Title         string   `json:"title"`
	Content       string   `json:"content"`
	Description   string   `json:"description"`
	Author        SafeUser `json:"author"`
	Date          string   `json:"date"`
	CommentsCount int      `json:"comments_count"`
	LikesCount    int      `json:"likes_count"`
	DislikesCount int      `json:"dislikes_count"`
	Categories    string   `json:"categories"`
	GroupId       *int     `json:"group_id,omitempty"`
}

type SafeComment struct {
	Id            int      `json:"id"`
	Content       string   `json:"content"`
	Author        SafeUser `json:"author"`
	Date          string   `json:"date"`
	LikesCount    int      `json:"likes_count"`
	DislikesCount int      `json:"dislikes_count"`
}

type SafeReaction struct {
	Reaction      int `json:"reaction"`
	LikesCount    int `json:"likes_count"`
	DislikesCount int `json:"dislikes_count"`
}

type SafeEvent struct {
	Id          int          `json:"id"`
	GroupId     int          `json:"group_id"`
	CreatedBy   int          `json:"created_by"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	TimeAndDate string       `json:"time_and_date"`
	Timestamp   string       `json:"timestamp"`
	Status      InviteStatus `json:"status"`
}

func isPresent(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
