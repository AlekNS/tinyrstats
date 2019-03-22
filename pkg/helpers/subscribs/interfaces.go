package subscribs

type (
	// EventHandler .
	EventHandler interface {
		// Emit event data.
		Emit(args ...interface{}) error
		// Subscribe on event data.
		On(*HandlerOnFunc) (HandlerOffFunc, error)
		// Unsubscribe from event data.
		Off(*HandlerOnFunc) error
		// Unsubscribe all subscribers.
		OffAll() error
	}
)
