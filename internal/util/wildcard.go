package util

import (
	"regexp"
	"sort"
	"strings"
)

// Match performs glob-style wildcard matching.
// * matches any characters except /
// ** matches any characters including /
// ? matches single character
// Backslashes are normalized to forward slashes.
func Match(str string, pattern string) bool {
	if str != "" {
		str = strings.ReplaceAll(str, "\\", "/")
	}
	if pattern != "" {
		pattern = strings.ReplaceAll(pattern, "\\", "/")
	}

	var result strings.Builder
	i := 0
	for i < len(pattern) {
		if i+1 < len(pattern) && pattern[i] == '*' && pattern[i+1] == '*' {
			result.WriteString("(.*)")
			i += 2
		} else if pattern[i] == '*' {
			result.WriteString("(.*)")
			i++
		} else if pattern[i] == '?' {
			result.WriteByte('.')
			i++
		} else if pattern[i] == '\\' || pattern[i] == '+' || pattern[i] == '^' || pattern[i] == '$' || pattern[i] == '|' || pattern[i] == '(' || pattern[i] == ')' || pattern[i] == '[' || pattern[i] == ']' || pattern[i] == '{' || pattern[i] == '}' {
			result.WriteByte('\\')
			result.WriteByte(pattern[i])
			i++
		} else {
			result.WriteByte(pattern[i])
			i++
		}
	}

	escaped := result.String()

	if strings.HasSuffix(escaped, " (.*)") {
		escaped = strings.TrimSuffix(escaped, " (.*)") + "( (.*))?"
	}

	re := regexp.MustCompile("^" + escaped + "$")
	return re.MatchString(str)
}

// MatchAll finds the first matching pattern in patterns (sorted by key length ascending).
// Returns the matched value and true, or nil and false if no match.
func MatchAll(input string, patterns map[string]any) (value any, matched bool) {
	keys := make([]string, 0, len(patterns))
	for k := range patterns {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if len(keys[i]) != len(keys[j]) {
			return len(keys[i]) < len(keys[j])
		}
		return keys[i] < keys[j]
	})

	var result any
	for _, k := range keys {
		if Match(input, k) {
			result = patterns[k]
			matched = true
		}
	}
	return result, matched
}

// MatchCommand matches input against command patterns.
// Input is split by whitespace into words.
// First word matches first pattern segment.
// Remaining words match remaining pattern segments.
// * in pattern segment matches anything.
func MatchCommand(input string, patterns map[string]any) (value any, matched bool) {
	keys := make([]string, 0, len(patterns))
	for k := range patterns {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if len(keys[i]) != len(keys[j]) {
			return len(keys[i]) < len(keys[j])
		}
		return keys[i] < keys[j]
	})

	words := strings.Fields(input)
	if len(words) == 0 {
		words = []string{}
	}

	var result any
	for _, k := range keys {
		parts := strings.Fields(k)
		if len(parts) == 0 {
			continue
		}

		if len(words) == 0 {
			continue
		}

		head := words[0]
		tail := words[1:]

		if !Match(head, parts[0]) {
			continue
		}

		if matchSequence(tail, parts[1:]) {
			result = patterns[k]
			matched = true
		}
	}
	return result, matched
}

// matchSequence matches a sequence of items against pattern segments.
// A pattern segment "*" matches one or more items.
// Other segments must match exactly one item each.
func matchSequence(items, patterns []string) bool {
	if len(patterns) == 0 {
		return true
	}

	pattern, rest := patterns[0], patterns[1:]

	if pattern == "*" {
		// * matches one or more items, then rest matches the rest
		// Try consuming i items (1 to len(items)) and check if rest matches
		for i := 1; i <= len(items); i++ {
			if matchSequence(items[i:], rest) {
				return true
			}
		}
		// Also try * matching zero items (for patterns like "rm *")
		return matchSequence(items, rest)
	}

	// Non-* pattern must match items one by one
	for i := 0; i < len(items); i++ {
		if Match(items[i], pattern) && matchSequence(items[i+1:], rest) {
			return true
		}
	}
	return false
}
