package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/freecode/freecode/internal/agent"
	"github.com/freecode/freecode/internal/session"
)

type AgentHandler struct {
	engine  *agent.Engine
	session *session.Manager
}

func NewAgentHandler(engine *agent.Engine, session *session.Manager) *AgentHandler {
	return &AgentHandler{
		engine:  engine,
		session: session,
	}
}

type RunRequest struct {
	SessionID string            `json:"session_id"`
	Prompt   string             `json:"prompt"`
	Agent    string             `json:"agent"`
	Tools    []string           `json:"tools"`
}

func (h *AgentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "POST":
		var req RunRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		ctx := context.Background()
		result, err := h.engine.Run(ctx, agent.Request{
			SessionID: req.SessionID,
			AgentName: req.Agent,
			Message: agent.Message{
				Role:    "user",
				Content: req.Prompt,
			},
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(result)
	}
}
