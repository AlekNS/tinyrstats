package webhttp

import "errors"

var (
	// Errors raised by validators.

	// ErrNoParameters .
	ErrNoParameters = errors.New("no parameters was passed")
	// ErrTooManyParameters .
	ErrTooManyParameters = errors.New("too many request parameters")
	// ErrInvalidParameterValue .
	ErrInvalidParameterValue = errors.New("invalid parameter value")
)
