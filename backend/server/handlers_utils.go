package server

import (
	"encoding/json"
	"fmt"
	"forum/database"
	"log"
	"net/http"
)

type SafeUser struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type SafePost struct {
	Id            int      `json:"id"`
	Title         string   `json:"title"`
	Content       string   `json:"content"`
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

// errorResponse responses with specified error code in format "404 Not Found"
func errorResponse(w http.ResponseWriter, code int) {
	http.Error(w, fmt.Sprintf("%v %v", code, http.StatusText(code)), code)
}

// SendObject sends object to http.ResponseWriter
//
// calls errorResponse(500) if error happened
func SendObject(w http.ResponseWriter, object any) {
	w.Header().Set("Content-Type", "application/json")
	objJson, err := json.Marshal(object)
	if err != nil {
		log.Println(err)
		errorResponse(w, 500)
		return
	}
	_, err = w.Write(objJson)
	if err != nil {
		log.Println(err)
		errorResponse(w, 500)
		return
	}
}

func cutPostContentForLists(post *database.Post) {
	if len(post.Content) > 200 {
		post.Content = post.Content[:200] + "..."
	}
}
