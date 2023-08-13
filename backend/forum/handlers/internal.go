package handlers

import (
	"backend/common/integrations/auth"
	"backend/common/server"
	"encoding/json"
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

// events is an SSE stream of auth.Event objects
func (h *Handlers) events(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	flusher.Flush()

	subscriber := make(chan auth.Event)

	h.lock.Lock()
	h.eventSubscribers = append(h.eventSubscribers, &subscriber)
	h.lock.Unlock()

loop:
	for {
		select {
		case event := <-subscriber:
			eventJSON, err := json.Marshal(event)
			if err != nil {
				log.Println(err)
				break loop
			}

			_, err = w.Write(append(eventJSON, '\n'))
			if err != nil {
				log.Println(err)
				break loop
			}
			flusher.Flush()
		case <-r.Context().Done():
			break loop
		}
	}

	for i, s := range h.eventSubscribers {
		if s == &subscriber {
			h.lock.Lock()
			h.eventSubscribers = append(h.eventSubscribers[:i], h.eventSubscribers[i+1:]...)
			h.lock.Unlock()
			close(subscriber)
			break
		}
	}
}
