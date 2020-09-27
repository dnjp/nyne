package event

type Handler func(*Event) *Event

type OpenHandler func(*Win)

type Hook struct {
	Op      AcmeOp
	Handler Handler
}

type OpenHook struct {
	Op AcmeOp
	Handler OpenHandler
}

func (a *Acme) RegisterHook(hook Hook) {
	hooks := a.hooks[hook.Op]
	hooks = append(hooks, hook)
	a.hooks[hook.Op] = hooks
}

func (a *Acme) RegisterOpenHook(hook OpenHook) {
	hooks := a.openHooks[hook.Op]
	hooks = append(hooks, hook)
	a.openHooks[hook.Op] = hooks
}

func (a *Acme) runOpenHooks(w *Win) {
	hooks := a.openHooks[NEW]
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

	hooks := a.hooks[*event.Builtin]
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