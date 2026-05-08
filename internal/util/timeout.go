package util

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"time"
)

func WithTimeout[T any](promise func() (T, error), ms int) (T, error) {
	var zero T
	if ms <= 0 {
		return promise()
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(ms)*time.Millisecond)
	defer cancel()

	resultCh := make(chan T, 1)
	errCh := make(chan error, 1)
	doneCh := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			return
		case <-doneCh:
			return
		default:
		}

		result, err := promise()
		select {
		case <-ctx.Done():
			return
		case <-doneCh:
			return
		case resultCh <- result:
			if err != nil {
				select {
				case errCh <- err:
				case <-ctx.Done():
				case <-doneCh:
				}
			}
		}
	}()

	select {
	case <-ctx.Done():
		wg.Wait()
		return zero, errors.New("Operation timed out after " + strconv.Itoa(ms) + "ms")
	case err := <-errCh:
		<-resultCh
		return zero, err
	case result := <-resultCh:
		close(doneCh)
		return result, nil
	}
}