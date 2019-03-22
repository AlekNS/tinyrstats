package monitor

import "errors"

var (
	// ErrTaskNotFound raise when task is not found in the repository.
	ErrTaskNotFound = errors.New("task not found")

	// ErrResourceRequestNotHTTP raise when url not starts with http: or https:
	ErrResourceRequestNotHTTP = errors.New("resource is not http")
)
