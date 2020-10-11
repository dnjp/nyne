package event

import (
	"log"
	"unicode/utf8"
)

// EventHandler listens for acme Events
type EventHandler func(Event) Event

// WinHandler listens for new acme Windows
type WinHandler func(*Win)

// KeyCmdHandler modifies keyboard mappings
type KeyCmdHandler func(Event) Event

// Condition is a function that returns under what condition to run event
type Condition func(Event) bool

// EventHook contains properties for event handlers
type EventHook struct {
	Handler EventHandler
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

// RegisterPHook registers hook on acme 'Put' events
func (a *Acme) RegisterPHook(hook EventHook) {
	if a.debug {
		log.Println("registered Put hook")
	}
	hooks := a.eventHooks[PUT]
	hooks = append(hooks, hook)
	a.eventHooks[PUT] = hooks
}

// RegisterNHook registers the hook on acme 'New' events
func (a *Acme) RegisterNHook(hook WinHook) {
	if a.debug {
		log.Println("registered New hook")
	}
	hooks := a.winHooks[NEW]
	hooks = append(hooks, hook)
	a.winHooks[NEW] = hooks
}

// RedisterKeyCmdHook registers hook for key events
func (a *Acme) RegisterKeyCmdHook(hook KeyCmdHook) {
	if a.debug {
		log.Println("registered key cmd hook")
	}
	a.keyCmdHooks[hook.Key] = &hook
}

// RegisterPHook registers hook on acme 'Put' events
func (f *FileLoop) RegisterPHook(hook EventHook) {
	if f.debug {
		log.Println("registered Put hook")
	}
	hooks := f.eventHooks[PUT]
	hooks = append(hooks, hook)
	f.eventHooks[PUT] = hooks
}

// RegisterNHook registers the hook on acme 'New' events
func (f *FileLoop) RegisterNHook(hook WinHook) {
	if f.debug {
		log.Println("registered New hook")
	}
	hooks := f.winHooks[NEW]
	hooks = append(hooks, hook)
	f.winHooks[NEW] = hooks
}

// RedisterKeyCmdHook registers hook for key events
func (f *FileLoop) RegisterKeyCmdHook(hook KeyCmdHook) {
	if f.debug {
		log.Println("registered key cmd hook")
	}
	f.keyCmdHooks[hook.Key] = &hook
}

func (f *FileLoop) runWinHooks(w *Win) {
	if f.debug {
		log.Println("running win hooks")
	}
	hooks := f.winHooks[NEW]
	if len(hooks) == 0 {
		return
	}

	for _, hook := range hooks {
		fn := hook.Handler
		fn(w)
	}
}

func (f *FileLoop) runKeyCmdHooks(event Event) Event {
	if f.debug {
		log.Println("running key cmd hooks")
	}
	r, _ := utf8.DecodeRune(event.Text)
	keyCmdHook := f.keyCmdHooks[r]
	if keyCmdHook == nil {
		return event
	}
	if f.debug {
		log.Printf("found key cmd for %c", r)
	}
	condition := keyCmdHook.Condition
	if condition(event) {
		if f.debug {
			log.Printf("%c condition met")
		}
		fn := keyCmdHook.Handler
		evt := fn(event)
		return evt
	}
	return event
}

func (f *FileLoop) runEventHooks(event Event) Event {
	if f.debug {
		log.Println("running event hooks")
	}

	hooks := f.eventHooks[event.Builtin]
	if len(hooks) == 0 {
		return event
	}

	// allow progressive mutation of event
	newEvent := event
	for _, hook := range hooks {
		fn := hook.Handler
		newEvent = fn(newEvent)
	}
	return newEvent
}
