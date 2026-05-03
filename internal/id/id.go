package id

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"sync"
	"time"
)

var counter uint64
var mu sync.Mutex

func New() string {
	mu.Lock()
	counter++
	mu.Unlock()
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), counter)
}

func Random() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func Short() string {
	b := make([]byte, 8)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
