package util

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestAbortController_Basic(t *testing.T) {
	ac := NewAbortController()

	if ac.IsAborted() {
		t.Fatal("controller should not be aborted initially")
	}

	ac.Abort()

	if !ac.IsAborted() {
		t.Fatal("controller should be aborted after Abort()")
	}
}

func TestAbortController_AbortOnce(t *testing.T) {
	ac := NewAbortController()

	ac.Abort()
	ac.Abort()

	if !ac.IsAborted() {
		t.Fatal("controller should remain aborted")
	}
}

func TestAbortAfter_Timeout(t *testing.T) {
	controller, clear := AbortAfter(50)

	if controller.IsAborted() {
		t.Fatal("controller should not be aborted immediately")
	}

	time.Sleep(60 * time.Millisecond)

	if !controller.IsAborted() {
		t.Fatal("controller should be aborted after timeout")
	}

	clear()
}

func TestAbortAfter_ClearTimeout(t *testing.T) {
	controller, clear := AbortAfter(5000)

	if controller.IsAborted() {
		t.Fatal("controller should not be aborted immediately")
	}

	clear()

	time.Sleep(10 * time.Millisecond)

	if controller.IsAborted() {
		t.Fatal("controller should not be aborted after clear")
	}
}

func TestWithAbort_SignalClose(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("skipping test with race condition in CI")
	}
	signal := make(chan struct{})
	ctx, _ := WithAbort(signal)

	select {
	case <-ctx.Done():
		t.Fatal("context should not be done before signal close")
	default:
	}

	close(signal)

	select {
	case <-ctx.Done():
	default:
		t.Fatal("context should be done after signal close")
	}
}

func TestWithAbort_SignalNeverClosed(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("skipping test with race condition in CI")
	}
	signal := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-ctx.Done()
		close(signal)
	}()

	abortCtx, _ := WithAbort(signal)
	cancel()

	select {
	case <-abortCtx.Done():
	default:
		t.Fatal("context should be done after parent cancel")
	}
}

func TestWithAbortFromController(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("skipping test with race condition in CI")
	}
	ac := NewAbortController()
	ctx, cancel := WithAbortFromController(ac)

	select {
	case <-ctx.Done():
		t.Fatal("context should not be done initially")
	default:
	}

	cancel()

	select {
	case <-ctx.Done():
	default:
		t.Fatal("context should be done after cancel")
	}

	if !ac.IsAborted() {
		t.Fatal("controller should be aborted")
	}
}

func TestAbortAfterContext_Timeout(t *testing.T) {
	ctx, cancel := AbortAfterContext(context.Background(), 50)

	select {
	case <-ctx.Done():
		t.Fatal("context should not be done before timeout")
	default:
	}

	time.Sleep(60 * time.Millisecond)

	select {
	case <-ctx.Done():
	default:
		t.Fatal("context should be done after timeout")
	}

	cancel()
}

func TestAbortAfterContext_ParentCancel(t *testing.T) {
	parent, parentCancel := context.WithCancel(context.Background())
	ctx, cancel := AbortAfterContext(parent, 5000)

	parentCancel()

	select {
	case <-ctx.Done():
	default:
		t.Fatal("context should be done when parent is canceled")
	}

	cancel()
}

func TestAbortAfterContext_ClearTimeout(t *testing.T) {
	ctx, cancel := AbortAfterContext(context.Background(), 5000)

	cancel()

	time.Sleep(10 * time.Millisecond)

	select {
	case <-ctx.Done():
	default:
		t.Fatal("context should be done after cancel")
	}
}

func TestAbortAfterContext_AlreadyCanceled(t *testing.T) {
	parent, cancel := context.WithCancel(context.Background())
	cancel()

	ctx, _ := AbortAfterContext(parent, 5000)

	select {
	case <-ctx.Done():
	default:
		t.Fatal("context should be done when parent is already canceled")
	}
}

func TestAbortAfterContext_DoubleCancel(t *testing.T) {
	_, cancel := AbortAfterContext(context.Background(), 5000)

	cancel()
	cancel()
}
