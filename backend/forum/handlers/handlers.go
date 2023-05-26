package handlers

import (
	"backend/forum/database"
	"database/sql"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type Handlers struct {
	DB *database.DB
}

// New connects database to Handlers
func New(db *sql.DB) *Handlers {
	return &Handlers{DB: database.New(db)}
}

// Handler returns http.Handler with all routes registered
func (handlers *Handlers) Handler() http.Handler {

	var routes = []Route{
		// method GET endpoints
		newRoute(http.MethodGet, `/api/me`, handlers.me),

		newRoute(http.MethodGet, `/api/me/posts`, handlers.mePosts),
		newRoute(http.MethodGet, `/api/me/posts/liked`, handlers.mePostsLiked),

		newRoute(http.MethodGet, `/api/user/(\d+)`, handlers.userId),
		newRoute(http.MethodGet, `/api/user/(\d+)/posts`, handlers.userIdPosts),

		newRoute(http.MethodGet, `/api/posts`, handlers.posts),
		newRoute(http.MethodGet, `/api/posts/(\d+)`, handlers.postsId),

		newRoute(http.MethodGet, `/api/posts/(\d+)/reaction`, handlers.postsIdReaction),
		newRoute(http.MethodGet, `/api/posts/(\d+)/comments`, handlers.postsIdComments),

		newRoute(http.MethodGet, `/api/posts/(\d+)/comment/(\d+)/reaction`, handlers.postsIdCommentIdReaction),

		newRoute(http.MethodGet, `/api/posts/categories`, handlers.postsCategories),
		newRoute(http.MethodGet, `/api/posts/categories/([[:alnum:]]+)`, handlers.postsCategoriesName),

		newRoute(http.MethodGet, `/api/internal/check-session`, handlers.checkSession),

		// method POST endpoints
		newRoute(http.MethodPost, `/api/posts/create`, handlers.postsCreate),
		newRoute(http.MethodPost, `/api/posts/(\d+)/like`, handlers.postsIdLike),
		newRoute(http.MethodPost, `/api/posts/(\d+)/dislike`, handlers.postsIdDislike),

		newRoute(http.MethodPost, `/api/login`, handlers.login),
		newRoute(http.MethodPost, `/api/signup`, handlers.signup),
		newRoute(http.MethodPost, `/api/logout`, handlers.logout),

		newRoute(http.MethodPost, `/api/posts/(\d+)/comment`, handlers.postsIdCommentCreate),

		newRoute(http.MethodPost, `/api/posts/(\d+)/comment/(\d+)/like`, handlers.postsIdCommentIdLike),
		newRoute(http.MethodPost, `/api/posts/(\d+)/comment/(\d+)/dislike`, handlers.postsIdCommentIdDislike),
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("ERROR %d. %v\n", http.StatusInternalServerError, err)
				errorResponse(w, http.StatusInternalServerError) // 500 ERROR
			}
		}()
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")

		var requestRoute Route

		for _, route := range routes {
			if route.Pattern.MatchString(r.URL.Path) {
				requestRoute = route
				break
			}
		}

		if requestRoute.Handler == nil {
			errorResponse(w, http.StatusNotFound)
			return
		}

		if r.Method != requestRoute.Method {
			errorResponse(w, http.StatusMethodNotAllowed)
			return
		}

		requestRoute.Handler(w, r)
	})
}

type Route struct {
	Method  string
	Pattern *regexp.Regexp
	Handler http.HandlerFunc
}

func newRoute(method, pattern string, handler http.HandlerFunc) Route {
	return Route{method, regexp.MustCompile("^" + pattern + "$"), handler}
}
