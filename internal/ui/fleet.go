package ui

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/freecode/freecode/internal/controlplane"
)

type FleetAgent struct {
	ID        string
	Name      string
	Status    string
	Load      float64
	Connected time.Time
}

type FleetTask struct {
	ID        string
	Type      string
	Status    string
	AgentID   string
	Created   time.Time
	Completed time.Time
}

type FleetPanel struct {
	width     int
	height    int
	isOpen    bool
	agents    []FleetAgent
	tasks     []FleetTask
	selected  int
	view      string
	cp        controlplane.ControlPlane
	refreshed time.Time
}

func NewFleetPanel() *FleetPanel {
	return &FleetPanel{
		width:     80,
		height:    20,
		isOpen:    false,
		agents:    []FleetAgent{},
		tasks:     []FleetTask{},
		selected:  0,
		view:      "agents",
		cp:        nil,
		refreshed: time.Now(),
	}
}

func (f *FleetPanel) SetWidth(w int) {
	f.width = w
}

func (f *FleetPanel) SetHeight(h int) {
	f.height = h
}

func (f *FleetPanel) Toggle() {
	f.isOpen = !f.isOpen
	if f.isOpen {
		f.refresh()
	}
}

func (f *FleetPanel) Open() {
	f.isOpen = true
	f.refresh()
}

func (f *FleetPanel) Close() {
	f.isOpen = false
}

func (f *FleetPanel) IsOpen() bool {
	return f.isOpen
}

func (f *FleetPanel) SetControlPlane(cp controlplane.ControlPlane) {
	f.cp = cp
}

func (f *FleetPanel) refresh() {
	f.refreshed = time.Now()
	if f.cp == nil {
		return
	}

	ctx := context.Background()
	agentList, err := f.cp.ListAgents(ctx)
	if err == nil {
		f.agents = make([]FleetAgent, 0, len(agentList))
		for _, a := range agentList {
			f.agents = append(f.agents, FleetAgent{
				ID:        a.ID,
				Name:      a.Name,
				Status:    a.Status,
				Load:      a.Load,
				Connected: a.LastSeen,
			})
		}
	}

	taskList, err := f.cp.ListTasks(ctx, "")
	if err == nil {
		f.tasks = make([]FleetTask, 0, len(taskList))
		for _, t := range taskList {
			f.tasks = append(f.tasks, FleetTask{
				ID:        t.ID,
				Type:      t.Type,
				Status:    t.Status,
				AgentID:   t.AgentID,
				Created:   t.Created,
				Completed: t.Completed,
			})
		}
	}
}

func (f *FleetPanel) HandleKey(msg string) bool {
	if !f.isOpen {
		return false
	}

	switch msg {
	case "escape":
		f.Close()
		return true
	case "r":
		f.refresh()
		return true
	case "tab":
		if f.view == "agents" {
			f.view = "tasks"
		} else {
			f.view = "agents"
		}
		f.selected = 0
		return true
	case "up", "k":
		f.moveUp()
		return true
	case "down", "j":
		f.moveDown()
		return true
	case "enter":
		f.handleSelect()
		return true
	}
	return false
}

func (f *FleetPanel) moveUp() {
	maxItems := f.maxItems()
	if f.selected > 0 {
		f.selected--
	}
	_ = maxItems
}

func (f *FleetPanel) moveDown() {
	maxItems := f.maxItems()
	if f.selected < maxItems-1 {
		f.selected++
	}
}

func (f *FleetPanel) maxItems() int {
	if f.view == "agents" {
		return len(f.agents)
	}
	return len(f.tasks)
}

func (f *FleetPanel) handleSelect() {
}

func (f *FleetPanel) Render() string {
	if !f.isOpen {
		return ""
	}

	dialogStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#0D1117")).
		Border(lipgloss.HiddenBorder()).
		Width(f.width).
		Height(f.height)

	return dialogStyle.Render(f.renderContent())
}

func (f *FleetPanel) renderContent() string {
	var lines []string

	lines = append(lines, f.renderHeader())
	lines = append(lines, "")

	if f.view == "agents" {
		lines = append(lines, f.renderAgents()...)
	} else {
		lines = append(lines, f.renderTasks()...)
	}

	lines = append(lines, "")
	lines = append(lines, f.renderHints())

	return strings.Join(lines, "\n")
}

