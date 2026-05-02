package hook

type ContinuationHooks struct {
	registry *Registry
}

func NewContinuationHooks(r *Registry) *ContinuationHooks {
	return &ContinuationHooks{registry: r}
}

func (h *ContinuationHooks) OnIterate(fn ContinuationHook) {
	h.registry.RegisterContinuationHook(fn)
}

func (h *ContinuationHooks) OnRetry(fn ContinuationHook) {
	h.registry.RegisterContinuationHook(fn)
}

func (h *ContinuationHooks) OnFallback(fn ContinuationHook) {
	h.registry.RegisterContinuationHook(fn)
}

func (h *ContinuationHooks) OnEscalate(fn ContinuationHook) {
	h.registry.RegisterContinuationHook(fn)
}

func (h *ContinuationHooks) OnDelegate(fn ContinuationHook) {
	h.registry.RegisterContinuationHook(fn)
}

func (h *ContinuationHooks) OnTimeout(fn ContinuationHook) {
	h.registry.RegisterContinuationHook(fn)
}

func (h *ContinuationHooks) OnIdle(fn ContinuationHook) {
	h.registry.RegisterContinuationHook(fn)
}
