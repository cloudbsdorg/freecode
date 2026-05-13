package ui

import (
	"strings"

	"github.com/freecode/freecode/internal/style"
)

type QuestionOption struct {
	Label       string
	Description string
}

type QuestionItem struct {
	Question   string
	Header     string
	Options    []QuestionOption
	Multiple   bool
	Custom     bool
}

type QuestionRequest struct {
	ID        string
	SessionID string
	Questions []QuestionItem
}

type QuestionAnswer struct {
	QuestionIndex int
	Answers      []string
}

type QuestionDialogState struct {
	Request  *QuestionRequest
	Tab      int
	Answers  [][]string
	Custom   []string
	Selected int
	Editing  bool
	Width    int
}

type QuestionDialog struct {
	state QuestionDialogState
}

func NewQuestionDialog() *QuestionDialog {
	return &QuestionDialog{
		state: QuestionDialogState{
			Tab:      0,
			Answers:  [][]string{},
			Custom:   []string{},
			Selected: 0,
			Editing:  false,
			Width:    60,
		},
	}
}

func (q *QuestionDialog) SetWidth(w int) {
	q.state.Width = w
}

func (q *QuestionDialog) SetRequest(req *QuestionRequest) {
	q.state.Request = req
	q.state.Tab = 0
	q.state.Selected = 0
	q.state.Editing = false

	numQuestions := len(req.Questions)
	q.state.Answers = make([][]string, numQuestions)
	q.state.Custom = make([]string, numQuestions)
	for i := 0; i < numQuestions; i++ {
		q.state.Answers[i] = []string{}
	}
}

func (q *QuestionDialog) GetRequest() *QuestionRequest {
	return q.state.Request
}

func (q *QuestionDialog) IsVisible() bool {
	return q.state.Request != nil
}

func (q *QuestionDialog) Clear() {
	q.state.Request = nil
	q.state.Tab = 0
	q.state.Answers = nil
	q.state.Custom = nil
	q.state.Selected = 0
	q.state.Editing = false
}

func (q *QuestionDialog) IsSingle() bool {
	if q.state.Request == nil {
		return true
	}
	return len(q.state.Request.Questions) == 1 && !q.state.Request.Questions[0].Multiple
}

func (q *QuestionDialog) IsConfirm() bool {
	return !q.IsSingle() && q.state.Tab == len(q.state.Request.Questions)
}

func (q *QuestionDialog) CurrentQuestion() *QuestionItem {
	if q.state.Request == nil || q.state.Tab >= len(q.state.Request.Questions) {
		return nil
	}
	return &q.state.Request.Questions[q.state.Tab]
}

func (q *QuestionDialog) TabCount() int {
	if q.state.Request == nil {
		return 0
	}
	if q.IsSingle() {
		return 1
	}
	return len(q.state.Request.Questions) + 1
}

func (q *QuestionDialog) NextTab() {
	tabs := q.TabCount()
	q.state.Tab = (q.state.Tab + 1) % tabs
	q.state.Selected = 0
}

func (q *QuestionDialog) PrevTab() {
	tabs := q.TabCount()
	q.state.Tab = (q.state.Tab - 1 + tabs) % tabs
	q.state.Selected = 0
}

func (q *QuestionDialog) NextOption() {
	question := q.CurrentQuestion()
	if question == nil {
		return
	}
	total := len(question.Options)
	if question.Custom {
		total++
	}
	q.state.Selected = (q.state.Selected + 1) % total
}

func (q *QuestionDialog) PrevOption() {
	question := q.CurrentQuestion()
	if question == nil {
		return
	}
	total := len(question.Options)
	if question.Custom {
		total++
	}
	q.state.Selected = (q.state.Selected - 1 + total) % total
}

func (q *QuestionDialog) IsOtherSelected() bool {
	question := q.CurrentQuestion()
	if question == nil {
		return false
	}
	return question.Custom && q.state.Selected == len(question.Options)
}

func (q *QuestionDialog) ToggleCurrentOption() {
	question := q.CurrentQuestion()
	if question == nil {
		return
	}

	if question.Multiple {
		q.toggleMulti()
	} else {
		q.pickSingle()
	}
}

func (q *QuestionDialog) toggleMulti() {
	question := q.CurrentQuestion()
	if question == nil {
		return
	}

	options := question.Options
	custom := question.Custom

	if q.state.Selected < len(options) {
		label := options[q.state.Selected].Label
		q.toggleAnswer(label)
	} else if custom && q.state.Selected == len(options) {
		if q.state.Custom[q.state.Tab] != "" {
			q.toggleAnswer(q.state.Custom[q.state.Tab])
		} else {
			q.state.Editing = true
		}
	}
}

