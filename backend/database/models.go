package database

import "time"

type Session struct {
	Id        int64     `json:"session_id"`
	Uuid      string    `json:"uuid"`
	ExpiredAt time.Time `json:"session_expired_at"`
	UserId    int64     `json:"session_user_id"`
}

type Post struct {
	Id             int      `json:"id"`
	Title          string   `json:"title"`
	Content        string   `json:"content"`
	Author         int      `json:"author_id"` // this is a user id
	Date           string   `json:"date"`
	Likes          int      `json:"likes"`
	Dislikes       int      `json:"dislikes"`
	UsersReactions string   `json:"users_reactions"`
	CommentsCount  int      `json:"comments_count"`
	Comments       []int    `json:"comments_ids"`
	Categories     []string `json:"categories"`
}

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Comment struct {
	Id            int         `json:"id"`             // comment id
	PostId        int         `json:"post_id"`        // post id
	Author        int         `json:"author_id"`      // author id
	Text          string      `json:"text"`           // comment text
	Date          string      `json:"date"`           // comment date
	Likes         int         `json:"likes"`          // comment likes
	Dislikes      int         `json:"dislikes"`       // comment dislikes
	UserReactions map[int]int `json:"user_reactions"` // map of user id to reaction value of -1 or 1
}
