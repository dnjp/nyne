package event

import (
	"unicode/utf8"
)

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

// RegisterPutHook registers hook on acme 'Put' events
func (a *Acme) RegisterPutHook(hook PutHook) {
	hooks := a.eventHooks[PUT]
	hooks = append(hooks, hook)
	a.eventHooks[PUT] = hooks
}

// RegisterWinHook registers the hook on acme 'New' events
func (a *Acme) RegisterWinHook(hook WinHook) {
	hooks := a.winHooks[NEW]
	hooks = append(hooks, hook)
	a.winHooks[NEW] = hooks
}

// RegisterKeyCmdHook registers hook for key events
func (a *Acme) RegisterKeyCmdHook(hook KeyCmdHook) {
	a.keyCmdHooks[hook.Key] = &hook
}

// RegisterPutHook registers hook on acme 'Put' events
func (b *Buf) RegisterPutHook(hook PutHook) {
	hooks := b.eventHooks[PUT]
	hooks = append(hooks, hook)
	b.eventHooks[PUT] = hooks
}

// RegisterWinHook registers the hook on acme 'New' events
func (b *Buf) RegisterWinHook(hook WinHook) {
	hooks := b.winHooks[NEW]
	hooks = append(hooks, hook)
	b.winHooks[NEW] = hooks
}

// RegisterKeyCmdHook registers hook for key events
func (b *Buf) RegisterKeyCmdHook(hook KeyCmdHook) {
	b.keyCmdHooks[hook.Key] = &hook
}

func (b *Buf) runWinHooks(w *Win) {
	hooks := b.winHooks[NEW]
	if len(hooks) == 0 {
		return
	}
	for _, hook := range hooks {
		fn := hook.Handler
		fn(w)
	}
}

func (b *Buf) runKeyCmdHooks(event Event) Event {
	r, _ := utf8.DecodeRune(event.Text)
	keyCmdHook := b.keyCmdHooks[r]
	if keyCmdHook == nil {
		return event
	}
	if keyCmdHook.Condition(event) {
		evt := keyCmdHook.Handler(event)
		return evt
	}
	return event
}

func (b *Buf) runPutHooks(event Event) Event {
	hooks := b.eventHooks[event.Builtin]
	if len(hooks) == 0 {
		return event
	}
	// allow progressive mutation of event
	newEvent := event
	for _, hook := range hooks {
		newEvent = hook.Handler(newEvent)
	}
	return newEvent
}
