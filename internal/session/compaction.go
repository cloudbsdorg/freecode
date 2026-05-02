package session

import (
	"strings"
)

type Compactor struct {
	thresholdTokens int
	keepLastN       int
}

func NewCompactor(thresholdTokens, keepLastN int) *Compactor {
	return &Compactor{
		thresholdTokens: thresholdTokens,
		keepLastN:       keepLastN,
	}
}

func (c *Compactor) Compact(sess *Session) error {
	totalTokens := c.estimateTokens(sess.Messages)
	if totalTokens < c.thresholdTokens {
		return nil
	}

	if len(sess.Messages) <= c.keepLastN {
		return nil
	}

	keepMessages := sess.Messages[len(sess.Messages)-c.keepLastN:]
	systemMessages := c.extractSystemMessages(sess.Messages)

	compacted := make([]Message, 0, len(systemMessages)+len(keepMessages))
	compacted = append(compacted, systemMessages...)
	compacted = append(compacted, keepMessages...)

	sess.Messages = compacted
	return nil
}

func (c *Compactor) extractSystemMessages(messages []Message) []Message {
	system := make([]Message, 0)
	for _, m := range messages {
		if m.Role == "system" {
			system = append(system, m)
		}
	}
	return system
}

func (c *Compactor) estimateTokens(messages []Message) int {
	total := 0
	for _, m := range messages {
		total += len(strings.Fields(m.Content)) * 4 / 3
	}
	return total
}

type History struct {
	maxSize int
}

func NewHistory(maxSize int) *History {
	return &History{maxSize: maxSize}
}

func (h *History) Add(sessions []*Session, msg Message) []*Session {
	if len(sessions) == 0 {
		return sessions
	}

	last := sessions[len(sessions)-1]
	last.Messages = append(last.Messages, msg)

	if len(last.Messages) > h.maxSize {
		last.Messages = last.Messages[h.maxSize/2:]
	}

	return sessions
}

func (h *History) Search(sessions []*Session, query string) []Message {
	results := make([]Message, 0)
	for _, sess := range sessions {
		for _, msg := range sess.Messages {
			if strings.Contains(strings.ToLower(msg.Content), strings.ToLower(query)) {
				results = append(results, msg)
			}
		}
	}
	return results
}
