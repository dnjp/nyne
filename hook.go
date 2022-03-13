package nyne

// Condition is a function that returns under what condition to run event
type Condition func(Event) bool

// Handler transforms an event
type Handler func(Event) Event

// WinHandler transforms the window
type WinHandler func(*Win)
