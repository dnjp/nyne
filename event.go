package nyne

import (
	"9fans.net/go/acme"
)

// Hook executes a function on the event after it has been
// written to the acme log
type Hook func(Event) error

// Event contains metadata for each Acme event
//
// The message includes the text if it is less than 256 chars. If it
// is longer than that, the fourth number is 0 and the body must be read
// through the data file
type Event struct {
	// Log
	ID                       int
	File                     string
	Origin                   Origin
	Action                   Action
	Text                     Text
	Flag                     Flag
	SelBegin, SelEnd         int
	OrigSelBegin, OrigSelEnd int
	NumBytes                 int
	NumRunes                 int
	ChordArg                 []byte
	ChordLoc                 []byte
	// Hooks
	WriteHooks []Hook
}

// Text contains the default acme event types
type Text string

const (
	// New represents window creation
	New Text = "New"
	// Zerox reprents window creation via zerox
	Zerox Text = "Zerox"
	// Get loads/reloads the file in the window
	Get Text = "Get"
	// Put writes window to the named file
	Put Text = "Put"
	// Del deletes the window
	Del Text = "Del"
	// Focus is received when the window focused
	Focus Text = "Focus"
)

// NewText constructs a builtin from the event text
func NewText(text []byte) Text {
	return Text(text)
}

// String returns Text as a string
func (b Text) String() string {
	return string(b)
}

// Bytes returns Text as a slice of bytes
func (b Text) Bytes() []byte {
	return []byte(b)
}

// Origin is the entity that originated the action
type Origin rune

const (
	// BodyOrTag represents an action received in the body or tag
	BodyOrTag Origin = 'E'
	// WindowFiles represents an action taken in the file
	WindowFiles Origin = 'F'
	// Keyboard represents an action taken by the keyboard
	Keyboard Origin = 'K'
	// Mouse represents an action taken by the mouse
	Mouse Origin = 'M'
	// Delete represents a delete event
	Delete Origin = 0x0
)

// NewOrigin constructs an origin from c1
func NewOrigin(c1 rune) Origin {
	return Origin(c1)
}

// Rune returns the Origin as a rune
func (o Origin) Rune() rune {
	if o == Delete {
		return 0
	}
	return rune(o)
}

// Action describes what kind of action was taken
type Action rune

const (
	// BodyDelete is a deletion in the window body
	BodyDelete Action = 'D'
	// TagDelete is a deletion in the window tag
	TagDelete Action = 'd'
	// BodyInsert is an insertion into the window body
	BodyInsert Action = 'I'
	// TagInsert is an insert into the window tag
	TagInsert Action = 'i'
	// B3Body is a right click event in the window body
	B3Body Action = 'L'
	// B3Tag is a right click event in the window tag
	B3Tag Action = 'l'
	// B2Body is a middle click event in the window body
	B2Body Action = 'X'
	// B2Tag is a middle click event in the window tag
	B2Tag Action = 'x'
	// DelType represents a delete event
	DelType Action = 0x0
)

// NewAction constructs a new action
func NewAction(c2 rune) Action {
	return Action(c2)
}

// Rune returns the Action as a rune
func (a Action) Rune() rune {
	if a == DelType {
		return 0
	}
	return rune(a)
}

// Flag contains the flag for the event. For BodyDelete, TagDelete,
// BodyInsert, and TagInsert the flag is always zero. For messages with
// the 1 bit on in the flag, writing the message back to the event file,
// but with the flag, count, and text omitted, will cause the action to be
// applied to the file exactly as it would have been if the event file had
// not been open.
type Flag int

const (
	// when action is B2Body and B2Tag

	// IsBuiltin represents a built-in command
	IsBuiltin Flag = iota

	// IsNull represents if the text is a null string that has a
	// non-null expansion; if so, another complete message will
	// follow describing the expansion exactly as if it had been
	// indicated explicitly (its flag will always be 0)
	IsNull

	// HasChordedArg says if the command has an extra (chorded)
	// argument; if so, two more complete messages will follow
	// reporting the argument (with all numbers 0 except the
	// character count) and where it originated, in the form of a
	// fully-qualified button 3 style address.
	HasChordedArg

	// when action is B3Body or B3Tag

	// NoReloadNeeded says if acme can interpret the action without
	// loading a new file
	NoReloadNeeded

	// PostExpandFollows says if a second (post-expansion) message
	// follows, analogous to that with X messages
	PostExpandFollows

	// IsFileOrWindow says If the text is a file or window name
	// (perhaps with address) rather than plain literal text.
	IsFileOrWindow

	// For messages with the 1 bit on in the flag, writing the message
	// back to the event file, but with the flag, count, and text omitted,
	// will cause the action to be applied to the file exactly as it would
	// have been if the event file had not been open.
)

// NewFlag constructs a Flag
func NewFlag(a Action, rawFlag int) Flag {
	if a == B2Body || a == B2Tag {
		switch rawFlag {
		case 1:
			return IsBuiltin
		case 2:
			return IsNull
		case 3, 8:
			return HasChordedArg
		}
	}

	if a == B3Body || a == B3Tag {
		switch rawFlag {
		case 1:
			return NoReloadNeeded
		case 2:
			return PostExpandFollows
		case 4:
			return IsFileOrWindow
		}
	}
	return Flag(rawFlag)
}

// NewEvent constructs an Event from a raw acme event
func NewEvent(event *acme.Event, id int, file string) (Event, error) {
	e := Event{
		ID:           id,
		File:         file,
		Origin:       NewOrigin(event.C1),
		Action:       NewAction(event.C2),
		Text:         NewText(event.Text),
		SelBegin:     event.Q0,
		SelEnd:       event.Q1,
		OrigSelBegin: event.OrigQ0,
		OrigSelEnd:   event.OrigQ1,
		NumBytes:     event.Nb,
		NumRunes:     event.Nr,
		ChordArg:     event.Arg,
		ChordLoc:     event.Loc,
	}
	e.Flag = NewFlag(e.Action, event.Flag)
	return e, nil
}

// Log returns a raw acme log event for the Event type
func (e *Event) Log() (*acme.Event, error) {
	return &acme.Event{
		C1:     e.Origin.Rune(),
		C2:     e.Action.Rune(),
		Q0:     e.SelBegin,
		Q1:     e.SelEnd,
		OrigQ0: e.OrigSelBegin,
		OrigQ1: e.OrigSelEnd,
		Flag:   int(e.Flag),
		Nb:     e.NumBytes,
		Nr:     e.NumRunes,
		Text:   e.Text.Bytes(),
		Arg:    e.ChordArg,
		Loc:    e.ChordLoc,
	}, nil
}
