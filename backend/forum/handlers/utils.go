package handlers

type SafeUser struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type SafePost struct {
	Id            int      `json:"id"`
	Title         string   `json:"title"`
	Content       string   `json:"content"`
	Description   string   `json:"description"`
	Author        SafeUser `json:"author"`
	Date          string   `json:"date"`
	CommentsCount int      `json:"comments_count"`
	LikesCount    int      `json:"likes_count"`
	DislikesCount int      `json:"dislikes_count"`
	Categories    string   `json:"categories"`
}

type SafeComment struct {
	Id            int      `json:"id"`
	Content       string   `json:"content"`
	Author        SafeUser `json:"author"`
	Date          string   `json:"date"`
	LikesCount    int      `json:"likes_count"`
	DislikesCount int      `json:"dislikes_count"`
}

type SafeReaction struct {
	Reaction      int `json:"reaction"`
	LikesCount    int `json:"likes_count"`
	DislikesCount int `json:"dislikes_count"`
}

func isPresent(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
