package subscribs

import "sync"

// syncEventHandler works only in sync sequential mode.
type syncEventHandler struct {
	mtx sync.RWMutex

	handlers []*HandlerOnFunc
}

// Emit event data.
func (o *syncEventHandler) Emit(args ...interface{}) error {
	o.mtx.RLock()
	handlersSnapshot := append([]*HandlerOnFunc(nil), o.handlers...)
	defer o.mtx.RUnlock()

	// Sync sequential call handlers.
	for _, obs := range handlersSnapshot {
		(*obs)(args...)
	}

	return nil
}

// On subscribes for event.
func (o *syncEventHandler) On(fn *HandlerOnFunc) (HandlerOffFunc, error) {
	o.mtx.Lock()

	subscriberIndex := len(o.handlers)
	o.handlers = append(o.handlers, fn)

	o.mtx.Unlock()

	return func() {
		o.mtx.Lock()
		defer o.mtx.Unlock()
		o.remove(subscriberIndex)
	}, nil
}

func (o *syncEventHandler) remove(index int) {
	// order is not matter
	handlersCount := len(o.handlers)
	if handlersCount > 1 {
		lastHandler := o.handlers[handlersCount]
		o.handlers[index] = lastHandler
		o.handlers = o.handlers[:handlersCount-1]
	} else {
		o.handlers = []*HandlerOnFunc(nil)
	}
}

// Off remove event handler.
func (o *syncEventHandler) Off(fn *HandlerOnFunc) error {
	o.mtx.Lock()
	defer o.mtx.Unlock()

	handlersCount := len(o.handlers)
	for index := 0; index < handlersCount; index++ {
		if o.handlers[index] == fn {
			o.remove(index)
			return nil
		}
	}

	return ErrHandlerNotFound
}

// Off remove all event handlers.
func (o *syncEventHandler) OffAll() error {
	o.mtx.Lock()
	o.handlers = []*HandlerOnFunc(nil)
	o.mtx.Unlock()

	return nil
}

// NewSyncEventHandler creates sync event handler.
func NewSyncEventHandler() EventHandler {
	return &syncEventHandler{}
}
