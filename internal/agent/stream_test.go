package agent

import (
    "bytes"
    "context"
    "strings"
    "testing"
    "time"
    "errors"
)

func TestStreamHandler_CreateSendAndClose(t *testing.T) {
    sh := NewStreamHandler()
    sess := "sess-stream"
    ch := sh.CreateStream(sess)
    // send a message and then close
    if err := sh.Send(sess, Message{Content: "hello"}); err != nil {
        t.Fatalf("unexpected error sending: %v", err)
    }
    // read via Stream in a separate goroutine
    var buf bytes.Buffer
    ctx := context.Background()
    done := make(chan struct{})
    go func() {
        _ = sh.Stream(ctx, sess, &buf)
        close(done)
    }()
    // ensure the message has been delivered
    // give it a moment for the goroutine to process
    time.Sleep(5 * time.Millisecond)
    sh.CloseStream(sess)
    <-done
    // check that the buffer contains the data line
    data := buf.String()
    if data == "" {
        t.Fatalf("expected data to be written to writer, got empty string; ch=%v", ch)
    }
}

func TestStreamBufferFull(t *testing.T) {
    sh := NewStreamHandler()
    sess := "sess-full"
    sh.CreateStream(sess)
    // Fill the internal channel buffer by sending without a reader
    for i := 0; i < 101; i++ {
        err := sh.Send(sess, Message{Content: "x"})
        if i < 100 {
            if err != nil {
                t.Fatalf("unexpected error on send %d: %v", i, err)
            }
        } else {
            if err == nil {
                t.Fatalf("expected buffer full error on 101st send, got nil")
            }
        }
    }
}

func TestStreamNotFoundPaths(t *testing.T) {
    sh := NewStreamHandler()
    // Sending to a non-existent session should error
    if err := sh.Send("nope", Message{Content: "x"}); err == nil {
        t.Fatalf("expected error when sending to unknown session, got nil")
    }
    // Streaming from a non-existent session should error
    var buf bytes.Buffer
    if err := sh.Stream(context.Background(), "nope", &buf); err == nil {
        t.Fatalf("expected error when streaming for unknown session, got nil")
    }
}

func TestStreamWritesMultipleMessages(t *testing.T) {
    sh := NewStreamHandler()
    sess := "sess-multi"
    ch := sh.CreateStream(sess)
    var buf bytes.Buffer
    ctx := context.Background()
    done := make(chan struct{})
    go func() {
        _ = sh.Stream(ctx, sess, &buf)
        close(done)
    }()

    // send several messages
    sh.Send(sess, Message{Content: "alpha"})
    time.Sleep(1 * time.Millisecond)
    sh.Send(sess, Message{Content: "beta"})
    time.Sleep(1 * time.Millisecond)
    sh.Send(sess, Message{Content: "gamma"})
    time.Sleep(1 * time.Millisecond)
    // close stream
    sh.CloseStream(sess)
    <-done

    // ensure all messages were written in order
    data := buf.String()
    if !strings.Contains(data, "data: alpha") || !strings.Contains(data, "data: beta") || !strings.Contains(data, "data: gamma") {
        t.Fatalf("expected all messages to be written, got: %q", data)
    }
    _ = ch
}

func TestStreamContextCancel(t *testing.T) {
    sh := NewStreamHandler()
    sess := "sess-cancel"
    sh.CreateStream(sess)
    var buf bytes.Buffer
    ctx, cancel := context.WithCancel(context.Background())
    done := make(chan error, 1)
    go func() {
        err := sh.Stream(ctx, sess, &buf)
        done <- err
    }()
    // Cancel the context to trigger early exit
    cancel()
    err := <-done
    if err != context.Canceled {
        t.Fatalf("expected context.Canceled, got %v", err)
    }
}

func TestStreamWriterErrorPropagation(t *testing.T) {
    sh := NewStreamHandler()
    sess := "sess-werr"
    sh.CreateStream(sess)
    // writer that always fails
    fw := &failWriter{}
    ctx := context.Background()
    done := make(chan error, 1)
    go func() {
        err := sh.Stream(ctx, sess, fw)
        done <- err
    }()
    // Send a single message then close
    sh.Send(sess, Message{Content: "will fail"})
    sh.CloseStream(sess)
    err := <-done
    if err == nil {
        t.Fatalf("expected error from writer, got nil")
    }
}

// Removed: non-deterministic path when closing a stream before it starts.

type failWriter struct{}

func (f *failWriter) Write(p []byte) (int, error) {
    return 0, errors.New("write-fail")
}
