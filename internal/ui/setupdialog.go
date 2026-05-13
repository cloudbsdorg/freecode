package ui

import (
	"fmt"
	"strings"

	"github.com/freecode/freecode/internal/style"

	"github.com/freecode/freecode/internal/ui/dialog"
)

type SetupStep int

const (
	SetupStepWelcome SetupStep = iota
	SetupStepProvider
	SetupStepModel
	SetupStepAPIKey
	SetupStepDone
)

type ProviderInfo struct {
	ID    string
	Name  string
	Count int
}

type SetupDialog struct {
	step           SetupStep
	width          int
	height         int
	isOpen         bool
	providerDialog *SelectDialog
	modelDialog    *SelectDialog
	apiKeyInput    *dialog.TextInput
	providerID     string
	modelID        string
	loading        bool
	loadingMessage string
	errorMessage   string
	colors         dialog.Colors
}

func NewSetupDialog() *SetupDialog {
	s := &SetupDialog{
		step:        SetupStepWelcome,
		width:       60,
		height:      15,
		isOpen:      true,
		loading:     false,
		providerID:  "",
		modelID:     "",
		errorMessage: "",
		colors:      dialog.Dark,
	}

	s.providerDialog = &SelectDialog{
		width:  60,
		isOpen: false,
		colors: dialog.Dark,
	}
	s.providerDialog.list = dialog.NewSelectionList(
		func(d *dialog.SelectionList) {
			d.Title = "Select Provider"
			d.Width = 60
			d.Height = 12
			d.Colors = dialog.Dark
			d.SkipFilter = true
		},
	)

	s.modelDialog = &SelectDialog{
		width:  60,
		isOpen: false,
		colors: dialog.Dark,
	}
	s.modelDialog.list = dialog.NewSelectionList(
		func(d *dialog.SelectionList) {
			d.Title = "Select Model"
			d.Width = 60
			d.Height = 12
			d.Colors = dialog.Dark
			d.SkipFilter = true
		},
	)

	s.apiKeyInput = dialog.NewTextInput(
		func(t *dialog.TextInput) {
			t.Colors = s.colors
			t.Hidden = true
			t.MaxLen = 100
			t.Width = 50
		},
	)

	return s
}

func (s *SetupDialog) SetWidth(w int) {
	s.width = w
}

func (s *SetupDialog) SetProviders(providers []ProviderInfo) {
	opts := make([]SelectOption, len(providers))
	for i, p := range providers {
		opts[i] = SelectOption{
			Title: p.Name,
			Value: p.ID,
		}
	}
	s.providerDialog.SetOptions(opts)
}

func (s *SetupDialog) SetModels(models []string) {
	opts := make([]SelectOption, len(models))
	for i, m := range models {
		opts[i] = SelectOption{
			Title: m,
			Value: m,
		}
	}
	s.modelDialog.SetOptions(opts)
}

func (s *SetupDialog) IsOpen() bool {
	return s.isOpen
}

func (s *SetupDialog) Close() {
	s.isOpen = false
}

func (s *SetupDialog) GetStep() SetupStep {
	return s.step
}

func (s *SetupDialog) Next() {
	switch s.step {
	case SetupStepWelcome:
		s.step = SetupStepProvider
	case SetupStepProvider:
		item := s.providerDialog.GetSelected()
		if item != nil {
			s.providerID = item.Value
			s.step = SetupStepModel
		}
	case SetupStepModel:
		item := s.modelDialog.GetSelected()
		if item != nil {
			s.modelID = item.Value
			s.step = SetupStepAPIKey
		}
	case SetupStepAPIKey:
		s.step = SetupStepDone
	}
}

func (s *SetupDialog) Prev() {
	switch s.step {
	case SetupStepWelcome:
		s.isOpen = false
	case SetupStepProvider:
		s.step = SetupStepWelcome
	case SetupStepModel:
		s.step = SetupStepProvider
	case SetupStepAPIKey:
		s.step = SetupStepModel
	case SetupStepDone:
		s.step = SetupStepAPIKey
	}
}

func (s *SetupDialog) MoveUp() {
	s.providerDialog.Prev()
	s.modelDialog.Prev()
}

func (s *SetupDialog) MoveDown() {
	s.providerDialog.Next()
	s.modelDialog.Next()
}

func (s *SetupDialog) GetSelection() (providerID, modelID, apiKey string) {
	return s.providerID, s.modelID, s.apiKeyInput.GetValue()
}

