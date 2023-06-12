package handlers

import (
	"backend/common/server"
	"context"
	"net/http"
	"os"
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
				os.Getenv("FORUM_BACKEND_SECRET") == "" {
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
