package subscribs

import "errors"

var (
	// ErrHandlerNotFound raises when handler func not found.
	ErrHandlerNotFound = errors.New("event handler not found")
)
