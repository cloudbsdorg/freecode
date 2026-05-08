package controlplane

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestMemoryControlPlane_Register(t *testing.T) {
	cp := NewMemoryControlPlane().(*memoryControlPlane)
	ctx := context.Background()

	agent := Agent{
		ID:       "test-agent-1",
		Name:     "Test Agent",
		Endpoint: "http://localhost:8080",
	}

	err := cp.Register(ctx, agent)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	if len(cp.agents) != 1 {
		t.Fatalf("Expected 1 agent, got %d", len(cp.agents))
	}
}

func TestMemoryControlPlane_Unregister(t *testing.T) {
	cp := NewMemoryControlPlane().(*memoryControlPlane)
	ctx := context.Background()

	agent := Agent{
		ID:       "test-agent-1",
		Name:     "Test Agent",
		Endpoint: "http://localhost:8080",
	}

	cp.Register(ctx, agent)

	err := cp.Unregister(ctx, "test-agent-1")
	if err != nil {
		t.Fatalf("Unregister failed: %v", err)
	}

	if len(cp.agents) != 0 {
		t.Fatalf("Expected 0 agents, got %d", len(cp.agents))
	}
}

func TestMemoryControlPlane_ListAgents(t *testing.T) {
	cp := NewMemoryControlPlane().(*memoryControlPlane)
	ctx := context.Background()

	cp.Register(ctx, Agent{ID: "agent-1", Name: "Agent 1"})
	cp.Register(ctx, Agent{ID: "agent-2", Name: "Agent 2"})

	agents, err := cp.ListAgents(ctx)
	if err != nil {
		t.Fatalf("ListAgents failed: %v", err)
	}

	if len(agents) != 2 {
		t.Fatalf("Expected 2 agents, got %d", len(agents))
	}
}

func TestMemoryControlPlane_GetAgent(t *testing.T) {
	cp := NewMemoryControlPlane().(*memoryControlPlane)
	ctx := context.Background()

	cp.Register(ctx, Agent{ID: "test-agent", Name: "Test Agent"})

	agent, err := cp.GetAgent(ctx, "test-agent")
	if err != nil {
		t.Fatalf("GetAgent failed: %v", err)
	}

	if agent == nil {
		t.Fatal("Expected agent, got nil")
	}

	if agent.ID != "test-agent" {
		t.Fatalf("Expected agent ID 'test-agent', got '%s'", agent.ID)
	}
}

func TestMemoryControlPlane_GetAgent_NotFound(t *testing.T) {
	cp := NewMemoryControlPlane().(*memoryControlPlane)
	ctx := context.Background()

	agent, err := cp.GetAgent(ctx, "nonexistent")
	if err != nil {
		t.Fatalf("GetAgent failed: %v", err)
	}

	if agent != nil {
		t.Fatal("Expected nil for nonexistent agent")
	}
}

func TestMemoryControlPlane_UpdateAgentStatus(t *testing.T) {
	cp := NewMemoryControlPlane().(*memoryControlPlane)
	ctx := context.Background()

	cp.Register(ctx, Agent{ID: "test-agent", Name: "Test Agent"})

	err := cp.UpdateAgentStatus(ctx, "test-agent", "busy")
	if err != nil {
		t.Fatalf("UpdateAgentStatus failed: %v", err)
	}

	agent, _ := cp.GetAgent(ctx, "test-agent")
	if agent.Status != "busy" {
		t.Fatalf("Expected status 'busy', got '%s'", agent.Status)
	}
}

func TestMemoryControlPlane_Heartbeat(t *testing.T) {
	cp := NewMemoryControlPlane().(*memoryControlPlane)
	ctx := context.Background()

	cp.Register(ctx, Agent{ID: "test-agent", Name: "Test Agent"})

	err := cp.Heartbeat(ctx, "test-agent")
	if err != nil {
		t.Fatalf("Heartbeat failed: %v", err)
	}

	agent, _ := cp.GetAgent(ctx, "test-agent")
	if agent.Status != "online" {
		t.Fatalf("Expected status 'online' after heartbeat, got '%s'", agent.Status)
	}
}

