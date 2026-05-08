package controlplane

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/google/uuid"
)

// Agent represents a fleet agent instance
type Agent struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Endpoint  string    `json:"endpoint"`
	Status    string    `json:"status"`
	Metadata  Metadata  `json:"metadata,omitempty"`
	Registered time.Time `json:"registered"`
	LastSeen  time.Time `json:"last_seen"`
	Capabilities []string `json:"capabilities,omitempty"`
	Load      float64   `json:"load"`
}

// Metadata contains agent metadata
type Metadata struct {
	OS      string `json:"os,omitempty"`
	Arch    string `json:"arch,omitempty"`
	Version string `json:"version,omitempty"`
}

// Task represents a distributed task
type Task struct {
	ID        string       `json:"id"`
	Type      string       `json:"type"`
	Payload   interface{}  `json:"payload,omitempty"`
	Priority  int          `json:"priority"`
	AgentID   string       `json:"agent_id,omitempty"`
	Status    string       `json:"status"`
	Created   time.Time    `json:"created"`
	Started   time.Time    `json:"started,omitempty"`
	Completed time.Time    `json:"completed,omitempty"`
	Result    interface{}  `json:"result,omitempty"`
	Error     string       `json:"error,omitempty"`
}

// ControlPlane is the fleet orchestration interface
type ControlPlane interface {
	// Agent management
	Register(ctx context.Context, agent Agent) error
	Unregister(ctx context.Context, id string) error
	ListAgents(ctx context.Context) ([]Agent, error)
	GetAgent(ctx context.Context, id string) (*Agent, error)
	UpdateAgentStatus(ctx context.Context, id string, status string) error
	Heartbeat(ctx context.Context, id string) error

	// Task management
	SubmitTask(ctx context.Context, task Task) (string, error)
	GetTask(ctx context.Context, id string) (*Task, error)
	ListTasks(ctx context.Context, status string) ([]Task, error)
	AssignTask(ctx context.Context, taskID, agentID string) error
	CompleteTask(ctx context.Context, id string, result interface{}) error
	FailTask(ctx context.Context, id string, err error) error

	// Fleet head server
	StartServer(ctx context.Context, addr string) error
	StopServer() error
}

type memoryControlPlane struct {
	mu      sync.RWMutex
	agents  map[string]Agent
	tasks   map[string]Task
	server  *http.Server
	upgrader websocket.Upgrader
	conns    map[string]*websocket.Conn
}

// NewMemoryControlPlane creates a new in-memory control plane
func NewMemoryControlPlane() ControlPlane {
	return &memoryControlPlane{
		agents:  make(map[string]Agent),
		tasks:   make(map[string]Task),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // TODO: proper origin check
			},
		},
		conns: make(map[string]*websocket.Conn),
	}
}

func (cp *memoryControlPlane) Register(ctx context.Context, agent Agent) error {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	if agent.ID == "" {
		agent.ID = uuid.New().String()
	}
	agent.Registered = time.Now()
	agent.LastSeen = time.Now()
	agent.Status = "online"

	cp.agents[agent.ID] = agent
	return nil
}

func (cp *memoryControlPlane) Unregister(ctx context.Context, id string) error {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	delete(cp.agents, id)
	delete(cp.tasks, id)

	// Close websocket connection if exists
	if conn, ok := cp.conns[id]; ok {
		conn.Close()
		delete(cp.conns, id)
	}

	return nil
}

func (cp *memoryControlPlane) ListAgents(ctx context.Context) ([]Agent, error) {
	cp.mu.RLock()
	defer cp.mu.RUnlock()

	var result []Agent
	for _, a := range cp.agents {
		result = append(result, a)
	}
	return result, nil
}

func (cp *memoryControlPlane) GetAgent(ctx context.Context, id string) (*Agent, error) {
	cp.mu.RLock()
	defer cp.mu.RUnlock()

	if a, ok := cp.agents[id]; ok {
		return &a, nil
	}
	return nil, nil
}

func (cp *memoryControlPlane) UpdateAgentStatus(ctx context.Context, id string, status string) error {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	if a, ok := cp.agents[id]; ok {
		a.Status = status
		a.LastSeen = time.Now()
		cp.agents[id] = a
		return nil
	}
	return fmt.Errorf("agent not found: %s", id)
}

