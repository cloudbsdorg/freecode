package api

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	tools map[string]ToolHandler
}

type ToolHandler func(args map[string]interface{}) (interface{}, error)

func NewHandler() *Handler {
	return &Handler{
		tools: make(map[string]ToolHandler),
	}
}

func (h *Handler) RegisterTool(name string, desc string, handler ToolHandler) {
	h.tools[name] = handler
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	resp := h.handle(&req)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) handle(req *Request) *Response {
	switch req.Method {
	case "initialize":
		return &Response{
			JSONRPC: "2.0",
			Result: map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"capabilities": map[string]interface{}{
					"tools":     true,
					"resources": true,
				},
				"serverInfo": map[string]interface{}{
					"name":    "freecode-api",
					"version": "0.1.0",
				},
			},
			ID: req.ID,
		}

	case "tools/list":
		var tools []Tool
		for name, handler := range h.tools {
			tools = append(tools, Tool{
				Name:        name,
				Description: "Tool: " + name,
			})
			_ = handler
		}
		return &Response{
			JSONRPC: "2.0",
			Result: map[string]interface{}{
				"tools": tools,
			},
			ID: req.ID,
		}

	default:
		return &Response{
			JSONRPC: "2.0",
			Error: &Error{
				Code:    -32601,
				Message: "Method not found: " + req.Method,
			},
			ID: req.ID,
		}
	}
}
