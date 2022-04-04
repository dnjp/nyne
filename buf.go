package nyne

import (
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
	EventHooks map[Builtin][]Handler
	WinHooks   map[Builtin][]WinHandler
	KeyHooks   map[rune]Handler
}

// NewBuf constructs an event loop
func NewBuf(id int, file string) *Buf {
	return &Buf{
		id:         id,
		file:       file,
		EventHooks: make(map[Builtin][]Handler),
		WinHooks:   make(map[Builtin][]WinHandler),
		KeyHooks:   make(map[rune]Handler),
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
	stop := make(chan struct{})
	defer func() { stop <- struct{}{} }()
	events, errs := b.win.EventChan(b.id, b.file, stop)
	for {
		select {
		case event := <-events:

			var ok bool
			if event.Origin == Keyboard {
				w.Lastpoint = event.SelBegin
				event, ok = b.keyEvent(event)
			} else {
				if event.Origin == Delete && event.Action == DelType {
					b.win.WriteEvent(event)
					b.win.Close()
					return nil
				}
				event, ok = b.execEvent(event)
			}
			if !ok {
				continue
			}

			if b.debug {
				log.Printf("EVENT(after): %+v\n", event)
			}

			b.win.WriteEvent(event)

			for _, h := range event.WriteHooks {
				if err := h(event); err != nil {
					return err
				}
			}

			// maintain current address after formatting buffer
			if event.Builtin == Put {
				body, err := b.win.Body()
				if err != nil {
					return err
				}
				if len(body) < w.Lastpoint {
					w.Lastpoint = len(body)
				}
				if err := b.win.SetAddr("#%d", w.Lastpoint); err != nil {
					return err
				}
				if err := b.win.SelectionFromAddr(); err != nil {
					return err
				}
				if err := b.win.Show(); err != nil {
					return err
				}
			}

		case err := <-errs:
			return err
		}
	}
}

func (b *Buf) winEvent(w *Win, event Event) {
	for _, hook := range b.WinHooks[event.Builtin] {
		hook(w)
	}
}

func (b *Buf) keyEvent(event Event) (Event, bool) {
	r, _ := utf8.DecodeRune(event.Text)
	hook, ok := b.KeyHooks[r]
	if !ok {
		return event, true
	}
	return hook(event)
}

func (b *Buf) execEvent(event Event) (Event, bool) {
	origEvent := event
	newEvent := origEvent
	ok := true
	for _, hook := range b.EventHooks[event.Builtin] {
		newEvent, ok = hook(newEvent)
		if !ok {
			return origEvent, true
		}
	}
	return newEvent, ok
}
