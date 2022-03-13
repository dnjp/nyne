package event

import (
	"fmt"
	"log"
	"unicode/utf8"
)

// Buf implements the BufListener interface and runs on opened
// acme buffers
type Buf struct {
	id         int
	file       string
	win        *Win
	debug      bool
	eventHooks map[AcmeOp][]Hook
	winHooks   map[AcmeOp][]WinHook
	keyHooks   map[rune]KeyHook
}

// NewBuf constructs an event loop
func NewBuf(id int, file string) *Buf {
	return &Buf{
		id:         id,
		file:       file,
		eventHooks: make(map[AcmeOp][]Hook),
		winHooks:   make(map[AcmeOp][]WinHook),
		keyHooks:   make(map[rune]KeyHook),
	}
}

// File returns the buffers active file
func (b *Buf) File() string {
	return b.file
}

//Win returns the active acme Window
func (b *Buf) Win() *Win {
	return b.win
}

// Start begins the event listener for the window
func (b *Buf) Start() error {
	w, err := OpenWin(b.id, b.file)
	if err != nil {
		return err
	}
	b.win = w

	// runs hooks for acme 'new' event
	b.winEvent(b.win, Event{Builtin: New})

	for e := range b.win.OpenEventChan() {
		if b.debug {
			log.Printf("EVENT(before): %+v\n", *e)
		}

		event, err := NewEvent(e, b.id, b.file)
		if err != nil {
			return err
		}

		if event.Origin == Keyboard {
			w.Lastpoint = event.SelBegin
			event = b.keyEvent(event)
		} else {
			if event.Origin == DelOrigin && event.Type == DelType {
				b.win.WriteEvent(event)
				b.win.Close()
				return nil
			}
			event = b.execEvent(event)
		}

		if b.debug {
			log.Printf("EVENT(after): %+v\n", event)
		}

		b.win.WriteEvent(event)

		for _, h := range event.PostHooks {
			if err := h(event); err != nil {
				return err
			}
		}

		// maintain current address after formatting buffer
		if event.Builtin == Put {
			body, err := b.win.ReadBody()
			if err != nil {
				return err
			}
			if len(body) < w.Lastpoint {
				w.Lastpoint = len(body)
			}
			if err := b.win.SetAddr(fmt.Sprintf("#%d", w.Lastpoint)); err != nil {
				return err
			}
			if err := b.win.SetTextToAddr(); err != nil {
				return err
			}
			if err := b.win.ExecShow(); err != nil {
				return err
			}
		}
	}
	return nil
}

// RegisterHook registers hook on acme 'Put' events
func (b *Buf) RegisterHook(hook Hook) {
	hooks := b.eventHooks[Put]
	hooks = append(hooks, hook)
	b.eventHooks[Put] = hooks
}

// RegisterWinHook registers the hook on acme 'New' events
func (b *Buf) RegisterWinHook(hook WinHook) {
	hooks := b.winHooks[New]
	hooks = append(hooks, hook)
	b.winHooks[New] = hooks
}

// RegisterKeyHook registers hook for key events
func (b *Buf) RegisterKeyHook(hook KeyHook) {
	b.keyHooks[hook.Key] = hook
}

func (b *Buf) winEvent(w *Win, event Event) {
	for _, hook := range b.winHooks[event.Builtin] {
		hook.Handler(w)
	}
}

func (b *Buf) keyEvent(event Event) Event {
	r, _ := utf8.DecodeRune(event.Text)
	hook, ok := b.keyHooks[r]
	if !ok {
		return event
	}

	if hook.Condition(event) {
		return hook.Handler(event)
	}

	return event
}

func (b *Buf) execEvent(event Event) Event {
	newEvent := event
	for _, hook := range b.eventHooks[event.Builtin] {
		newEvent = hook.Handler(newEvent)
	}
	return newEvent
}
