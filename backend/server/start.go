package server

import (
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
	router.HandleFunc("/api/posts/", server.postsHandler)

	// router.HandleFunc("/api/create", server.createHandler)

	router.HandleFunc("/api/auth/login", server.loginHandler)
	router.HandleFunc("/api/auth/signup", server.signUpHandler)
	router.HandleFunc("/api/auth/logout", server.logoutHandler)

	return router
}
