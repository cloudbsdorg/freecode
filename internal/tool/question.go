package tool

import (
	"context"
	"fmt"
)

type QuestionTool struct{}

func NewQuestionTool() *QuestionTool {
	return &QuestionTool{}
}

func (t *QuestionTool) Name() string {
	return "question"
}

func (t *QuestionTool) Description() string {
	return "Ask the user a question"
}

func (t *QuestionTool) Schema() ToolSchema {
	return ToolSchema{
		Name:        "question",
		Description: "Ask the user a question",
		Parameters: map[string]Parameter{
			"question": {
				Type:        "string",
				Description: "Question to ask",
				Required:    true,
			},
			"header": {
				Type:        "string",
				Description: "Question header/label",
			},
			"options": {
				Type:        "array",
				Description: "Answer options",
				Items: &Parameter{
					Type: "string",
				},
			},
			"multiple": {
				Type:        "boolean",
				Description: "Allow multiple selections",
				Default:     false,
			},
		},
	}
}

func (t *QuestionTool) Execute(ctx context.Context, req Request) (*Response, error) {
	question, ok := req.Arguments["question"].(string)
	if !ok {
		return nil, fmt.Errorf("question must be a string")
	}

	header := ""
	if h, ok := req.Arguments["header"].(string); ok {
		header = h
	}

	return &Response{
		Result: fmt.Sprintf("Question: %s (header=%s) - awaiting user response", question, header),
	}, nil
}
