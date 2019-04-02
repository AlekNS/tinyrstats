package subscribs

import "testing"

func TestEventHandlerOnAndEmit(t *testing.T) {
	evh := NewSyncEventHandler()
	counter := 0

	evh.Emit(10)
	if counter != 0 {
		t.Error("expected counter to be 0, got", counter)
	}

	handler := func(val ...interface{}) {
		counter += val[0].(int)
	}

	evh.On(&handler)

	evh.Emit(10)
	if counter != 10 {
		t.Error("expected counter to be 10, got", counter)
	}
}

func TestEventHandlerOffAndEmit(t *testing.T) {
	evh := NewSyncEventHandler()
	counter := 0

	handler1 := func(val ...interface{}) {
		counter += val[0].(int)
	}

	handler2 := func(val ...interface{}) {
		counter += val[0].(int)+1
	}

	if err := evh.Off(&handler1); err != ErrHandlerNotFound {
		t.Error("expected ErrHandlerNotFound, got", err)
	}

	evh.On(&handler1)
	evh.On(&handler2)

	if err := evh.Off(&handler1); err != nil {
		t.Error("expected nil error, got", err)
	}

	evh.Emit(10)
	if counter != 11 {
		t.Error("expected counter to be 11, got", counter)
	}
}

func TestEventHandlerOffAllAndEmit(t *testing.T) {
	evh := NewSyncEventHandler()
	counter := 0

	handler := func(val ...interface{}) {
		counter += val[0].(int)
	}

	evh.On(&handler)

	if err := evh.OffAll(); err != nil {
		t.Error("expected nil error, got", err)
	}

	evh.Emit(10)
	if counter != 0 {
		t.Error("expected counter to be 0, got", counter)
	}
}