func (cp *memoryControlPlane) Heartbeat(ctx context.Context, id string) error {
	return cp.UpdateAgentStatus(ctx, id, "online")
}

func (cp *memoryControlPlane) SubmitTask(ctx context.Context, task Task) (string, error) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	if task.ID == "" {
		task.ID = uuid.New().String()
	}
	task.Status = "pending"
	task.Created = time.Now()

	cp.tasks[task.ID] = task
	return task.ID, nil
}

func (cp *memoryControlPlane) GetTask(ctx context.Context, id string) (*Task, error) {
	cp.mu.RLock()
	defer cp.mu.RUnlock()

	if t, ok := cp.tasks[id]; ok {
		return &t, nil
	}
	return nil, nil
}

func (cp *memoryControlPlane) ListTasks(ctx context.Context, status string) ([]Task, error) {
	cp.mu.RLock()
	defer cp.mu.RUnlock()

	var result []Task
	for _, t := range cp.tasks {
		if status == "" || t.Status == status {
			result = append(result, t)
		}
	}
	return result, nil
}

func (cp *memoryControlPlane) AssignTask(ctx context.Context, taskID, agentID string) error {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	task, ok := cp.tasks[taskID]
	if !ok {
		return fmt.Errorf("task not found: %s", taskID)
	}

	if _, ok := cp.agents[agentID]; !ok {
		return fmt.Errorf("agent not found: %s", agentID)
	}

	task.AgentID = agentID
	task.Status = "assigned"
	task.Started = time.Now()
	cp.tasks[taskID] = task

	return nil
}

func (cp *memoryControlPlane) CompleteTask(ctx context.Context, id string, result interface{}) error {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	task, ok := cp.tasks[id]
	if !ok {
		return fmt.Errorf("task not found: %s", id)
	}

	task.Status = "completed"
	task.Completed = time.Now()
	task.Result = result
	cp.tasks[id] = task

	return nil
}

func (cp *memoryControlPlane) FailTask(ctx context.Context, id string, err error) error {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	task, ok := cp.tasks[id]
	if !ok {
		return fmt.Errorf("task not found: %s", id)
	}

	task.Status = "failed"
	task.Completed = time.Now()
	if err != nil {
		task.Error = err.Error()
	}
	cp.tasks[id] = task

	return nil
}

// StartServer starts the fleet head HTTP+WebSocket server
func (cp *memoryControlPlane) StartServer(ctx context.Context, addr string) error {
	mux := http.NewServeMux()

	// REST endpoints
	mux.HandleFunc("POST /api/agents", cp.handleRegisterAgent)
	mux.HandleFunc("DELETE /api/agents/{id}", cp.handleUnregisterAgent)
	mux.HandleFunc("GET /api/agents", cp.handleListAgents)
	mux.HandleFunc("GET /api/agents/{id}", cp.handleGetAgent)
	mux.HandleFunc("POST /api/agents/{id}/heartbeat", cp.handleHeartbeat)

	mux.HandleFunc("POST /api/tasks", cp.handleSubmitTask)
	mux.HandleFunc("GET /api/tasks", cp.handleListTasks)
	mux.HandleFunc("GET /api/tasks/{id}", cp.handleGetTask)
	mux.HandleFunc("POST /api/tasks/{id}/assign", cp.handleAssignTask)
	mux.HandleFunc("POST /api/tasks/{id}/complete", cp.handleCompleteTask)
	mux.HandleFunc("POST /api/tasks/{id}/fail", cp.handleFailTask)

	// WebSocket endpoint for agent connections
	mux.HandleFunc("WS /ws/agents/{id}", cp.handleWebSocket)

	cp.server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		cp.server.Shutdown(context.Background())
	}()

	return cp.server.ListenAndServe()
}

// StopServer stops the fleet head server
func (cp *memoryControlPlane) StopServer() error {
	if cp.server != nil {
		return cp.server.Shutdown(context.Background())
	}
	return nil
}

// HTTP handlers

