package runner

import (
	context "context"
	"errors"
	"testing"
	"time"
)

type mockChannelConsumer struct {
	ch  chan []interface{}
	err error
}

type mockHandlerError struct {
	ch chan error
}

func (mh *mockHandlerError) OnError(err error) error {
	mh.ch <- err
	return nil
}

// Accept .
func (mc *mockChannelConsumer) Accept(ctx context.Context, results ...interface{}) error {
	mc.ch <- results
	return mc.err
}

func TestConcurrentProcessorSuccess(t *testing.T) {
	ctx := context.Background()

	consumer := &mockChannelConsumer{ch: make(chan []interface{})}

	proc := NewConcurrentProcessor(1, 2, 4)
	proc.Start(ctx, consumer, NewPassErrorHandler())

	proc.Enqueue(ctx, "task1")
	proc.Enqueue(ctx, "task2")

	v1 := <-consumer.ch
	if len(v1) != 1 || (v1[0].(string) != "task1" && v1[0].(string) != "task2") {
		t.Error("expect one argument and equal task1 or task2, got", v1)
	}

	v2 := <-consumer.ch
	if len(v2) != 1 || (v2[0].(string) != "task1" && v2[0].(string) != "task2") ||
		v1[0].(string) == v2[0].(string) {
		t.Error("expect one argument and equal task1 or task2, got", v2)
	}

	proc.Stop()
	proc.Wait()
}

func TestConcurrentProcessorError(t *testing.T) {
	ctx := context.Background()

	errHandler := &mockHandlerError{ch: make(chan error)}
	consumer := &mockChannelConsumer{ch: make(chan []interface{})}

	proc := NewConcurrentProcessor(1, 2, 4)
	proc.Start(ctx, consumer, errHandler)

	consumer.err = errors.New("error")
	proc.Enqueue(ctx, "task1")

	<-consumer.ch

	var err error
	select {
	case <-time.After(time.Millisecond * 10):
		t.Error("expect error but timeout happened")
	case err = <-errHandler.ch:
	}

	if err == nil || err.Error() != "error" {
		t.Error("expect error, got", err)
	}

	proc.Stop()
	proc.Wait()
}
