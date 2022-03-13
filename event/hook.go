package event

// Condition is a function that returns under what condition to run event
type Condition func(Event) bool

// Handler transforms an event
type Handler func(Event) Event

// Hook contains properties for event handlers
type Hook struct {
	Op      AcmeOp
	Handler Handler
}

// WinHook contains properties for window handlers
type WinHook struct {
	Op      AcmeOp
	Handler func(*Win)
}

// KeyHook contains properties for key handlers
type KeyHook struct {
	Key       rune
	Condition Condition
	Handler   Handler
}
