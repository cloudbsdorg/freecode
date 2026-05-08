package util

import (
	"context"
	"errors"
	"strconv"
	"time"
)

func WithTimeout[T any](promise func() (T, error), ms int) (T, error) {
	var zero T
	if ms <= 0 {
		return promise()
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(ms)*time.Millisecond)
	defer cancel()

	type resultType struct {
		result T
		err    error
	}
	resultCh := make(chan resultType, 1)

	go func() {
		result, err := promise()
		select {
		case resultCh <- resultType{result: result, err: err}:
		case <-ctx.Done():
		}
	}()

	select {
	case <-ctx.Done():
		return zero, errors.New("Operation timed out after " + strconv.Itoa(ms) + "ms")
	case rt := <-resultCh:
		return rt.result, rt.err
	}
}