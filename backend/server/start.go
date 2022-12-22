package server

import (
	"net/http"
)

type Server struct {
	posts string
}

func Start() http.Handler {
	server := Server{posts: "POSTS"}

	router := http.NewServeMux()

	router.HandleFunc("/api/posts", server.postsHandler)
	router.HandleFunc("/api/posts/", server.postHandler)

	router.HandleFunc("/api/create", server.createHandler)

	router.HandleFunc("/api/auth/login", server.postHandler)
	router.HandleFunc("/api/auth/signup", server.loginHandler)
	router.HandleFunc("/api/auth/logout", server.logoutHandler)

	return router
}
