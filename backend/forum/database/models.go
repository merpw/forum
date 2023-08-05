package database

import "database/sql"

type Post struct {
	Id            int
	Title         string
	Content       string
	AuthorId      int
	Date          string
	LikesCount    int
	DislikesCount int
	CommentsCount int
	Categories    string
	Description   string
}

type User struct {
	Id        int
	Username  string
	Email     string
	Password  string
	FirstName sql.NullString
	LastName  sql.NullString
	DoB       sql.NullString
	Gender    sql.NullString
	Avatar    sql.NullString
	Bio       sql.NullString
}

type Comment struct {
	Id            int
	PostId        int
	AuthorId      int
	Content       string
	Date          string
	LikesCount    int
	DislikesCount int
}

type Session struct {
	Id     int
	Token  string
	Expire int
	UserId int
}
