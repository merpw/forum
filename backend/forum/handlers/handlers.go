package handlers

import (
	"backend/common/server"
	"backend/forum/database"
	"database/sql"
	"log"
	"net/http"
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

	var routes = []server.Route{
		// method GET endpoints
		server.NewRoute(http.MethodGet, `/api/me`, handlers.me),

		server.NewRoute(http.MethodGet, `/api/me/posts`, handlers.mePosts),
		server.NewRoute(http.MethodGet, `/api/me/posts/liked`, handlers.mePostsLiked),

		server.NewRoute(http.MethodGet, `/api/user/(\d+)`, handlers.userId),
		server.NewRoute(http.MethodGet, `/api/user/(\d+)/posts`, handlers.userIdPosts),

		server.NewRoute(http.MethodGet, `/api/posts`, handlers.posts),
		server.NewRoute(http.MethodGet, `/api/posts/(\d+)`, handlers.postsId),

		server.NewRoute(http.MethodGet, `/api/posts/(\d+)/reaction`, handlers.postsIdReaction),
		server.NewRoute(http.MethodGet, `/api/posts/(\d+)/comments`, handlers.postsIdComments),

		server.NewRoute(http.MethodGet, `/api/posts/(\d+)/comment/(\d+)/reaction`, handlers.postsIdCommentIdReaction),

		server.NewRoute(http.MethodGet, `/api/posts/categories`, handlers.postsCategories),
		server.NewRoute(http.MethodGet, `/api/posts/categories/([[:alnum:]]+)`, handlers.postsCategoriesName),

		server.NewRoute(http.MethodGet, `/api/internal/check-session`, handlers.checkSession),

		// method POST endpoints
		server.NewRoute(http.MethodPost, `/api/posts/create`, handlers.postsCreate),
		server.NewRoute(http.MethodPost, `/api/posts/(\d+)/like`, handlers.postsIdLike),
		server.NewRoute(http.MethodPost, `/api/posts/(\d+)/dislike`, handlers.postsIdDislike),

		server.NewRoute(http.MethodPost, `/api/login`, handlers.login),
		server.NewRoute(http.MethodPost, `/api/signup`, handlers.signup),
		server.NewRoute(http.MethodPost, `/api/logout`, handlers.logout),

		server.NewRoute(http.MethodPost, `/api/posts/(\d+)/comment`, handlers.postsIdCommentCreate),

		server.NewRoute(http.MethodPost, `/api/posts/(\d+)/comment/(\d+)/like`, handlers.postsIdCommentIdLike),
		server.NewRoute(http.MethodPost, `/api/posts/(\d+)/comment/(\d+)/dislike`, handlers.postsIdCommentIdDislike),
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("ERROR %d. %v\n", http.StatusInternalServerError, err)
				server.ErrorResponse(w, http.StatusInternalServerError) // 500 ERROR
			}
		}()
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")

		var requestRoute server.Route

		for _, route := range routes {
			if route.Pattern.MatchString(r.URL.Path) {
				requestRoute = route
				break
			}
		}

		if requestRoute.Handler == nil {
			server.ErrorResponse(w, http.StatusNotFound)
			return
		}

		if r.Method != requestRoute.Method {
			server.ErrorResponse(w, http.StatusMethodNotAllowed)
			return
		}

		requestRoute.Handler(w, r)
	})
}
