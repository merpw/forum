package database

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
}

type User struct {
	Id       int
	Name     string
	Email    string
	Password string
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

// TODO: remove this later

// chat section
type Chat struct {
	Id              int
	LastMessageDate string
}

type Membership struct {
	Id     int
	ChatId int
	UserId int
	Date   string
}

type Message struct {
	Id       int
	ChatId   int
	SenderId int
	Body     string
	Date     string
}
