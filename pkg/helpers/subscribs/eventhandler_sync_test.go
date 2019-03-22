package subscribs

import "testing"

func TestEventHandlerOnAndEmit(t *testing.T) {
	evh := &syncEventHandler{}
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
	evh := &syncEventHandler{}
	counter := 0

	handler := func(val ...interface{}) {
		counter += val[0].(int)
	}

	if err := evh.Off(&handler); err != ErrHandlerNotFound {
		t.Error("expected ErrHandlerNotFound, got", err)
	}

	evh.On(&handler)

	if err := evh.Off(&handler); err != nil {
		t.Error("expected nil error, got", err)
	}

	evh.Emit(10)
	if counter != 0 {
		t.Error("expected counter to be 0, got", counter)
	}
}

func TestEventHandlerOffAllAndEmit(t *testing.T) {
	evh := &syncEventHandler{}
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
