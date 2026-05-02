package hook

type ToolHooks struct {
	registry *Registry
}

func NewToolHooks(r *Registry) *ToolHooks {
	return &ToolHooks{registry: r}
}

func (h *ToolHooks) OnBefore(name string, fn ToolHook) {
	h.registry.RegisterToolHook("before:"+name, fn)
}

func (h *ToolHooks) OnAfter(name string, fn ToolHook) {
	h.registry.RegisterToolHook("after:"+name, fn)
}

func (h *ToolHooks) OnError(name string, fn ToolHook) {
	h.registry.RegisterToolHook("error:"+name, fn)
}

func (h *ToolHooks) OnTransform(fn TransformHook) {
	h.registry.RegisterTransformHook(fn)
}

func (h *ToolHooks) OnBash(fn ToolHook) {
	h.registry.RegisterToolHook("before:bash", fn)
}

func (h *ToolHooks) OnRead(fn ToolHook) {
	h.registry.RegisterToolHook("before:read", fn)
}

func (h *ToolHooks) OnWrite(fn ToolHook) {
	h.registry.RegisterToolHook("before:write", fn)
}

func (h *ToolHooks) OnEdit(fn ToolHook) {
	h.registry.RegisterToolHook("before:edit", fn)
}

func (h *ToolHooks) OnGlob(fn ToolHook) {
	h.registry.RegisterToolHook("before:glob", fn)
}

func (h *ToolHooks) OnGrep(fn ToolHook) {
	h.registry.RegisterToolHook("before:grep", fn)
}

func (h *ToolHooks) OnWebFetch(fn ToolHook) {
	h.registry.RegisterToolHook("before:webfetch", fn)
}

func (h *ToolHooks) OnWebSearch(fn ToolHook) {
	h.registry.RegisterToolHook("before:websearch", fn)
}

func (h *ToolHooks) OnTask(fn ToolHook) {
	h.registry.RegisterToolHook("before:task", fn)
}

func (h *ToolHooks) OnSkill(fn ToolHook) {
	h.registry.RegisterToolHook("before:skill", fn)
}
