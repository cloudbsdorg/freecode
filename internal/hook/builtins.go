package hook

import (
	"context"
	"log"
)

func RegisterBuiltinHooks(r *Registry) {
	r.RegisterSessionHook("session.start", func(ctx context.Context, evt SessionEvent) error {
		log.Printf("[hook] session started: %s", evt.SessionID)
		return nil
	})

	r.RegisterSessionHook("session.end", func(ctx context.Context, evt SessionEvent) error {
		log.Printf("[hook] session ended: %s", evt.SessionID)
		return nil
	})

	r.RegisterSessionHook("session.error", func(ctx context.Context, evt SessionEvent) error {
		if err, ok := evt.Data["error"].(string); ok {
			log.Printf("[hook] session error: %s - %s", evt.SessionID, err)
		}
		return nil
	})

	r.RegisterToolHook("tool.execute.before", func(ctx context.Context, evt ToolEvent) (error, bool) {
		return nil, false
	})

	r.RegisterToolHook("tool.execute.after", func(ctx context.Context, evt ToolEvent) (error, bool) {
		return nil, false
	})

	r.RegisterToolHook("tool.execute.error", func(ctx context.Context, evt ToolEvent) (error, bool) {
		return nil, false
	})

	r.RegisterTransformHook(func(msg *Message) (*Message, error) {
		return msg, nil
	})

	r.RegisterContinuationHook(func(ctx context.Context, session *SessionData) (*ContinueSignal, error) {
		return &ContinueSignal{Continue: true}, nil
	})

	r.RegisterRalphHook(func(ctx context.Context, input string) (string, error) {
		return input, nil
	})
}
