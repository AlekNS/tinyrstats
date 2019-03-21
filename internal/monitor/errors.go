package worker

import "errors"

var (
	// ErrTaskNotFound raise when task is not found in the repository.
	ErrTaskNotFound = errors.New("task not found")
)
