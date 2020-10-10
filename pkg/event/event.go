package event

import (
	"9fans.net/go/acme"
	"fmt"
	"strings"
)

// AcmeOp contains the acme events that are available
type AcmeOp int

const (
	// NEW represents window creation
	NEW AcmeOp = iota
	// ZEROX reprents window creation via zerox
	ZEROX
	// GET loadr file into window
	GET
	// PUT writes window to the named file
	PUT
	// DEL deletes the window
	DEL
	// FOCUS is received when the window is in focus
	FOCUS
)

// Event contains metadata for each Acme event
type Event struct {
	ID      int
	File    string
	Origin  ActionOrigin
	Type    ActionType
	Text    []byte
	Builtin AcmeOp
	Flag    Flag
	SelBegin, SelEnd int
	OrigSelBegin, OrigSelEnd int
	NumBytes int
	NumRunes int
	ChordArg []byte
	ChordLoc []byte
}

// ActionOrigin is the entity that originated the action
type ActionOrigin int

const (
	// BodyOrTag represents an action received in the body or tag
	BodyOrTag ActionOrigin = iota
	// WindowFiles represents an action taken in the file
	WindowFiles
	// Keyboard represents an action taken by the keyboard
	Keyboard
	// Mouse represents an action taken by the mouse
	Mouse
	// DelOrigin represents a delete event
	DelOrigin
)

func (e *Event) setActionOrigin(event *acme.Event) error {
	var origin ActionOrigin
	c := rune(event.C1)
	switch c {
	case '0':
		origin = DelOrigin
	case 'E':
		origin = BodyOrTag
	case 'F':
		origin = WindowFiles
	case 'K':
		origin = Keyboard
	case 'M':
		origin = Mouse
	default:
		return fmt.Errorf("%c is not a known ActionOrigin", c)
	}
	e.Origin = origin
	return nil
}

func (e *Event) getActionOriginCode() rune {
	var origin rune
	switch e.Origin {
	case BodyOrTag:
		origin = 'E'
	case WindowFiles:
		origin = 'F'
	case Keyboard:
		origin = 'K'
	case Mouse:
		origin = 'M'
	case DelOrigin:
		origin = '0'
	}
	return origin
}

// ActionType describes what kind of action was taken
type ActionType int

const (
	// BodyDelete is a deletion in the window body
	BodyDelete ActionType = iota
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
	var action ActionType
	c := rune(event.C2)
	switch c {
	case 'D':
		action = BodyDelete
	case 'd':
		action = TagDelete
	case 'I':
		action = BodyInsert
	case 'i':
		action = TagInsert
	case 'L':
		action = B3Body
	case 'l':
		action = B3Tag
	case 'X':
		action = B2Body
	case 'x':
		action = B2Tag
	case '0':
		action = DelType
	default:
		return fmt.Errorf("'%c' is not a known ActionType", c)
	}
	e.Type = action
	return nil
}

func (e *Event) getActionTypeCode() rune {
	var action rune
	switch e.Type {
	case BodyDelete:
		action = 'D'
	case TagDelete:
		action = 'd'
	case BodyInsert:
		action = 'I'
	case TagInsert:
		action = 'i'
	case B3Body:
		action = 'L'
	case B3Tag:
		action = 'l'
	case B2Body:
		action = 'X'
	case B2Tag:
		action = 'x'
	case DelType:
		action = '0'
	}
	return action
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
	var flag Flag
	if e.Type == B2Body || e.Type == B2Tag {
		switch event.Flag {
		case 1:
			flag = IsBuiltin
		case 2:
			flag = IsNull
		case 8:
			flag = HasChordedArg
		}
	}

	if e.Type == B3Body || e.Type == B3Tag {
		switch event.Flag {
		case 1:
			flag = NoReloadNeeded
		case 2:
			flag = PostExpandFollows
		case 4:
			flag = IsFileOrWindow
		}
	}
	e.Flag = flag
}

func (e *Event) setBuiltin(event *acme.Event) error {
	text := string(event.Text)
	action := strings.ToLower(text)
	
	var op AcmeOp
	switch action {
	case "new":
		op = NEW
	case "zerox":
		op = ZEROX
	case "get":
		op = GET
	case "put":
		op = PUT
	case "del":
		op = DEL
	}
	e.Builtin = op
	return nil
}

func (e *Event) getBuiltin() string {
	var op string

	switch e.Builtin {
	case NEW:
		op = "new"
	case ZEROX:
		op = "zerox"
	case GET:
		op = "get"
	case PUT:
		op = "put"
	case DEL:
		op = "del"
	}
	return op
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

// If the relevant text has less than 256 characters,
// it is included in the message; otherwise it is elided, the fourth number
// is 0, and the program must read it from the data file if needed. No text
// is sent on a D or d message.
func TokenizeEvent(event *acme.Event, id int, file string) (Event, error) {
	e := Event{		
		ID:      id,	
		File:    file,
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

func (e *Event) GetLog() acme.Event {
	return acme.Event{
		C1: e.getActionOriginCode(),
		C2: e.getActionTypeCode(),
		Q0: e.SelBegin,
		Q1: e.SelEnd,
		OrigQ0: e.OrigSelBegin,
		OrigQ1: e.OrigSelEnd ,
		Flag: int(e.Flag),
		Nb: e.NumBytes,
		Nr: e.NumRunes,
		Text: e.Text,
		Arg: e.ChordArg,
		Loc: e.ChordLoc,	
	}
}
