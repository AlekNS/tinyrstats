package subscribs

type (
	// HandlerOnFunc .
	HandlerOnFunc = func(args ...interface{})

	// HandlerOffFunc .
	HandlerOffFunc = func()
)
