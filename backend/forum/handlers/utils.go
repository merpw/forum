package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

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

// shortenContent shortens content to 200 characters, adds "..." at the end
func shortenContent(content string) string {
	if len(content) > 200 {
		return content[:200] + "..."
	}
	return content
}

func isPresent(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// revalidateURL creates POST request to frontend to revalidate url
//
// Uses environment variables FRONTEND_REVALIDATE_URL and optional FRONTEND_REVALIDATE_TOKEN
//
// Does nothing if FRONTEND_REVALIDATE_URL is not set
func revalidateURL(url string) error {
	apiURL := os.Getenv("FRONTEND_REVALIDATE_URL")
	if apiURL == "" {
		return nil
	}
	req, err := http.NewRequest(http.MethodPost, apiURL, nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("url", url)
	q.Add("token", os.Getenv("FRONTEND_REVALIDATE_TOKEN"))
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("revalidation failed: %s, %s", res.Status, bodyBytes)
	}

	return nil
}