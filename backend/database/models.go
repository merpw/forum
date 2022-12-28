package database

import "time"

type Post struct {
	Id            int
	Title         string
	Content       string
	AuthorId      int
	Author        *User
	Date          string
	LikesCount    int
	DislikesCount int
	CommentsCount int
}

type User struct {
	Id       int
	Name     string
	Email    string
	Password string
}

type Comment struct {
	Id       int
	PostId   int
	Author   int
	Text     string
	Date     string
	Likes    int
	Dislikes int
}

type Session struct {
	Id        int64
	Uuid      string
	ExpiredAt time.Time
	UserId    int64
}
