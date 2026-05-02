package agent

import (
	"context"
	"fmt"

	"github.com/freecode/freecode/internal/tool"
)

type ToolCaller struct{}

func NewToolCaller() *ToolCaller {
	return &ToolCaller{}
}

func (tc *ToolCaller) Call(ctx context.Context, name string, args map[string]interface{}, eng *Engine) (*tool.Response, error) {
	t, ok := eng.GetTool(name)
	if !ok {
		return nil, fmt.Errorf("tool not found: %s", name)
	}

	req := tool.Request{
		Name:      name,
		Arguments: args,
	}

	return t.Execute(ctx, req)
}

func (tc *ToolCaller) CallBatch(ctx context.Context, calls []ToolCall, eng *Engine) ([]*tool.Response, error) {
	results := make([]*tool.Response, len(calls))
	for i, call := range calls {
		resp, err := tc.Call(ctx, call.Name, call.Arguments, eng)
		if err != nil {
			results[i] = &tool.Response{Error: err}
		} else {
			results[i] = resp
		}
	}
	return results, nil
}
