package nyne

import (
	"fmt"
	"strings"

	"9fans.net/go/acme"
)

// Builtin contains the default acme event types
type Builtin int

const (
	// New represents window creation
	New Builtin = iota
	// Zerox reprents window creation via zerox
	Zerox
	// Get loads/reloads the file in the window
	Get
	// Put writes window to the named file
	Put
	// Del deletes the window
	Del
	// Focus is received when the window focused
	Focus
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
	Text                     []byte
	Builtin                  Builtin
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

// Origin is the entity that originated the action
type Origin int

const (
	// BodyOrTag represents an action received in the body or tag
	BodyOrTag Origin = iota
	// WindowFiles represents an action taken in the file
	WindowFiles
	// Keyboard represents an action taken by the keyboard
	Keyboard
	// Mouse represents an action taken by the mouse
	Mouse
	// Delete represents a delete event
	Delete
)

func (e *Event) setActionOrigin(event *acme.Event) error {
	var o Origin
	c := rune(event.C1)
	switch c {
	case 0x0:
		o = Delete
	case 'E':
		o = BodyOrTag
	case 'F':
		o = WindowFiles
	case 'K':
		o = Keyboard
	case 'M':
		o = Mouse
	default:
		return fmt.Errorf("%#U %c is not a known ActionOrigin", c, c)
	}
	e.Origin = o
	return nil
}

func (e *Event) actionOriginCode() rune {
	var o rune
	switch e.Origin {
	case BodyOrTag:
		o = 'E'
	case WindowFiles:
		o = 'F'
	case Keyboard:
		o = 'K'
	case Mouse:
		o = 'M'
	case Delete:
		o = '0'
	}
	return o
}

// Action describes what kind of action was taken
type Action int

const (
	// BodyDelete is a deletion in the window body
	BodyDelete Action = iota
	// TagDelete is a deletion in the window tag
	TagDelete
	// BodyInsert is an insertion into the window body
	BodyInsert
	// TagInsert is an insert into the window tag
	TagInsert
	// B3Body is a right click event in the window body
	B3Body
	// B3Tag is a right click event in the window tag
	B3Tag
	// B2Body is a middle click event in the window body
	B2Body
	// B2Tag is a middle click event in the window tag
	B2Tag
	// DelType represents a delete event
	DelType
)

func (e *Event) setActionType(event *acme.Event) error {
	var a Action
	c := rune(event.C2)
	switch c {
	case 'D':
		a = BodyDelete
	case 'd':
		a = TagDelete
	case 'I':
		a = BodyInsert
	case 'i':
		a = TagInsert
	case 'L':
		a = B3Body
	case 'l':
		a = B3Tag
	case 'X':
		a = B2Body
	case 'x':
		a = B2Tag
	case 0x0:
		a = DelType
	default:
		return fmt.Errorf("'%c' is not a known ActionType", c)
	}
	e.Action = a
	return nil
}

func (e *Event) actionTypeCode() rune {
	var a rune
	switch e.Action {
	case BodyDelete:
		a = 'D'
	case TagDelete:
		a = 'd'
	case BodyInsert:
		a = 'I'
	case TagInsert:
		a = 'i'
	case B3Body:
		a = 'L'
	case B3Tag:
		a = 'l'
	case B2Body:
		a = 'X'
	case B2Tag:
		a = 'x'
	case DelType:
		a = '0'
	}
	return a
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

	// TODO: determine what flag with value '3' means
)

func (e *Event) setFlag(event *acme.Event) {
	var f Flag
	if e.Action == B2Body || e.Action == B2Tag {
		switch event.Flag {
		case 1:
			f = IsBuiltin
		case 2:
			f = IsNull
		case 8:
			f = HasChordedArg
		}
	}

	if e.Action == B3Body || e.Action == B3Tag {
		switch event.Flag {
		case 1:
			f = NoReloadNeeded
		case 2:
			f = PostExpandFollows
		case 4:
			f = IsFileOrWindow
		}
	}
	e.Flag = f
}

func (e *Event) setBuiltin(event *acme.Event) error {
	text := string(event.Text)
	action := strings.ToLower(text)

	var op Builtin
	switch action {
	case "new":
		op = New
	case "zerox":
		op = Zerox
	case "get":
		op = Get
	case "put":
		op = Put
	case "del":
		op = Del
	}
	e.Builtin = op
	return nil
}

func (e *Event) setMeta(event *acme.Event) {
	e.Text = event.Text
	e.SelBegin = event.Q0
	e.SelEnd = event.Q1
	e.OrigSelBegin = event.OrigQ0
	e.OrigSelEnd = event.OrigQ1
	e.NumBytes = event.Nb
	e.NumRunes = event.Nr
	e.ChordArg = event.Arg
	e.ChordLoc = event.Loc
}

// NewEvent constructs an Event from a raw acme event
func NewEvent(event *acme.Event, id int, file string) (Event, error) {
	e := Event{
		ID:   id,
		File: file,
	}
	if err := e.setActionOrigin(event); err != nil {
		return e, err
	}
	if err := e.setActionType(event); err != nil {
		return e, err
	}
	if err := e.setBuiltin(event); err != nil {
		return e, err
	}
	e.setFlag(event)
	e.setMeta(event)
	return e, nil
}

// Log returns a raw acme log event for the Event type
func (e *Event) Log() acme.Event {
	return acme.Event{
		C1:     e.actionOriginCode(),
		C2:     e.actionTypeCode(),
		Q0:     e.SelBegin,
		Q1:     e.SelEnd,
		OrigQ0: e.OrigSelBegin,
		OrigQ1: e.OrigSelEnd,
		Flag:   int(e.Flag),
		Nb:     e.NumBytes,
		Nr:     e.NumRunes,
		Text:   e.Text,
		Arg:    e.ChordArg,
		Loc:    e.ChordLoc,
	}
}
