package event

import (
	"fmt"
	"log"
	"unicode/utf8"
)

// Buf implements the BufListener interface and runs on opened
// acme buffers
type Buf struct {
	id          int
	file        string
	win         *Win
	debug       bool
	eventHooks  map[AcmeOp][]PutHook
	winHooks    map[AcmeOp][]WinHook
	keyCmdHooks map[rune]*KeyCmdHook
}

// NewBuf constructs an event loop
func NewBuf(id int, file string) *Buf {
	return &Buf{
		id:          id,
		file:        file,
		eventHooks:  make(map[AcmeOp][]PutHook),
		winHooks:    make(map[AcmeOp][]WinHook),
		keyCmdHooks: make(map[rune]*KeyCmdHook),
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
	b.runWinHooks(b.win)

	for e := range b.win.OpenEventChan() {
		if b.debug {
			log.Printf("RAW: %+v\n", *e)
		}

		event, err := NewEvent(e, b.id, b.file)
		if err != nil {
			return err
		}

		if event.Origin == Keyboard {
			w.Lastpoint = event.SelBegin
			event = b.runKeyCmdHooks(event)
		} else {
			if event.Origin == DelOrigin && event.Type == DelType {
				b.win.WriteEvent(event)
				b.win.Close()
				return nil
			}
			event = b.runPutHooks(event)
		}

		if b.debug {
			log.Printf("TOKEN: %+v\n", event)
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

// RegisterPutHook registers hook on acme 'Put' events
func (b *Buf) RegisterPutHook(hook PutHook) {
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

// RegisterKeyCmdHook registers hook for key events
func (b *Buf) RegisterKeyCmdHook(hook KeyCmdHook) {
	b.keyCmdHooks[hook.Key] = &hook
}

func (b *Buf) runWinHooks(w *Win) {
	hooks := b.winHooks[New]
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
