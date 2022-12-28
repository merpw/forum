package database

import "time"

type ApiMe struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ApiMeLikedPosts struct {
}

type ApiUserId struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Session struct {
	Id        int64     `json:"session_id"`
	Uuid      string    `json:"uuid"`
	ExpiredAt time.Time `json:"session_expired_at"`
	UserId    int64     `json:"session_user_id"`
}

// type ApiUserIdPosts []Post

// This holds all the categories on the forum
var AllCategories = []string{"Technology", "Facts", "Rumors", "Science", "Politics", "Sports", "Entertainment", "Health", "Business", "Other"}

type Post struct {
	Id            int         `json:"id"`
	Title         string      `json:"title"`
	Content       string      `json:"content"`
	Author        int         `json:"author_id"` // this is a user id
	Date          string      `json:"date"`
	Likes         int         `json:"likes"`
	Dislikes      int         `json:"dislikes"`
	UserReactions map[int]int `json:"user_reactions"` // map of user id to reaction value of -1 or 1
	CommentsCount int         `json:"comments_count"`
	Comments      []int       `json:"comments_ids"`
	Categories    []string    `json:"categories"`
}

type UserReaction struct {
	UserId         int `json:"user_id"`
	Reaction_Value int `json:"reaction_value"` // 1 for like, -1 for dislike, 0 for no reaction
}
type Author struct {
	Id   int    `json:"user_id"`
	Name string `json:"user_name"`
}

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Posts    []int  `json:"posts_id"`
	Comments []int  `json:"comments_id"`
}

type ApiPosts []Post

type ApiPostsCategories struct {
}

type ApiPostsCategoryFacts struct {
}

type ApiPostsCategoryRumors struct {
}

type ApiPostsId struct {
	Id           int       `json:"id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	Author       int       `json:"author_id"`
	Date         string    `json:"date"`
	Likes        int       `json:"likes"`
	Dislikes     int       `json:"dislikes"`
	UserReaction int       `json:"user_reaction"`
	Comments     []Comment `json:"comments"`
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

type ApiPostsCreate struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ApiSignup struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ApiLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type ApiLogout struct {
	// according to notion, looks like empty, or no needed
}

type ApiPostsIdLike struct {
}

type ApiPostsIdDislike struct {
}

type ApiPostsIdComment struct {
}

type ApiPostsIdCommentIdLike struct {
}

type ApiPostsIdCommentIdDislike struct {
}
