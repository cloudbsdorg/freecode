package hook

import "context"

type SessionHooks struct {
	registry *Registry
}

func NewSessionHooks(r *Registry) *SessionHooks {
	return &SessionHooks{registry: r}
}

func (h *SessionHooks) OnStart(fn SessionHook) {
	h.registry.RegisterSessionHook("start", fn)
}

func (h *SessionHooks) OnEnd(fn SessionHook) {
	h.registry.RegisterSessionHook("end", fn)
}

func (h *SessionHooks) OnPause(fn SessionHook) {
	h.registry.RegisterSessionHook("pause", fn)
}

func (h *SessionHooks) OnResume(fn SessionHook) {
	h.registry.RegisterSessionHook("resume", fn)
}

func (h *SessionHooks) OnSave(fn SessionHook) {
	h.registry.RegisterSessionHook("save", fn)
}

func (h *SessionHooks) OnLoad(fn SessionHook) {
	h.registry.RegisterSessionHook("load", fn)
}

func (h *SessionHooks) OnError(fn SessionHook) {
	h.registry.RegisterSessionHook("error", fn)
}

func (h *SessionHooks) OnTimeout(fn SessionHook) {
	h.registry.RegisterSessionHook("timeout", fn)
}

func (h *SessionHooks) OnIdle(fn SessionHook) {
	h.registry.RegisterSessionHook("idle", fn)
}

func (h *SessionHooks) OnActive(fn SessionHook) {
	h.registry.RegisterSessionHook("active", fn)
}

func (h *SessionHooks) OnMessage(fn SessionHook) {
	h.registry.RegisterSessionHook("message", fn)
}

func (h *SessionHooks) OnToolCall(fn SessionHook) {
	h.registry.RegisterSessionHook("tool_call", fn)
}

func (h *SessionHooks) OnAgentCall(fn SessionHook) {
	h.registry.RegisterSessionHook("agent_call", fn)
}

func (h *SessionHooks) OnUserInput(fn SessionHook) {
	h.registry.RegisterSessionHook("user_input", fn)
}

func (h *SessionHooks) OnResponse(fn SessionHook) {
	h.registry.RegisterSessionHook("response", fn)
}

func (h *SessionHooks) OnCompaction(fn SessionHook) {
	h.registry.RegisterSessionHook("compaction", fn)
}

func (h *SessionHooks) OnTabCreate(fn SessionHook) {
	h.registry.RegisterSessionHook("tab_create", fn)
}

func (h *SessionHooks) OnTabClose(fn SessionHook) {
	h.registry.RegisterSessionHook("tab_close", fn)
}

func (h *SessionHooks) OnTabSwitch(fn SessionHook) {
	h.registry.RegisterSessionHook("tab_switch", fn)
}

func (h *SessionHooks) OnSplitCreate(fn SessionHook) {
	h.registry.RegisterSessionHook("split_create", fn)
}

func (h *SessionHooks) OnSplitClose(fn SessionHook) {
	h.registry.RegisterSessionHook("split_close", fn)
}

func (h *SessionHooks) OnFleetConnect(fn SessionHook) {
	h.registry.RegisterSessionHook("fleet_connect", fn)
}

func (h *SessionHooks) OnFleetDisconnect(fn SessionHook) {
	h.registry.RegisterSessionHook("fleet_disconnect", fn)
}

func (h *SessionHooks) OnFleetError(fn SessionHook) {
	h.registry.RegisterSessionHook("fleet_error", fn)
}

var _ context.Context = nil
