package app

import (
	"github.com/alekns/tinyrstats/internal/monitor"
	"github.com/alekns/tinyrstats/pkg/helpers/subscribs"
)

// eventsImpl implementation of Events
type eventsImpl struct {
	taskQueriedByResource subscribs.EventHandler

	taskQueriedByMinResponse subscribs.EventHandler
	taskQueriedByMaxResponse subscribs.EventHandler
}

// TaskQueriedByResource .
func (e *eventsImpl) TaskQueriedByResource() subscribs.EventHandler {
	return e.taskQueriedByResource
}

// TaskQueriedByMinResponse .
func (e *eventsImpl) TaskQueriedByMinResponse() subscribs.EventHandler {
	return e.taskQueriedByMinResponse
}

// TaskQueriedByMaxResponse .
func (e *eventsImpl) TaskQueriedByMaxResponse() subscribs.EventHandler {
	return e.taskQueriedByMaxResponse
}

// NewSyncEventsImpl creates sync events.
func NewSyncEventsImpl() monitor.Events {
	return &eventsImpl{
		taskQueriedByResource: subscribs.NewSyncEventHandler(),

		taskQueriedByMinResponse: subscribs.NewSyncEventHandler(),
		taskQueriedByMaxResponse: subscribs.NewSyncEventHandler(),
	}
}