func TestMemoryControlPlane_SubmitTask(t *testing.T) {
	cp := NewMemoryControlPlane().(*memoryControlPlane)
	ctx := context.Background()

	task := Task{
		Type:     "test-task",
		Priority: 1,
		Payload:  "test data",
	}

	id, err := cp.SubmitTask(ctx, task)
	if err != nil {
		t.Fatalf("SubmitTask failed: %v", err)
	}

	if id == "" {
		t.Fatal("Expected non-empty task ID")
	}
}

func TestMemoryControlPlane_GetTask(t *testing.T) {
	cp := NewMemoryControlPlane().(*memoryControlPlane)
	ctx := context.Background()

	cp.SubmitTask(ctx, Task{ID: "task-1", Type: "test"})

	task, err := cp.GetTask(ctx, "task-1")
	if err != nil {
		t.Fatalf("GetTask failed: %v", err)
	}

	if task == nil {
		t.Fatal("Expected task, got nil")
	}

	if task.ID != "task-1" {
		t.Fatalf("Expected task ID 'task-1', got '%s'", task.ID)
	}
}

func TestMemoryControlPlane_ListTasks(t *testing.T) {
	cp := NewMemoryControlPlane().(*memoryControlPlane)
	ctx := context.Background()

	cp.SubmitTask(ctx, Task{ID: "task-1", Type: "test"})
	cp.SubmitTask(ctx, Task{ID: "task-2", Type: "test"})

	tasks, err := cp.ListTasks(ctx, "")
	if err != nil {
		t.Fatalf("ListTasks failed: %v", err)
	}

	if len(tasks) != 2 {
		t.Fatalf("Expected 2 tasks, got %d", len(tasks))
	}
}

func TestMemoryControlPlane_ListTasks_FilterByStatus(t *testing.T) {
	cp := NewMemoryControlPlane().(*memoryControlPlane)
	ctx := context.Background()

	cp.SubmitTask(ctx, Task{ID: "task-1", Type: "test"})
	cp.SubmitTask(ctx, Task{ID: "task-2", Type: "test"})

	tasks, err := cp.ListTasks(ctx, "pending")
	if err != nil {
		t.Fatalf("ListTasks failed: %v", err)
	}

	if len(tasks) != 2 {
		t.Fatalf("Expected 2 pending tasks, got %d", len(tasks))
	}
}

func TestMemoryControlPlane_AssignTask(t *testing.T) {
	cp := NewMemoryControlPlane().(*memoryControlPlane)
	ctx := context.Background()

	cp.Register(ctx, Agent{ID: "agent-1", Name: "Agent 1"})
	cp.SubmitTask(ctx, Task{ID: "task-1", Type: "test"})

	err := cp.AssignTask(ctx, "task-1", "agent-1")
	if err != nil {
		t.Fatalf("AssignTask failed: %v", err)
	}

	task, _ := cp.GetTask(ctx, "task-1")
	if task.AgentID != "agent-1" {
		t.Fatalf("Expected agent ID 'agent-1', got '%s'", task.AgentID)
	}
	if task.Status != "assigned" {
		t.Fatalf("Expected status 'assigned', got '%s'", task.Status)
	}
}

func TestMemoryControlPlane_CompleteTask(t *testing.T) {
	cp := NewMemoryControlPlane().(*memoryControlPlane)
	ctx := context.Background()

	cp.SubmitTask(ctx, Task{ID: "task-1", Type: "test"})

	err := cp.CompleteTask(ctx, "task-1", "result data")
	if err != nil {
		t.Fatalf("CompleteTask failed: %v", err)
	}

	task, _ := cp.GetTask(ctx, "task-1")
	if task.Status != "completed" {
		t.Fatalf("Expected status 'completed', got '%s'", task.Status)
	}
	if task.Result != "result data" {
		t.Fatalf("Expected result 'result data', got '%v'", task.Result)
	}
}

func TestMemoryControlPlane_FailTask(t *testing.T) {
	cp := NewMemoryControlPlane().(*memoryControlPlane)
	ctx := context.Background()

	cp.SubmitTask(ctx, Task{ID: "task-1", Type: "test"})

	err := cp.FailTask(ctx, "task-1", nil)
	if err != nil {
		t.Fatalf("FailTask failed: %v", err)
	}

	task, _ := cp.GetTask(ctx, "task-1")
	if task.Status != "failed" {
		t.Fatalf("Expected status 'failed', got '%s'", task.Status)
	}
}

