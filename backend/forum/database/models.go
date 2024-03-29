package database

import (
	"database/sql"
)

type InviteStatus uint8

const (
	InviteStatusUnset InviteStatus = iota
	InviteStatusAccepted
	InviteStatusPending
)

type InviteType uint8

const (
	InviteTypeFollowUser InviteType = iota
	InviteTypeGroupInvite
	InviteTypeGroupJoin
	InviteTypeEvent
)

type Privacy uint8

const (
	Public Privacy = iota
	Private
	SuperPrivate
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
	Privacy       Privacy
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
	CreatorId   int
}

type Event struct {
	Id          int
	Title       string
	Description string
	TimeAndDate string
	Timestamp   string
	CreatedBy   int
	GroupId     int
}
