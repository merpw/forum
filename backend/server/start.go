package server

import (
	"log"
	"net/http"
)

type Server struct {
	posts string
}

// Start returns http.Handler with all routes
func Start() http.Handler {
	server := Server{posts: "POSTS"}

	router := http.NewServeMux()

	// Master-handler for:
	// /api/posts, /api/posts/{id}, /api/posts/{id}/like, /api/posts/{id}/dislike
	// /api/posts/{id}/comment, /api/posts/{id}/comment/{id}/like, /api/posts/{id}/comment/{id}/dislike
	router.HandleFunc("/api/posts/", server.postsHandler)

	// Master-handler for:
	// /api/user, /api/user/{id}, /api/user/{id}/posts, /api/user/{id}
	router.HandleFunc("/api/user/", server.usersHandler)

	router.HandleFunc("/api/auth/login", server.loginHandler)
	router.HandleFunc("/api/auth/signup", server.signUpHandler)
	router.HandleFunc("/api/auth/logout", server.logoutHandler)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("ERROR %d. %v\n", http.StatusInternalServerError, err)
				errorResponse(w, http.StatusInternalServerError) // 500 ERROR
			}
		}()

		switch {
		case GetRegexp.MatchString(r.URL.Path):
			if r.Method != http.MethodGet {
				errorResponse(w, http.StatusMethodNotAllowed)
				return
			}
		case PostRegexp.MatchString(r.URL.Path):
			if r.Method != http.MethodPost {
				errorResponse(w, http.StatusMethodNotAllowed)
				return
			}
		default:
			errorResponse(w, http.StatusNotFound)
			return
		}
		router.ServeHTTP(w, r)
	})
}