func (cp *memoryControlPlane) handleRegisterAgent(w http.ResponseWriter, r *http.Request) {
	var agent Agent
	if err := json.NewDecoder(r.Body).Decode(&agent); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if err := cp.Register(r.Context(), agent); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(agent)
}

func (cp *memoryControlPlane) handleUnregisterAgent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := cp.Unregister(r.Context(), id); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(204)
}

func (cp *memoryControlPlane) handleListAgents(w http.ResponseWriter, r *http.Request) {
	agents, err := cp.ListAgents(r.Context())
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(agents)
}

func (cp *memoryControlPlane) handleGetAgent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	agent, err := cp.GetAgent(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if agent == nil {
		http.Error(w, "agent not found", 404)
		return
	}
	json.NewEncoder(w).Encode(agent)
}

func (cp *memoryControlPlane) handleHeartbeat(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := cp.Heartbeat(r.Context(), id); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(204)
}

func (cp *memoryControlPlane) handleSubmitTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	id, err := cp.SubmitTask(r.Context(), task)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write([]byte(id))
}

func (cp *memoryControlPlane) handleListTasks(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	tasks, err := cp.ListTasks(r.Context(), status)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(tasks)
}

func (cp *memoryControlPlane) handleGetTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	task, err := cp.GetTask(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if task == nil {
		http.Error(w, "task not found", 404)
		return
	}
	json.NewEncoder(w).Encode(task)
}

func (cp *memoryControlPlane) handleAssignTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req struct {
		AgentID string `json:"agent_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if err := cp.AssignTask(r.Context(), id, req.AgentID); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(204)
}

func (cp *memoryControlPlane) handleCompleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req struct {
		Result interface{} `json:"result"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if err := cp.CompleteTask(r.Context(), id, req.Result); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(204)
}

func (cp *memoryControlPlane) handleFailTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req struct {
		Error string `json:"error"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if err := cp.FailTask(r.Context(), id, fmt.Errorf(req.Error)); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(204)
}

func (cp *memoryControlPlane) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	conn, err := cp.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer conn.Close()

	cp.mu.Lock()
	cp.conns[id] = conn
	cp.mu.Unlock()

	// Keep connection alive and handle messages
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		// Handle incoming messages (tasks, heartbeats, results)
		var message WebSocketMessage
		if err := json.Unmarshal(msg, &message); err != nil {
			continue
		}

		cp.handleWebSocketMessage(r.Context(), id, message)
	}

	cp.mu.Lock()
	delete(cp.conns, id)
	cp.mu.Unlock()
}

// WebSocketMessage represents a WebSocket message
type WebSocketMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload,omitempty"`
}

func (cp *memoryControlPlane) handleWebSocketMessage(ctx context.Context, agentID string, msg WebSocketMessage) {
	switch msg.Type {
	case "heartbeat":
		cp.Heartbeat(ctx, agentID)
	case "task_result":
		// Handle task result from agent
		if payload, ok := msg.Payload.(map[string]interface{}); ok {
			if taskID, ok := payload["task_id"].(string); ok {
				if result, ok := payload["result"]; ok {
					cp.CompleteTask(ctx, taskID, result)
				}
			}
		}
	case "status_update":
		// Handle agent status update
		if payload, ok := msg.Payload.(map[string]interface{}); ok {
			if status, ok := payload["status"].(string); ok {
				cp.UpdateAgentStatus(ctx, agentID, status)
			}
			if load, ok := payload["load"].(float64); ok {
				cp.mu.Lock()
				if a, ok := cp.agents[agentID]; ok {
					a.Load = load
					cp.agents[agentID] = a
				}
				cp.mu.Unlock()
			}
		}
	}
}

// BroadcastTask sends a task to a specific agent via WebSocket
func (cp *memoryControlPlane) BroadcastTask(agentID string, task Task) error {
	cp.mu.RLock()
	conn, ok := cp.conns[agentID]
	cp.mu.RUnlock()

	if !ok {
		return fmt.Errorf("agent not connected: %s", agentID)
	}

	msg := WebSocketMessage{
		Type:    "task",
		Payload: task,
	}

	return conn.WriteJSON(msg)
}