func (s *SetupDialog) GetSelectedProviderID() string {
	item := s.providerDialog.GetSelected()
	if item == nil {
		return ""
	}
	return item.Value
}

func (s *SetupDialog) SetLoading(loading bool, message string) {
	s.loading = loading
	s.loadingMessage = message
}

func (s *SetupDialog) SetError(message string) {
	s.errorMessage = message
}

func (s *SetupDialog) AppendToAPIKey(ch rune) {
	s.apiKeyInput.Append(ch)
}

func (s *SetupDialog) BackspaceAPIKey() {
	s.apiKeyInput.Backspace()
}

func (s *SetupDialog) ClearAPIKey() {
	s.apiKeyInput.Clear()
}

func (s *SetupDialog) Render() string {
	if !s.isOpen {
		return ""
	}

	content := s.renderContent()
	return style.NewStyle().
		Background(style.Color(s.colors.Background)).
		Width(s.width).
		BorderStyle(style.RoundedBorder()).
		BorderForeground(style.Color(s.colors.Border)).
		Padding(1).
		Render(content)
}

func (s *SetupDialog) renderContent() string {
	switch s.step {
	case SetupStepWelcome:
		return s.renderWelcome()
	case SetupStepProvider:
		return s.renderProviderSelection()
	case SetupStepModel:
		return s.renderModelSelection()
	case SetupStepAPIKey:
		return s.renderAPIKeyInput()
	case SetupStepDone:
		return s.renderDone()
	}
	return ""
}

func (s *SetupDialog) renderWelcome() string {
	lines := []string{
		dialog.Header("Welcome to Freecode Setup", s.colors),
		"",
		"This wizard will help you configure Freecode.",
		"",
		"You'll need:",
		"  • An API key from your AI provider",
		"  • A model selection",
		"",
		"Press ENTER to continue or ESC to exit.",
		"",
	}

	if s.loading {
		lines = append(lines, dialog.Muted(fmt.Sprintf("%s...", s.loadingMessage), s.colors))
	}

	if s.errorMessage != "" {
		lines = append(lines, "", dialog.ErrorText(fmt.Sprintf("Error: %s", s.errorMessage), s.colors))
	}

	content := strings.Join(lines, "\n")
	return style.NewStyle().
		Width(s.width).
		Background(style.Color(s.colors.Background)).
		BorderStyle(style.RoundedBorder()).
		BorderForeground(style.Color(s.colors.Border)).
		Padding(1).
		Render(content)
}

func (s *SetupDialog) renderProviderSelection() string {
	return s.providerDialog.Render()
}

func (s *SetupDialog) renderModelSelection() string {
	return s.modelDialog.Render()
}

func (s *SetupDialog) renderAPIKeyInput() string {
	lines := []string{
		dialog.Header("Enter API Key", s.colors),
		"",
		fmt.Sprintf("Provider: %s  |  Model: %s", s.providerID, s.modelID),
		"",
		"Type your API key:",
		"",
		s.apiKeyInput.RenderWithPrefix(""),
		"",
		dialog.Muted("Press ENTER when done, ESC to go back", s.colors),
		dialog.Muted("Use Backspace to delete", s.colors),
	}

	if s.errorMessage != "" {
		lines = append(lines, "", dialog.ErrorText(fmt.Sprintf("Error: %s", s.errorMessage), s.colors))
	}

	content := strings.Join(lines, "\n")
	return style.NewStyle().
		Width(s.width).
		Background(style.Color(s.colors.Background)).
		BorderStyle(style.RoundedBorder()).
		BorderForeground(style.Color(s.colors.Border)).
		Padding(1).
		Render(content)
}

func (s *SetupDialog) renderDone() string {
	lines := []string{
		dialog.Header("Setup Complete!", s.colors),
		"",
		fmt.Sprintf("Provider: %s", s.providerID),
		fmt.Sprintf("Model:    %s", s.modelID),
		"",
		"Configuration saved successfully.",
		"",
		"Press ENTER to start using Freecode.",
	}

	content := strings.Join(lines, "\n")
	return style.NewStyle().
		Width(s.width).
		Background(style.Color(s.colors.Background)).
		BorderStyle(style.RoundedBorder()).
		BorderForeground(style.Color(s.colors.Border)).
		Padding(1).
		Render(content)
}