func (q *QuestionDialog) toggleAnswer(answer string) {
	answers := q.state.Answers[q.state.Tab]
	found := false
	for i, a := range answers {
		if a == answer {
			q.state.Answers[q.state.Tab] = append(answers[:i], answers[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		q.state.Answers[q.state.Tab] = append(answers, answer)
	}
}

func (q *QuestionDialog) pickSingle() {
	question := q.CurrentQuestion()
	if question == nil {
		return
	}

	options := question.Options
	custom := question.Custom

	if q.state.Selected < len(options) {
		label := options[q.state.Selected].Label
		q.state.Answers[q.state.Tab] = []string{label}
		if q.IsSingle() {
			return
		}
		q.NextTab()
	} else if custom && q.state.Selected == len(options) {
		if q.state.Custom[q.state.Tab] != "" {
			q.state.Answers[q.state.Tab] = []string{q.state.Custom[q.state.Tab]}
			if q.IsSingle() {
				return
			}
			q.NextTab()
		} else {
			q.state.Editing = true
		}
	}
}

func (q *QuestionDialog) SetCustomInput(text string) {
	q.state.Custom[q.state.Tab] = text
}

func (q *QuestionDialog) GetCustomInput() string {
	return q.state.Custom[q.state.Tab]
}

func (q *QuestionDialog) SetEditing(editing bool) {
	q.state.Editing = editing
}

func (q *QuestionDialog) IsEditing() bool {
	return q.state.Editing
}

func (q *QuestionDialog) Submit() [][]string {
	return q.state.Answers
}

func (q *QuestionDialog) Render() string {
	if !q.IsVisible() {
		return ""
	}

	req := q.state.Request
	if req == nil {
		return ""
	}

	dialogStyle := style.NewStyle().
		Background(style.Color("#1E1E1E")).
		BorderStyle(style.HiddenBorder()).
		Width(q.state.Width)

	return dialogStyle.Render(q.renderContent())
}

func (q *QuestionDialog) renderContent() string {
	var lines []string

	if !q.IsSingle() {
		lines = append(lines, q.renderTabs())
		lines = append(lines, "")
	}

	if q.IsConfirm() {
		lines = append(lines, q.renderConfirmStage()...)
	} else {
		lines = append(lines, q.renderQuestionStage()...)
	}

	lines = append(lines, "")
	lines = append(lines, q.renderHints())

	return strings.Join(lines, "\n")
}

func (q *QuestionDialog) renderTabs() string {
	req := q.state.Request
	var tabs []string

	for i, question := range req.Questions {
		active := i == q.state.Tab
		answered := len(q.state.Answers[i]) > 0

		var tabStr string
		if active {
			tabStr = style.NewStyle().
				Background(style.Color("#007ACC")).
				Foreground(style.Color("#FFFFFF")).
				Render(" " + question.Header + " ")
		} else if answered {
			tabStr = style.NewStyle().
				Foreground(style.Color("#E0E0E0")).
				Render(" " + question.Header + " ")
		} else {
			tabStr = style.NewStyle().
				Foreground(style.Color("#808080")).
				Render(" " + question.Header + " ")
		}
		tabs = append(tabs, tabStr)
	}

	confirmActive := q.state.Tab == len(req.Questions)
	var confirmTab string
	if confirmActive {
		confirmTab = style.NewStyle().
			Background(style.Color("#007ACC")).
			Foreground(style.Color("#FFFFFF")).
			Render(" Confirm ")
	} else {
		confirmTab = style.NewStyle().
			Foreground(style.Color("#808080")).
			Render(" Confirm ")
	}
	tabs = append(tabs, confirmTab)

	return strings.Join(tabs, " ")
}

func (q *QuestionDialog) renderQuestionStage() []string {
	var lines []string
	question := q.CurrentQuestion()
	if question == nil {
		return lines
	}

	multiText := ""
	if question.Multiple {
		multiText = " (select all that apply)"
	}

	lines = append(lines, style.NewStyle().
		Foreground(style.Color("#E0E0E0")).
		Render(question.Question+multiText))

	lines = append(lines, "")

	for i, opt := range question.Options {
		selected := i == q.state.Selected
		picked := q.isAnswerPicked(opt.Label)

		var label string
		if question.Multiple {
			check := " "
			if picked {
				check = "✓"
			}
			label = "[" + check + "] " + opt.Label
		} else {
			if picked {
				label = opt.Label + " ✓"
			} else {
				label = opt.Label
			}
		}

		numStr := style.NewStyle().
			Foreground(style.Color("#808080")).
			Render(string(rune('1' + i)) + ".")

		var labelStyle style.Style
		if selected {
			if picked {
				labelStyle = style.NewStyle().Foreground(style.Color("#4EC9B0"))
			} else {
				labelStyle = style.NewStyle().Foreground(style.Color("#DCDCAA"))
			}
		} else {
			labelStyle = style.NewStyle().Foreground(style.Color("#E0E0E0"))
		}

		lines = append(lines, "  "+numStr+" "+labelStyle.Render(label))

		if opt.Description != "" {
			descStyle := style.NewStyle().Foreground(style.Color("#808080"))
			lines = append(lines, "     "+descStyle.Render(opt.Description))
		}
	}

	if question.Custom {
		i := len(question.Options)
		selected := i == q.state.Selected
		customPicked := q.isAnswerPicked(q.state.Custom[q.state.Tab])

		var label string
		if question.Multiple {
			check := " "
			if customPicked {
				check = "✓"
			}
			label = "[" + check + "] Type your own answer"
		} else {
			if customPicked {
				label = "Type your own answer ✓"
			} else {
				label = "Type your own answer"
			}
		}

		numStr := style.NewStyle().
			Foreground(style.Color("#808080")).
			Render(string(rune('1' + i)) + ".")

		var labelStyle style.Style
		if selected {
			labelStyle = style.NewStyle().Foreground(style.Color("#DCDCAA"))
		} else {
			labelStyle = style.NewStyle().Foreground(style.Color("#E0E0E0"))
		}

		lines = append(lines, "  "+numStr+" "+labelStyle.Render(label))

		customInput := q.state.Custom[q.state.Tab]
		if q.state.Editing && selected {
			inputStyle := style.NewStyle().
				Background(style.Color("#3C3C3C")).
				Foreground(style.Color("#E0E0E0")).
				Padding(0, 1)
			lines = append(lines, "     "+inputStyle.Render(customInput+"_"))
		} else if customPicked && customInput != "" {
			descStyle := style.NewStyle().Foreground(style.Color("#808080"))
			lines = append(lines, "     "+descStyle.Render(customInput))
		}
	}

	return lines
}

func (q *QuestionDialog) renderConfirmStage() []string {
	var lines []string

	lines = append(lines, style.NewStyle().
		Foreground(style.Color("#E0E0E0")).
		Render("Review your answers"))

	lines = append(lines, "")

	for i, question := range q.state.Request.Questions {
		answer := ""
		if len(q.state.Answers[i]) > 0 {
			answer = strings.Join(q.state.Answers[i], ", ")
		} else {
			answer = "(not answered)"
		}

		headerStyle := style.NewStyle().Foreground(style.Color("#808080"))
		answerStyle := style.NewStyle().Foreground(style.Color("#E0E0E0"))

		if len(q.state.Answers[i]) == 0 {
			answerStyle = style.NewStyle().Foreground(style.Color("#F44747"))
		}

		lines = append(lines, "  "+headerStyle.Render(question.Header+":")+" "+answerStyle.Render(answer))
	}

	return lines
}

func (q *QuestionDialog) renderHints() string {
	hints := []string{}

	if !q.IsSingle() {
		hints = append(hints, "⇆ tab")
	}

	if !q.IsConfirm() {
		hints = append(hints, "↑↓ select")
	}

	if q.IsConfirm() {
		hints = append(hints, "enter submit")
	} else if q.CurrentQuestion() != nil && q.CurrentQuestion().Multiple {
		hints = append(hints, "enter toggle")
	} else if q.IsSingle() {
		hints = append(hints, "enter submit")
	} else {
		hints = append(hints, "enter confirm")
	}

	hints = append(hints, "esc dismiss")

	hintStyle := style.NewStyle().Foreground(style.Color("#808080"))
	result := ""
	for i, h := range hints {
		if i > 0 {
			result += "  "
		}
		result += hintStyle.Render(h)
	}
	return result
}

func (q *QuestionDialog) isAnswerPicked(answer string) bool {
	for _, a := range q.state.Answers[q.state.Tab] {
		if a == answer {
			return true
		}
	}
	return false
}