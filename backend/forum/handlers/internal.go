package handlers

import (
	"backend/common/server"
	"log"
	"net/http"
)

// checkSession checks if the session token is valid
//
// GET /api/internal/check-session?token=<token>
func (h *Handlers) checkSession(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		server.ErrorResponse(w, http.StatusBadRequest)
		return
	}
	userId := h.DB.CheckSession(token)
	if userId == -1 {
		server.SendObject(w, struct {
			Error string `json:"error"`
		}{
			Error: "Invalid token",
		})
		return
	}
	server.SendObject(w, userId)
}

// events returns an SSE stream of revoked sessions
func (h *Handlers) events(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	flusher.Flush()

	subscriber := make(chan string)

	h.lock.Lock()
	h.revokeSessionSubscribers = append(h.revokeSessionSubscribers, &subscriber)
	h.lock.Unlock()

loop:
	for {
		select {
		case token := <-subscriber:
			_, err := w.Write([]byte(token + "\n"))
			if err != nil {
				log.Println(err)
				break loop
			}
			flusher.Flush()
		case <-r.Context().Done():
			break loop
		}
	}

	for i, s := range h.revokeSessionSubscribers {
		if s == &subscriber {
			h.lock.Lock()
			h.revokeSessionSubscribers = append(h.revokeSessionSubscribers[:i], h.revokeSessionSubscribers[i+1:]...)
			h.lock.Unlock()
			close(subscriber)
			break
		}
	}
}
