package i18n

import (
	"strings"
)

var rtlLanguages = map[string]bool{
	"ar": true,
	"he": true,
	"ur": true,
	"fa": true,
	"yi": true,
}

func IsRTLLang(tag string) bool {
	return rtlLanguages[strings.ToLower(tag)]
}
