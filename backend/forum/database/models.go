package database

import (
	"database/sql"
)

type InviteStatus uint8

const (
	Inactive InviteStatus = iota
	Accepted
	Pending
)

type InviteType uint8

const (
	FollowUser InviteType = iota
	GroupInvite
	GroupJoin
	Event
)

type Privacy uint8

const (
	Public Privacy = iota
	Private
)

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
	GroupId       *int
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
	Privacy   Privacy
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

type Invitation struct {
	Id           int
	Type         InviteType
	FromUserId   int
	ToUserId     int
	AssociatedId sql.NullInt64
	TimeStamp    string
}

type Group struct {
	Id          int
	Title       string
	Description string
	Members     int
}
