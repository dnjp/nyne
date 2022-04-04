package nyne

// Hook executes a function on the event after it has been
// written to the acme log
type Hook func(Event) error

// Condition is a function that returns under what condition to run event
type Condition func(Event) bool

// Handler transforms an event
type Handler func(Event) (Event, bool)

// WinHandler transforms the window
type WinHandler func(*Win)
