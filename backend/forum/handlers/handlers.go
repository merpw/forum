package handlers

import (
	"backend/common/server"
	"backend/forum/database"
	"database/sql"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"sync"
)

type ctxKey int

const (
	userIdCtxKey ctxKey = iota
)

type Handlers struct {
	DB *database.DB

	revokeSession            chan string
	revokeSessionSubscribers []*chan string

	lock sync.Mutex
}

// New connects database to Handlers
func New(db *sql.DB) *Handlers {

	h := &Handlers{DB: database.New(db)}

	h.revokeSession = make(chan string)
	h.revokeSessionSubscribers = make([]*chan string, 0)

	go func() {
		for token := range h.revokeSession {
			log.Printf("Revoked session: %s %v", token, len(h.revokeSessionSubscribers))
			for _, subscriber := range h.revokeSessionSubscribers {
				*subscriber <- token
			}
		}
	}()

	return h
}

// Handler returns http.Handler with all routes registered
func (h *Handlers) Handler() http.Handler {

	var routes = []server.Route{
		// method GET endpoints
		server.NewRoute(http.MethodGet, `/api/users/(\d+)`, h.usersId),
		server.NewRoute(http.MethodGet, `/api/users/(\d+)/posts`, h.usersIdPosts),
		server.NewRoute(http.MethodGet, `/api/users`, h.usersAll),

		server.NewRoute(http.MethodGet, `/api/posts`, h.posts),
		server.NewRoute(http.MethodGet, `/api/posts/(\d+)`, h.postsId),

		server.NewRoute(http.MethodGet, `/api/posts/(\d+)/comments`, h.postsIdComments),

		server.NewRoute(http.MethodGet, `/api/posts/categories`, h.postsCategories),
		server.NewRoute(http.MethodGet, `/api/posts/categories/([[:alnum:]]+)`, h.postsCategoriesName),
	}

	var publicRoutes = []server.Route{
		// method POST endpoints
		server.NewRoute(http.MethodPost, `/api/login`, h.login),
		server.NewRoute(http.MethodPost, `/api/signup`, h.signup),
	}

	var authRoutes = []server.Route{
		// method GET endpoints
		server.NewRoute(http.MethodGet, `/api/me`, h.me),

		server.NewRoute(http.MethodGet, `/api/me/posts`, h.mePosts),
		server.NewRoute(http.MethodGet, `/api/me/posts/liked`, h.mePostsLiked),

		server.NewRoute(http.MethodGet, `/api/posts/(\d+)/reaction`, h.postsIdReaction),

		server.NewRoute(http.MethodGet, `/api/posts/(\d+)/comment/(\d+)/reaction`, h.postsIdCommentIdReaction),

		// method POST endpoints
		server.NewRoute(http.MethodPost, `/api/logout`, h.logout),

		server.NewRoute(http.MethodPost, `/api/posts/create`, h.postsCreate),
		server.NewRoute(http.MethodPost, `/api/posts/(\d+)/like`, h.postsIdLike),
		server.NewRoute(http.MethodPost, `/api/posts/(\d+)/dislike`, h.postsIdDislike),

		server.NewRoute(http.MethodPost, `/api/posts/(\d+)/comment`, h.postsIdCommentCreate),

		server.NewRoute(http.MethodPost, `/api/posts/(\d+)/comment/(\d+)/like`, h.postsIdCommentIdLike),
		server.NewRoute(http.MethodPost, `/api/posts/(\d+)/comment/(\d+)/dislike`, h.postsIdCommentIdDislike),
	}

	var internalRoutes = []server.Route{
		server.NewRoute(http.MethodGet, `/api/internal/check-session`, h.checkSession),
		server.NewRoute(http.MethodGet, `/api/internal/revoked-sessions`, h.revokedSessions),
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("ERROR %d. %v\n%s", http.StatusInternalServerError, err, debug.Stack())
				server.ErrorResponse(w, http.StatusInternalServerError) // 500 ERROR
			}
		}()

		// remove trailing slash
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")

		for _, route := range routes {
			if route.Pattern.MatchString(r.URL.Path) {
				if r.Method != route.Method {
					server.ErrorResponse(w, http.StatusMethodNotAllowed)
					return
				}

				if os.Getenv("FORUM_IS_PRIVATE") == "true" {
					h.withAuth(route.Handler)(w, r)
					return
				}

				route.Handler(w, r)
				return
			}
		}

		for _, route := range publicRoutes {
			if route.Pattern.MatchString(r.URL.Path) {
				if r.Method != route.Method {
					server.ErrorResponse(w, http.StatusMethodNotAllowed)
					return
				}

				route.Handler(w, r)
				return
			}
		}

		for _, route := range authRoutes {
			if route.Pattern.MatchString(r.URL.Path) {
				if r.Method != route.Method {
					server.ErrorResponse(w, http.StatusMethodNotAllowed)
					return
				}

				h.withAuth(route.Handler)(w, r)
				return
			}
		}

		for _, route := range internalRoutes {
			if route.Pattern.MatchString(r.URL.Path) {
				if r.Method != route.Method {
					server.ErrorResponse(w, http.StatusMethodNotAllowed)
					return
				}

				h.withInternal(route.Handler)(w, r)
				return
			}
		}

		// if we're still here, no route was found
		server.ErrorResponse(w, http.StatusNotFound)
	})
}
