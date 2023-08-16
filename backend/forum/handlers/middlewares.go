package handlers

import (
	"backend/common/server"
	"context"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// withAuth is a middleware that checks if the user is authenticated
func (h *Handlers) withAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userId int

		if r.Header.Get("Internal-Auth") == "" {
			userId = h.getUserId(w, r)
			if userId == -1 {
				server.ErrorResponse(w, http.StatusUnauthorized)
				return
			}
		} else {
			// Bypass auth if Internal-Auth header is set
			if r.Header.Get("Internal-Auth") == os.Getenv("FORUM_BACKEND_SECRET") ||
				os.Getenv("FORUM_BACKEND_SECRET") == "" || r.Header.Get("Internal-Auth") == "SSR" {
				userId = -1
			} else {
				server.ErrorResponse(w, http.StatusUnauthorized)
				return
			}
		}

		ctx := context.WithValue(r.Context(), userIdCtxKey, userId)
		r = r.WithContext(ctx)

		handler(w, r)
	}
}

// withInternal is a middleware that checks if the request has the valid Internal-Auth header set
func (h *Handlers) withInternal(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Internal-Auth") == "" {
			server.ErrorResponse(w, http.StatusForbidden)
			return
		}

		if r.Header.Get("Internal-Auth") == os.Getenv("FORUM_BACKEND_SECRET") ||
			os.Getenv("FORUM_BACKEND_SECRET") == "" {
			handler(w, r)
		} else {
			server.ErrorResponse(w, http.StatusUnauthorized)
		}
	}
}

// withPermissions is a middleware that checks if the user has permissions
func (h *Handlers) withPermissions(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// r.URL.Path does not match any posts endpoint, therefore it skips it
		if !regexp.MustCompile(`^/api/posts/\d+`).MatchString(r.URL.Path) {
			handler(w, r)
			return
		}

		postId, err := strconv.Atoi(strings.Split(r.URL.Path, "/")[3])
		if err != nil {
			server.ErrorResponse(w, http.StatusNotFound)
			return
		}

		if post := h.DB.GetPostById(postId); post == nil {
			server.ErrorResponse(w, http.StatusNotFound)
			return
		}

		userId := h.getUserId(w, r)

		ctx := context.WithValue(r.Context(), postIdCtxKey, postId)

		r = r.WithContext(ctx)

		// Forum is private, server side rendering.
		if os.Getenv("FORUM_IS_PRIVATE") == isPrivate && userId == -1 {
			r.Header.Set("Internal-Auth", "SSR")
			handler(w, r)
			return
		}

		// Forum is public, user is not logged in, only show public non-group posts
		if os.Getenv("FORUM_IS_PRIVATE") != isPrivate && userId == -1 {
			r.Header.Set("Internal-Auth", "Public")
			handler(w, r)
			return
		}

		// User does not have post permissions, get status forbidden
		if !h.DB.GetPostPermissions(userId, postId) {
			server.ErrorResponse(w, http.StatusForbidden)
			return
		}

		// user has permissions
		handler(w, r)
	}
}