func TestMemoryControlPlane_AgentAutoID(t *testing.T) {
	cp := NewMemoryControlPlane().(*memoryControlPlane)
	ctx := context.Background()

	agent := Agent{
		Name:     "Test Agent",
		Endpoint: "http://localhost:8080",
	}

	err := cp.Register(ctx, agent)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	if len(cp.agents) != 1 {
		t.Fatalf("Expected 1 agent, got %d", len(cp.agents))
	}

	for id, a := range cp.agents {
		if id == "" {
			t.Fatal("Expected auto-generated ID, got empty string")
		}
		if a.ID != id {
			t.Fatalf("Agent ID mismatch: expected '%s', got '%s'", id, a.ID)
		}
	}
}

func TestMemoryControlPlane_TaskAutoID(t *testing.T) {
	cp := NewMemoryControlPlane().(*memoryControlPlane)
	ctx := context.Background()

	task := Task{
		Type:     "test-task",
		Priority: 1,
	}

	id, err := cp.SubmitTask(ctx, task)
	if err != nil {
		t.Fatalf("SubmitTask failed: %v", err)
	}

	if id == "" {
		t.Fatal("Expected auto-generated ID, got empty string")
	}

	storedTask, _ := cp.GetTask(ctx, id)
	if storedTask.ID != id {
		t.Fatalf("Task ID mismatch: expected '%s', got '%s'", id, storedTask.ID)
	}
}

func TestMemoryControlPlane_TaskLifecycle(t *testing.T) {
	cp := NewMemoryControlPlane().(*memoryControlPlane)
	ctx := context.Background()

	cp.Register(ctx, Agent{ID: "agent-1", Name: "Agent 1"})

	taskID, _ := cp.SubmitTask(ctx, Task{Type: "build", Priority: 5})
	_ = cp.AssignTask(ctx, taskID, "agent-1")
	_ = cp.CompleteTask(ctx, taskID, "build successful")

	task, _ := cp.GetTask(ctx, taskID)

	if task.Status != "completed" {
		t.Fatalf("Expected status 'completed', got '%s'", task.Status)
	}
	if task.AgentID != "agent-1" {
		t.Fatalf("Expected agent 'agent-1', got '%s'", task.AgentID)
	}
	if task.Completed.IsZero() {
		t.Fatal("Expected Completed timestamp to be set")
	}
}

func TestMemoryControlPlane_ConcurrentAccess(t *testing.T) {
	cp := NewMemoryControlPlane().(*memoryControlPlane)
	ctx := context.Background()

	done := make(chan bool)

	go func() {
		for i := 0; i < 50; i++ {
			cp.Register(ctx, Agent{ID: fmt.Sprintf("agent-%d", i), Name: "Agent"})
		}
		done <- true
	}()

	go func() {
		for i := 50; i < 100; i++ {
			cp.Register(ctx, Agent{ID: fmt.Sprintf("agent-%d", i), Name: "Agent"})
		}
		done <- true
	}()

	<-done
	<-done

	if len(cp.agents) != 100 {
		t.Fatalf("Expected 100 agents, got %d", len(cp.agents))
	}
}

func TestAgent_Status(t *testing.T) {
	agent := Agent{
		ID:        "test-agent",
		Name:      "Test Agent",
		Endpoint:  "http://localhost:8080",
		Status:    "online",
		Load:      0.5,
		Registered: time.Now(),
		LastSeen:  time.Now(),
	}

	if agent.Status != "online" {
		t.Fatalf("Expected status 'online', got '%s'", agent.Status)
	}

	if agent.Load != 0.5 {
		t.Fatalf("Expected load 0.5, got %f", agent.Load)
	}
}

func TestTask_Status(t *testing.T) {
	task := Task{
		ID:        "test-task",
		Type:      "build",
		Priority:  1,
		Status:    "pending",
		Created:   time.Now(),
	}

	if task.Status != "pending" {
		t.Fatalf("Expected status 'pending', got '%s'", task.Status)
	}
}