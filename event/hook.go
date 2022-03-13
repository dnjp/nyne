package event

// PutHandler listens for acme Events
type PutHandler func(Event) Event

// WinHandler listens for new acme Windows
type WinHandler func(*Win)

// KeyCmdHandler modifies keyboard mappings
type KeyCmdHandler func(Event) Event

// Condition is a function that returns under what condition to run event
type Condition func(Event) bool

// PutHook contains properties for event handlers
type PutHook struct {
	Handler PutHandler
}

// WinHook contains properties for window handlers
type WinHook struct {
	Handler WinHandler
}

// KeyCmdHook contains properties for key handlers
type KeyCmdHook struct {
	Key       rune
	Condition Condition
	Handler   KeyCmdHandler
}
