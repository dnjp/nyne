package event

import "log"

// EventHandler listens for acme Events
type EventHandler func(*Event) *Event

// WinHandler listens for new acme Windows
type WinHandler func(*Win)

// EventHook contains properties for event handlers
type EventHook struct {
	Handler EventHandler
}

// WinHook contains properties for window handlers
type WinHook struct {
	Handler WinHandler
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

func (a *Acme) runWinHooks(w *Win) {
	if a.debug {
		log.Println("running win hooks")
	}
	hooks := a.winHooks[NEW]
	if len(hooks) == 0 {
		return
	}

	for _, hook := range hooks {
		fn := hook.Handler
		fn(w)
	}
}

func (a *Acme) runEventHooks(event *Event) *Event {
	if a.debug {
		log.Println("running event hooks")
	}
	if event.Builtin == nil {
		return event
	}

	hooks := a.eventHooks[*event.Builtin]
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
