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
	Id            int    `json:"id"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	Author        Author `json:"author"`
	Date          string `json:"date"`
	Likes         int    `json:"likes"`
	Dislikes      int    `json:"dislikes"`
	UserReaction  int    `json:"user_reaction"`
	CommentsCount int    `json:"comments_count"`
}
type Author struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
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
