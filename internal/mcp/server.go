package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Server struct {
	Port     int
	Tools    []Tool
	Handlers map[string]func(args map[string]interface{}) (interface{}, error)
}

type ServerRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string           `json:"method"`
	Params  json.RawMessage  `json:"params,omitempty"`
	ID     interface{}      `json:"id,omitempty"`
}

type ServerResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error  *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
	ID     interface{} `json:"id,omitempty"`
}

func NewServer(port int) *Server {
	return &Server{
		Port:     port,
		Tools:    make([]Tool, 0),
		Handlers: make(map[string]func(args map[string]interface{}) (interface{}, error)),
	}
}

func (s *Server) RegisterTool(name, description string, handler func(args map[string]interface{}) (interface{}, error)) {
	s.Tools = append(s.Tools, Tool{
		Name:        name,
		Description: description,
		InputSchema: make(map[string]interface{}),
	})
	s.Handlers[name] = handler
}

func (s *Server) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ServerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	resp := s.processRequest(&req)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) processRequest(req *ServerRequest) *ServerResponse {
	switch req.Method {
	case "initialize":
		return &ServerResponse{
			JSONRPC: "2.0",
			Result: map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"capabilities": map[string]interface{}{
					"tools": true,
				},
				"serverInfo": map[string]interface{}{
					"name":    "freecode-mcp",
					"version": "0.1.0",
				},
			},
			ID: req.ID,
		}

	case "tools/list":
		return &ServerResponse{
			JSONRPC: "2.0",
			Result: map[string]interface{}{
				"tools": s.Tools,
			},
			ID: req.ID,
		}

	case "tools/call":
		var params struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments,omitempty"`
		}
		if err := json.Unmarshal(req.Params, &params); err != nil {
			return &ServerResponse{
				JSONRPC: "2.0",
				Error: &struct {
					Code    int    `json:"code"`
					Message string `json:"message"`
				}{-32602, "Invalid params"},
				ID: req.ID,
			}
		}

		handler, ok := s.Handlers[params.Name]
		if !ok {
			return &ServerResponse{
				JSONRPC: "2.0",
				Error: &struct {
					Code    int    `json:"code"`
					Message string `json:"message"`
				}{-32602, fmt.Sprintf("Unknown tool: %s", params.Name)},
				ID: req.ID,
			}
		}

		result, err := handler(params.Arguments)
		if err != nil {
			return &ServerResponse{
				JSONRPC: "2.0",
				Result: map[string]interface{}{
					"content": []map[string]string{
						{"type": "text", "text": err.Error()},
					},
					"isError": true,
				},
				ID: req.ID,
			}
		}

		return &ServerResponse{
			JSONRPC: "2.0",
			Result: map[string]interface{}{
				"content": []map[string]string{
					{"type": "text", "text": fmt.Sprintf("%v", result)},
				},
			},
			ID: req.ID,
		}
	}

	return &ServerResponse{
		JSONRPC: "2.0",
		Error: &struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}{-32601, fmt.Sprintf("Method not found: %s", req.Method)},
		ID: req.ID,
	}
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.Handle)

	addr := fmt.Sprintf("127.0.0.1:%d", s.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		server.Close()
	}()

	return server.ListenAndServe()
}
