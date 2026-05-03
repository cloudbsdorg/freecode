package hook

type TransformHooks struct {
	registry *Registry
}

func NewTransformHooks(r *Registry) *TransformHooks {
	return &TransformHooks{registry: r}
}

func (h *TransformHooks) OnInput(fn TransformHook) {
	h.registry.RegisterTransformHook(fn)
}

func (h *TransformHooks) OnOutput(fn TransformHook) {
	h.registry.RegisterTransformHook(fn)
}

func (h *TransformHooks) OnToolArgs(fn TransformHook) {
	h.registry.RegisterTransformHook(fn)
}

func (h *TransformHooks) OnToolResult(fn TransformHook) {
	h.registry.RegisterTransformHook(fn)
}

func (h *TransformHooks) OnError(fn TransformHook) {
	h.registry.RegisterTransformHook(fn)
}
