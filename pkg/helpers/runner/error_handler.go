package runner

// passErrorHandler is not using a error and only pass it through your self.
type passErrorHandler struct {
}

// OnError pass.
func (eh *passErrorHandler) OnError(err error) error {
	return err
}

// NewPassErrorHandler create a pass through error handler.
func NewPassErrorHandler() ErrorHandler {
	return &passErrorHandler{}
}
