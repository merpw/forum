package handlers

import (
	"backend/common/server"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// postsCreate creates a new post in the database
func (handlers *Handlers) postsCreate(w http.ResponseWriter, r *http.Request) {
	userId := handlers.getUserId(w, r)
	if userId == -1 {
		server.ErrorResponse(w, http.StatusUnauthorized)
		return
	}

	requestBody := struct {
		Title       string   `json:"title"`
		Content     string   `json:"content"`
		Description string   `json:"description"`
		Categories  []string `json:"categories"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Body is not valid", http.StatusBadRequest)
		return
	}

	requestBody.Title = strings.TrimSpace(requestBody.Title)
	requestBody.Content = strings.TrimSpace(requestBody.Content)
	requestBody.Description = strings.TrimSpace(requestBody.Description)

	if requestBody.Title == "" {
		http.Error(w, "Title is too short", http.StatusBadRequest)
		return
	}
	if len(requestBody.Title) > 25 {
		http.Error(w, "Title is too long, maximum length is 25", http.StatusBadRequest)
		return
	}

	if requestBody.Content == "" {
		http.Error(w, "Content is too short", http.StatusBadRequest)
		return
	}
	if len(requestBody.Content) > 10000 {
		http.Error(w, "Content is too long, maximum length is 10000", http.StatusBadRequest)
		return
	}

	if requestBody.Description == "" {
		requestBody.Description = shortenContent(requestBody.Content)
	}
	if len(requestBody.Description) > 200 {
		requestBody.Description = shortenContent(requestBody.Description)
	}

	for i, cat := range requestBody.Categories {
		cat = strings.TrimSpace(cat)
		cat = strings.ToLower(cat)
		requestBody.Categories[i] = cat
	}

	isValid := true
	for _, cat := range requestBody.Categories {
		if !isPresent(categories, cat) {
			isValid = false
			break
		}
	}
	if len(requestBody.Categories) == 0 {
		isValid = false
	}

	if !isValid {
		http.Error(w, "Categories are not valid", http.StatusBadRequest)
		return
	}

	id := handlers.DB.AddPost(requestBody.Title, requestBody.Content, requestBody.Description,
		userId, strings.Join(requestBody.Categories, ","))
	server.SendObject(w, id)

	err = revalidateURL(fmt.Sprintf("/post/%v", id))
	if err != nil {
		log.Printf("Error while revalidating `/post/%v`: %v", id, err)
	}
	err = revalidateURL("/")
	if err != nil {
		log.Printf("Error while revalidating '/': %v", err)
	}
	for _, category := range categories {
		err = revalidateURL(fmt.Sprintf("/category/%v", category))
		if err != nil {
			log.Printf("Error while revalidating `/category/%v`: %v", category, err)
		}
	}
}
