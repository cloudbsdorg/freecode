package hook

import (
	"context"
	"fmt"
)

type Trigger struct {
	registry *Registry
}

func NewTrigger(r *Registry) *Trigger {
	return &Trigger{registry: r}
}

func (t *Trigger) SessionStart(ctx context.Context, sessionID string) error {
	return t.registry.EmitSessionEvent(ctx, "session.start", sessionID, nil)
}

func (t *Trigger) SessionEnd(ctx context.Context, sessionID string) error {
	return t.registry.EmitSessionEvent(ctx, "session.end", sessionID, nil)
}

func (t *Trigger) SessionError(ctx context.Context, sessionID string, err error) error {
	return t.registry.EmitSessionEvent(ctx, "session.error", sessionID, map[string]interface{}{
		"error": err.Error(),
	})
}

func (t *Trigger) SessionCreated(ctx context.Context, sessionID string) error {
	return t.registry.EmitSessionEvent(ctx, "session.created", sessionID, nil)
}

func (t *Trigger) SessionDeleted(ctx context.Context, sessionID string) error {
	return t.registry.EmitSessionEvent(ctx, "session.deleted", sessionID, nil)
}

func (t *Trigger) SessionTitleChanged(ctx context.Context, sessionID, title string) error {
	return t.registry.EmitSessionEvent(ctx, "session.title_changed", sessionID, map[string]interface{}{
		"title": title,
	})
}

func (t *Trigger) SessionRenamed(ctx context.Context, sessionID, oldName, newName string) error {
	return t.registry.EmitSessionEvent(ctx, "session.renamed", sessionID, map[string]interface{}{
		"old_name": oldName,
		"new_name": newName,
	})
}

func (t *Trigger) SessionForked(ctx context.Context, sessionID, newSessionID string) error {
	return t.registry.EmitSessionEvent(ctx, "session.forked", sessionID, map[string]interface{}{
		"new_session_id": newSessionID,
	})
}

func (t *Trigger) SessionMerged(ctx context.Context, sessionID, mergedSessionID string) error {
	return t.registry.EmitSessionEvent(ctx, "session.merged", sessionID, map[string]interface{}{
		"merged_session_id": mergedSessionID,
	})
}

func (t *Trigger) SessionShared(ctx context.Context, sessionID, shareURL string) error {
	return t.registry.EmitSessionEvent(ctx, "session.shared", sessionID, map[string]interface{}{
		"share_url": shareURL,
	})
}

func (t *Trigger) SessionExported(ctx context.Context, sessionID, exportPath string) error {
	return t.registry.EmitSessionEvent(ctx, "session.exported", sessionID, map[string]interface{}{
		"export_path": exportPath,
	})
}

func (t *Trigger) SessionImported(ctx context.Context, sessionID, importPath string) error {
	return t.registry.EmitSessionEvent(ctx, "session.imported", sessionID, map[string]interface{}{
		"import_path": importPath,
	})
}

func (t *Trigger) SessionNotification(ctx context.Context, sessionID, message string) error {
	return t.registry.EmitSessionEvent(ctx, "session.notification", sessionID, map[string]interface{}{
		"message": message,
	})
}

func (t *Trigger) SessionErrorRecovery(ctx context.Context, sessionID string, attempt int) error {
	return t.registry.EmitSessionEvent(ctx, "session.error_recovery", sessionID, map[string]interface{}{
		"attempt": attempt,
	})
}

func (t *Trigger) SessionContextExhausted(ctx context.Context, sessionID string, percentUsed float64) error {
	return t.registry.EmitSessionEvent(ctx, "session.context_exhausted", sessionID, map[string]interface{}{
		"percent_used": percentUsed,
	})
}

func (t *Trigger) SessionCompactionStart(ctx context.Context, sessionID string) error {
	return t.registry.EmitSessionEvent(ctx, "session.compaction_start", sessionID, nil)
}

func (t *Trigger) SessionCompactionEnd(ctx context.Context, sessionID string, tokensBefore, tokensAfter int) error {
	return t.registry.EmitSessionEvent(ctx, "session.compaction_end", sessionID, map[string]interface{}{
		"tokens_before": tokensBefore,
		"tokens_after":  tokensAfter,
	})
}

func (t *Trigger) SessionTabCreated(ctx context.Context, sessionID, tabID string) error {
	return t.registry.EmitSessionEvent(ctx, "session.tab_created", sessionID, map[string]interface{}{
		"tab_id": tabID,
	})
}

