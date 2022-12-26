package server

import (
	"log"
	"net/http"
	"regexp"
)

type Server struct {
	posts string
}

// Start returns http.Handler with all routes
func Start() http.Handler {
	// server := Server{posts: "POSTS"}

	router := http.NewServeMux()

	// /api/posts ... see the postHandler variable "commands" description
	// router.HandleFunc("/api/posts/", server.postsHandler) //trailing slash removed

	// router.HandleFunc("/api/user/", server.usersHandler)

	// router.HandleFunc("/api/auth/login", server.loginHandler)
	// router.HandleFunc("/api/auth/signup", server.signUpHandler)
	// router.HandleFunc("/api/auth/logout", server.logoutHandler)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("ERROR %d. %v\n", http.StatusInternalServerError, err)
				errorResponse(w, http.StatusInternalServerError) // 500 ERROR
			}
		}()

		// todo add to regexp all variants for urls for GET method, to check incoming request urls fast
		reGet := regexp.MustCompile("")

		// todo add to regexp all variants for urls for POST method, to check incoming request urls fast
		rePost := regexp.MustCompile("")

		// todo natkim, ask maxim what wrong here.
		// i need to sleep now :x . And this code section i plan to use for checking
		// the wrong method used for incoming requests. Every regexp include all
		// allowed url schemes for each of two methods, to check once. Perhaps it is wrong
		// place for this block of code
		switch r.Method {
		case http.MethodGet:
			if !reGet.Match([]byte(r.Method)) {
				errorResponse(w, http.StatusMethodNotAllowed)
			}
		case http.MethodPost:
			if !rePost.Match([]byte(r.Method)) {
				errorResponse(w, http.StatusMethodNotAllowed)
			}
		default:
			log.Printf("ERROR: %d\n", http.StatusMethodNotAllowed)
			errorResponse(w, http.StatusMethodNotAllowed) // 405 ERROR
		}

		router.ServeHTTP(w, r)
	})
}
