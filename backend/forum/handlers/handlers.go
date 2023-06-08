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
func (h *Handlers) Handler() http.Handler {

	var routes = []server.Route{
		// method GET endpoints
		server.NewRoute(http.MethodGet, `/api/me`, h.me),

		server.NewRoute(http.MethodGet, `/api/me/posts`, h.mePosts),
		server.NewRoute(http.MethodGet, `/api/me/posts/liked`, h.mePostsLiked),


		server.NewRoute(http.MethodGet, `/api/users`, h.userAll),
		server.NewRoute(http.MethodGet, `/api/users/(\d+)`, h.userId),
		server.NewRoute(http.MethodGet, `/api/users/(\d+)/posts`, h.userIdPosts),

		server.NewRoute(http.MethodGet, `/api/posts`, h.posts),
		server.NewRoute(http.MethodGet, `/api/posts/(\d+)`, h.postsId),

		server.NewRoute(http.MethodGet, `/api/posts/(\d+)/reaction`, h.postsIdReaction),
		server.NewRoute(http.MethodGet, `/api/posts/(\d+)/comments`, h.postsIdComments),

		server.NewRoute(http.MethodGet, `/api/posts/(\d+)/comment/(\d+)/reaction`, h.postsIdCommentIdReaction),

		server.NewRoute(http.MethodGet, `/api/posts/categories`, h.postsCategories),
		server.NewRoute(http.MethodGet, `/api/posts/categories/([[:alnum:]]+)`, h.postsCategoriesName),

		server.NewRoute(http.MethodGet, `/api/internal/check-session`, h.checkSession),

		// method POST endpoints
		server.NewRoute(http.MethodPost, `/api/posts/create`, h.postsCreate),
		server.NewRoute(http.MethodPost, `/api/posts/(\d+)/like`, h.postsIdLike),
		server.NewRoute(http.MethodPost, `/api/posts/(\d+)/dislike`, h.postsIdDislike),

		server.NewRoute(http.MethodPost, `/api/login`, h.login),
		server.NewRoute(http.MethodPost, `/api/signup`, h.signup),
		server.NewRoute(http.MethodPost, `/api/logout`, h.logout),

		server.NewRoute(http.MethodPost, `/api/posts/(\d+)/comment`, h.postsIdCommentCreate),

		server.NewRoute(http.MethodPost, `/api/posts/(\d+)/comment/(\d+)/like`, h.postsIdCommentIdLike),
		server.NewRoute(http.MethodPost, `/api/posts/(\d+)/comment/(\d+)/dislike`, h.postsIdCommentIdDislike),
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
