package util

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
)

func RandomID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func TrimSpace(s string) string {
	return strings.TrimSpace(s)
}

func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func Split(s, sep string) []string {
	return strings.Split(s, sep)
}

func Join(elems []string, sep string) string {
	return strings.Join(elems, sep)
}

func ToLower(s string) string {
	return strings.ToLower(s)
}

func ToUpper(s string) string {
	return strings.ToUpper(s)
}

func Replace(s, old, new string) string {
	return strings.ReplaceAll(s, old, new)
}

func Matches(s, pattern string) (bool, error) {
	matched, err := regexp.MatchString(pattern, s)
	return matched, err
}

func ParseInt(s string) (int64, error) {
	var n int64
	fmt.Sscanf(s, "%d", &n)
	return n, nil
}

func Hash(s string) uint64 {
	var h uint64
	for _, c := range s {
		h = h*31 + uint64(c)
	}
	return h
}
