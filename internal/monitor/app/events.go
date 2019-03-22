package app

import (
	"github.com/alekns/tinyrstats/internal/monitor"
	"github.com/alekns/tinyrstats/pkg/helpers/subscribs"
)

type eventsImpl struct {
	taskQueriedByURL subscribs.EventHandler

	taskQueriedByMinResponse subscribs.EventHandler
	taskQueriedByMaxResponse subscribs.EventHandler
}

func (e *eventsImpl) TaskQueriedByURL() subscribs.EventHandler {
	return e.taskQueriedByURL
}

func (e *eventsImpl) TaskQueriedByMinResponse() subscribs.EventHandler {
	return e.taskQueriedByMinResponse
}

func (e *eventsImpl) TaskQueriedByMaxResponse() subscribs.EventHandler {
	return e.taskQueriedByMaxResponse
}

// NewSyncEventsImpl creates sync events
func NewSyncEventsImpl() monitor.Events {
	return &eventsImpl{
		taskQueriedByURL: subscribs.NewSyncEventHandler(),

		taskQueriedByMinResponse: subscribs.NewSyncEventHandler(),
		taskQueriedByMaxResponse: subscribs.NewSyncEventHandler(),
	}
}
