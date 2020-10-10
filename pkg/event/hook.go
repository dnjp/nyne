package event

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
	hooks := a.eventHooks[PUT]
	hooks = append(hooks, hook)
	a.eventHooks[PUT] = hooks
}

// RegisterNHook registers the hook on acme 'New' events
func (a *Acme) RegisterNHook(hook WinHook) {
	hooks := a.winHooks[NEW]
	hooks = append(hooks, hook)
	a.winHooks[NEW] = hooks
}

func (a *Acme) runNHooks(w *Win) {
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
