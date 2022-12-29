package database

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
	Id     int
	Token  string
	Expire int
	UserId int
}

type Like struct {
	Id       int
	PostLike int // 1 (like for post) 0 (like for the comment)
	ParentId int // id of the liked node (post or comment). Not a foreign key, because two types of id from post or comment
	Value    int // 1 (like) -1 (dislike)
}
