package template

type DialogEngine interface {
	Show(id string)
	Hide(id string)
	SetComponentAttr(id, key, value string)
	SetComponentContent(id string, content interface{})
	OnDialogShown(func(dialog any))
	OnDialogHidden(func(dialog any))
}

type DialogPresenter struct {
	engine     DialogEngine
	dialogID   string
	messageID  string
	titleID    string
	contentID  string
	buttonsID  string
	isPresented bool
	onDismiss  func()
	onAction   func(action string)
}

func NewDialogPresenter(engine DialogEngine, dialogID, messageID, titleID, contentID, buttonsID string) *DialogPresenter {
	return &DialogPresenter{
		engine:    engine,
		dialogID:  dialogID,
		messageID: messageID,
		titleID:   titleID,
		contentID: contentID,
		buttonsID: buttonsID,
	}
}

func (dp *DialogPresenter) SetDismissHandler(handler func()) {
	dp.onDismiss = handler
}

func (dp *DialogPresenter) SetActionHandler(handler func(action string)) {
	dp.onAction = handler
}

func (dp *DialogPresenter) Show(options DialogOptions) {
	dp.isPresented = true

	if options.Title != "" {
		dp.engine.SetComponentAttr(dp.titleID, "value", options.Title)
	}
	if options.Message != "" {
		dp.engine.SetComponentAttr(dp.messageID, "value", options.Message)
	}
	if options.Content != nil {
		dp.engine.SetComponentContent(dp.contentID, options.Content)
	}

	dp.engine.Show(dp.dialogID)
}

func (dp *DialogPresenter) Hide() {
	dp.isPresented = false
	dp.engine.Hide(dp.dialogID)
	if dp.onDismiss != nil {
		dp.onDismiss()
	}
}

func (dp *DialogPresenter) Toggle(options DialogOptions) {
	if dp.isPresented {
		dp.Hide()
	} else {
		dp.Show(options)
	}
}

func (dp *DialogPresenter) Dismiss() {
	dp.Hide()
}

func (dp *DialogPresenter) PerformAction(action string) {
	if dp.onAction != nil {
		dp.onAction(action)
	}
}

func (dp *DialogPresenter) IsPresented() bool {
	return dp.isPresented
}

func (dp *DialogPresenter) DialogID() string {
	return dp.dialogID
}

func (dp *DialogPresenter) MessageID() string {
	return dp.messageID
}

type DialogOptions struct {
	Title   string
	Message string
	Content interface{}
	Type    DialogType
	Buttons []DialogButton
}

type DialogType string

const (
	DialogTypeInfo    DialogType = "info"
	DialogTypeSuccess DialogType = "success"
	DialogTypeWarning DialogType = "warning"
	DialogTypeError   DialogType = "error"
	DialogTypeConfirm DialogType = "confirm"
)

type DialogButton struct {
	ID     string
	Label  string
	Style  string
	IsClose bool
}

type AlertPresenter struct {
	engine    DialogEngine
	dialogID  string
	messageID string
}

func NewAlertPresenter(engine DialogEngine, dialogID, messageID string) *AlertPresenter {
	return &AlertPresenter{
		engine:    engine,
		dialogID:  dialogID,
		messageID: messageID,
	}
}

func (ap *AlertPresenter) ShowMessage(title, message string, dialogType DialogType) {
	dp := NewDialogPresenter(ap.engine, ap.dialogID, ap.messageID, "", "", "")
	dp.Show(DialogOptions{
		Title:   title,
		Message: message,
		Type:    dialogType,
	})
}

func (ap *AlertPresenter) ShowError(message string) {
	ap.ShowMessage("Error", message, DialogTypeError)
}

func (ap *AlertPresenter) ShowSuccess(message string) {
	ap.ShowMessage("Success", message, DialogTypeSuccess)
}

func (ap *AlertPresenter) ShowWarning(message string) {
	ap.ShowMessage("Warning", message, DialogTypeWarning)
}

func (ap *AlertPresenter) ShowInfo(message string) {
	ap.ShowMessage("Info", message, DialogTypeInfo)
}

type ConfirmPresenter struct {
	engine    DialogEngine
	dialogID  string
	messageID string
	confirmID string
	cancelID  string
	onConfirm func()
	onCancel  func()
}

func NewConfirmPresenter(engine DialogEngine, dialogID, messageID, confirmID, cancelID string) *ConfirmPresenter {
	return &ConfirmPresenter{
		engine:    engine,
		dialogID:  dialogID,
		messageID: messageID,
		confirmID: confirmID,
		cancelID:  cancelID,
	}
}

func (cp *ConfirmPresenter) SetConfirmHandler(handler func()) {
	cp.onConfirm = handler
}

func (cp *ConfirmPresenter) SetCancelHandler(handler func()) {
	cp.onCancel = handler
}

func (cp *ConfirmPresenter) ShowConfirm(title, message string, confirmLabel, cancelLabel string) {
	dp := NewDialogPresenter(cp.engine, cp.dialogID, cp.messageID, "", "", cp.confirmID)
	dp.Show(DialogOptions{
		Title:   title,
		Message: message,
		Type:    DialogTypeConfirm,
		Buttons: []DialogButton{
			{ID: cp.cancelID, Label: cancelLabel, Style: "secondary"},
			{ID: cp.confirmID, Label: confirmLabel, Style: "primary"},
		},
	})
}

func (cp *ConfirmPresenter) Confirm() {
	if cp.onConfirm != nil {
		cp.onConfirm()
	}
}

func (cp *ConfirmPresenter) Cancel() {
	if cp.onCancel != nil {
		cp.onCancel()
	}
}

func (dp *DialogPresenter) HandleButtonClick(buttonID string) {
	switch buttonID {
	case "confirm", "ok", "yes":
		dp.PerformAction("confirm")
		dp.Dismiss()
	case "cancel", "no", "close":
		dp.Dismiss()
	default:
		dp.PerformAction(buttonID)
	}
}
