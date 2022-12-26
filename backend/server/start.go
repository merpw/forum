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

	// /api/posts ... see the postHandler variable "commands" description
	router.HandleFunc("/api/posts/", server.postsHandler) //trailing slash removed

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

		router.ServeHTTP(w, r)
	})
}
