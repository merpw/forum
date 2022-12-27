package server

import (
	"fmt"
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
		// reGet := regexp.MustCompile(`\/api\/(?:user\/[[:digit:]]+(?:\/posts(?:\/liked)?)?|posts(?:\/(?:rumor|fact)s)?|me(?:\/liked)?)\/?$`)
		reGet := regexp.MustCompile(`^\/api\/(?:post(?:s\/category\/[a-zA-Z0-9_.-]+|\/[[:digit:]]+)\/?|posts\/categories\/?|me\/liked\/posts\/?|user\/[[:digit:]]+(?:\/(?:posts\/?)?)?|posts\/?|me\/?)$`)

		// todo add to regexp all variants for urls for POST method, to check incoming request urls fast
		// rePost := regexp.MustCompile(`\/api\/post(?:\/[[:digit:]]+/(?:(?:comment\/[[:digit:]]+\/)?dislike|(?:comment\/[[:digit:]]\/)?like|comment))?\/?$`) // perhaps this regexp, made by natkim does the same filtering
		rePost := regexp.MustCompile(`^\/api\/post(?:\/(?:[[:digit:]]+\/(?:(?:comment\/[[:digit:]]+\/)?dislike\/?|(?:comment\/[[:digit:]]+\/)?like\/?|comment\/?))?)?$`)

		// todo cococore, ask maxim what wrong here.
		// i need to sleep now :x . And this code section i plan to use for checking
		// the wrong method used for incoming requests. Every regexp include all
		// allowed url schemes for each of two methods, to check once. Perhaps it is wrong
		// place for this block of code

		//not completed, wrong errors. need to be fixed
		fmt.Println("r.URL.Path = ", r.URL.Path, "len", len(r.URL.Path))
		if reGet.Match([]byte(r.URL.Path)) {
			if r.Method != http.MethodGet {
				fmt.Println("===TEST inside GET regex fired\n") // todo remove later
				errorResponse(w, http.StatusMethodNotAllowed)
			}
		} else if rePost.Match([]byte(r.URL.Path)) {
			if r.Method != http.MethodPost {
				fmt.Println("===TEST inside POST regex fired\n") // todo remove later
				errorResponse(w, http.StatusMethodNotAllowed)
			}
		} else {
			fmt.Println("===TEST inside no get no post regex fired\n") // todo remove later
			errorResponse(w, http.StatusNotFound)
		}

		/*
			switch r.Method {
			case http.MethodGet:
				if !reGet.Match([]byte(r.Method)) {
					errorResponse(w, http.StatusNotFound)
				}
			case http.MethodPost:
				if !rePost.Match([]byte(r.Method)) {
					errorResponse(w, http.StatusMethodNotAllowed)
				}
			default:
				log.Printf("ERROR: %d\n", http.StatusMethodNotAllowed)
				errorResponse(w, http.StatusMethodNotAllowed) // 405 ERROR
			}
		*/
		router.ServeHTTP(w, r)
	})
}