func (f *FleetPanel) renderHeader() string {
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E0E0E0")).
		Background(lipgloss.Color("#161B22")).
		Padding(0, 1).
		Bold(true)

	agentsTab := "Agents"
	tasksTab := "Tasks"

	if f.view == "agents" {
		agentsTab = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#58A6FF")).
			Background(lipgloss.Color("#21262D")).
			Padding(0, 1).
			Render("● Agents")
		tasksTab = "  Tasks  "
	} else {
		tasksTab = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#58A6FF")).
			Background(lipgloss.Color("#21262D")).
			Padding(0, 1).
			Render("● Tasks")
		agentsTab = "  Agents  "
	}

	refreshStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
	refreshStr := refreshStyle.Render(fmt.Sprintf("refreshed %s", f.refreshed.Format("15:04:05")))

	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
	hintStr := hintStyle.Render("esc close  r refresh  tab switch view")

	return headerStyle.Render("Fleet") + agentsTab + tasksTab + "  " + refreshStr + "  " + hintStr
}

func (f *FleetPanel) renderAgents() []string {
	var lines []string

	if len(f.agents) == 0 {
		emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
		lines = append(lines, emptyStyle.Render("  No agents connected"))
		lines = append(lines, "")
		lines = append(lines, emptyStyle.Render("  Start a fleet head to connect agents"))
		return lines
	}

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8B949E")).
		Bold(true)

	lines = append(lines, headerStyle.Render("  ID                    NAME                STATUS    LOAD   LAST SEEN"))
	lines = append(lines, headerStyle.Render("  "+strings.Repeat("─", 70)))

	for i, agent := range f.agents {
		lines = append(lines, f.renderAgent(agent, i))
	}

	return lines
}

func (f *FleetPanel) renderAgent(agent FleetAgent, idx int) string {
	selected := idx == f.selected
	prefix := "  "
	if selected {
		prefix = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFCC00")).
			Render("▶")
	}

	idStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#8B949E"))
	idStr := idStyle.Render(truncate(agent.ID, 20))

	nameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#E0E0E0"))
	nameStr := nameStyle.Render(truncate(agent.Name, 16))

	var statusStyle lipgloss.Style
	statusStr := agent.Status
	switch agent.Status {
	case "online":
		statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#4EC9B0"))
	case "offline":
		statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#F44747"))
	case "busy":
		statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#DCDCAA"))
	default:
		statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
	}
	statusStr = statusStyle.Render(statusStr)

	loadStr := fmt.Sprintf("%.0f%%", agent.Load*100)
	loadStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#8B949E"))
	if agent.Load > 0.8 {
		loadStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#F44747"))
	} else if agent.Load > 0.5 {
		loadStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#DCDCAA"))
	} else {
		loadStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#4EC9B0"))
	}
	loadStr = loadStyle.Render(loadStr)

	connectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#8B949E"))
	connectedStr := connectedStyle.Render(agent.Connected.Format("15:04:05"))

	return prefix + " " + idStr + "  " + nameStr + "  " + statusStr + "  " + loadStr + "  " + connectedStr
}

func (f *FleetPanel) renderTasks() []string {
	var lines []string

	if len(f.tasks) == 0 {
		emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
		lines = append(lines, emptyStyle.Render("  No tasks in queue"))
		return lines
	}

	sort.Slice(f.tasks, func(i, j int) bool {
		return f.tasks[i].Created.After(f.tasks[j].Created)
	})

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8B949E")).
		Bold(true)

	lines = append(lines, headerStyle.Render("  ID                    TYPE          STATUS      AGENT               CREATED"))
	lines = append(lines, headerStyle.Render("  "+strings.Repeat("─", 70)))

	for i, task := range f.tasks {
		lines = append(lines, f.renderTask(task, i))
	}

	return lines
}

func (f *FleetPanel) renderTask(task FleetTask, idx int) string {
	selected := idx == f.selected
	prefix := "  "
	if selected {
		prefix = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFCC00")).
			Render("▶")
	}

	idStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#8B949E"))
	idStr := idStyle.Render(truncate(task.ID, 20))

	typeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#E0E0E0"))
	typeStr := typeStyle.Render(truncate(task.Type, 12))

	var statusStyle lipgloss.Style
	statusStr := task.Status
	switch task.Status {
	case "pending":
		statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#DCDCAA"))
	case "assigned":
		statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#58A6FF"))
	case "completed":
		statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#4EC9B0"))
	case "failed":
		statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#F44747"))
	default:
		statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
	}
	statusStr = statusStyle.Render(statusStr)

	agentStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#8B949E"))
	agentStr := agentStyle.Render(truncate(task.AgentID, 18))

	createdStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#8B949E"))
	createdStr := createdStyle.Render(task.Created.Format("15:04:05"))

	return prefix + " " + idStr + "  " + typeStr + "  " + statusStr + "  " + agentStr + "  " + createdStr
}

func (f *FleetPanel) renderHints() string {
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
	if f.view == "agents" {
		return hintStyle.Render("↑↓ select  enter details  r refresh  esc close")
	}
	return hintStyle.Render("↑↓ select  enter details  r refresh  tab agents  esc close")
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
