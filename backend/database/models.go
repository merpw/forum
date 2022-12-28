package database

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

// type ApiUserIdPosts []Post

type Post struct {
	Id            int            `json:"id"`
	Title         string         `json:"title"`
	Content       string         `json:"content"`
	Author        Author         `json:"author"`
	Date          string         `json:"date"`
	Likes         int            `json:"likes"`
	Dislikes      int            `json:"dislikes"`
	UserReactions []UserReaction `json:"user_reactions"`
	CommentsCount int            `json:"comments_count"`
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
	Author       Author    `json:"author"`
	Date         string    `json:"date"`
	Likes        int       `json:"likes"`
	Dislikes     int       `json:"dislikes"`
	UserReaction int       `json:"user_reaction"`
	Comments     []Comment `json:"comments"`
}
type Comment struct {
	Author Author `json:"author"`
	PostId int    `json:"post_id"`
	Text   string `json:"text"`
	Date   string `json:"date"`
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
