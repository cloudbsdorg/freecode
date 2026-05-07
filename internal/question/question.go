package question

import (
	"context"
)

type Question struct {
	ID       string
	Text     string
	Options  []Option
	Required bool
}

type Option struct {
	ID    string
	Label string
	Value string
}

type Answer struct {
	QuestionID string
	OptionID   string
	Value      string
}

type Flow interface {
	Ask(ctx context.Context, q *Question) (*Answer, error)
	Validate(ctx context.Context, answer *Answer) error
}

type memoryFlow struct{}

func NewMemoryFlow() Flow {
	return &memoryFlow{}
}

func (f *memoryFlow) Ask(ctx context.Context, q *Question) (*Answer, error) {
	if len(q.Options) > 0 {
		return &Answer{
			QuestionID: q.ID,
			OptionID:   q.Options[0].ID,
			Value:      q.Options[0].Value,
		}, nil
	}
	return &Answer{QuestionID: q.ID}, nil
}

func (f *memoryFlow) Validate(ctx context.Context, answer *Answer) error {
	if answer == nil {
		return nil
	}
	return nil
}