func (t *Trigger) SessionTabClosed(ctx context.Context, sessionID, tabID string) error {
	return t.registry.EmitSessionEvent(ctx, "session.tab_closed", sessionID, map[string]interface{}{
		"tab_id": tabID,
	})
}

func (t *Trigger) SessionTabChanged(ctx context.Context, sessionID, tabID string) error {
	return t.registry.EmitSessionEvent(ctx, "session.tab_changed", sessionID, map[string]interface{}{
		"tab_id": tabID,
	})
}

func (t *Trigger) SessionMessageAdded(ctx context.Context, sessionID, role, content string) error {
	return t.registry.EmitSessionEvent(ctx, "session.message_added", sessionID, map[string]interface{}{
		"role":    role,
		"content": content,
	})
}

func (t *Trigger) SessionTurnEnd(ctx context.Context, sessionID string, turnNumber int) error {
	return t.registry.EmitSessionEvent(ctx, "session.turn_end", sessionID, map[string]interface{}{
		"turn_number": turnNumber,
	})
}

func (t *Trigger) SessionUltraworkStart(ctx context.Context, sessionID string) error {
	return t.registry.EmitSessionEvent(ctx, "session.ultrawork_start", sessionID, nil)
}

func (t *Trigger) SessionUltraworkEnd(ctx context.Context, sessionID string) error {
	return t.registry.EmitSessionEvent(ctx, "session.ultrawork_end", sessionID, nil)
}

func (t *Trigger) ToolExecuteBefore(ctx context.Context, toolName, sessionID string, args map[string]interface{}) error {
	return t.emitToolEvent(ctx, "tool.execute.before", toolName, sessionID, args)
}

func (t *Trigger) ToolExecuteAfter(ctx context.Context, toolName, sessionID string, args map[string]interface{}, result interface{}) error {
	return t.emitToolEvent(ctx, "tool.execute.after", toolName, sessionID, args, result)
}

func (t *Trigger) ToolExecuteError(ctx context.Context, toolName, sessionID string, args map[string]interface{}, err error) error {
	return t.emitToolEvent(ctx, "tool.execute.error", toolName, sessionID, args, nil, err)
}

func (t *Trigger) ToolExecuteTimeout(ctx context.Context, toolName, sessionID string, args map[string]interface{}) error {
	return t.emitToolEvent(ctx, "tool.execute.timeout", toolName, sessionID, args)
}

func (t *Trigger) ToolConfirm(ctx context.Context, toolName, sessionID string, args map[string]interface{}) error {
	return t.emitToolEvent(ctx, "tool.confirm", toolName, sessionID, args)
}

func (t *Trigger) ToolConfirmDeny(ctx context.Context, toolName, sessionID string, args map[string]interface{}) error {
	return t.emitToolEvent(ctx, "tool.confirm.deny", toolName, sessionID, args)
}

func (t *Trigger) ToolConfirmAllow(ctx context.Context, toolName, sessionID string, args map[string]interface{}) error {
	return t.emitToolEvent(ctx, "tool.confirm.allow", toolName, sessionID, args)
}

func (t *Trigger) ToolRateLimit(ctx context.Context, toolName, sessionID string, retryAfter int) error {
	return t.registry.EmitSessionEvent(ctx, "tool.rate_limit", sessionID, map[string]interface{}{
		"tool_name":    toolName,
		"retry_after":  retryAfter,
	})
}

func (t *Trigger) emitToolEvent(ctx context.Context, eventType, toolName, sessionID string, args map[string]interface{}, result ...interface{}) error {
	r := t.registry
	r.mu.RLock()
	hooks := r.toolHooks[eventType]
	r.mu.RUnlock()

	if len(hooks) == 0 {
		return nil
	}

	evt := ToolEvent{
		Type:      eventType,
		ToolName:  toolName,
		SessionID: sessionID,
		Arguments: args,
	}

	if len(result) > 0 {
		evt.Result = result[0]
	}
	if len(result) > 1 {
		if err, ok := result[1].(error); ok {
			evt.Error = err
		}
	}

	for _, hook := range hooks {
		if err, handled := hook(ctx, evt); handled {
			if err != nil {
				return fmt.Errorf("tool hook error (%s): %w", eventType, err)
			}
			return nil
		}
	}
	return nil
}
