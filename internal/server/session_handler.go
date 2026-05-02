package server

import (
	"encoding/json"
	"net/http"

	"github.com/freecode/freecode/internal/session"
)

type SessionHandler struct {
	manager *session.Manager
}

func NewSessionHandler(manager *session.Manager) *SessionHandler {
	return &SessionHandler{
		manager: manager,
	}
}

func (h *SessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		sessions := h.manager.ListSessions()
		json.NewEncoder(w).Encode(sessions)

	case "POST":
		var req struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		sess, err := h.manager.CreateSession(req.Name, "", "")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(sess)
	}
}
