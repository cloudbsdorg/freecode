package agent

import (
	"context"
)

type MessageHandler struct{}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{}
}

func (h *MessageHandler) Handle(ctx context.Context, msg Message, eng *Engine) (*Response, error) {
	req := Request{
		AgentName: "sisyphus",
		Message:   msg,
	}
	return eng.Run(ctx, req)
}

func (h *MessageHandler) FormatResponse(resp *Response) string {
	if resp.Error != nil {
		return "Error: " + resp.Error.Error()
	}
	return resp.Message.Content
}
