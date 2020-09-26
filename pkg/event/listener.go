package event

import (
	"fmt"
	"log"
	"strings"
	"unicode/utf8"

	"9fans.net/go/acme"
)

type AcmeOp int

const (
	NEW   AcmeOp = iota // window creation
	ZEROX               // window creation via zerox
	GET                 // load file into window
	PUT                 // write window to the named file
	DEL                 // window deletion
)


type Handler func(*Event)

type Hook struct {
	Op      AcmeOp
	Handler Handler
}

type Listener interface {
	Listen() err
	RegisterHook(hook Hook)
}
type Acme struct {
	hooks map[AcmeOp][]Hook
}

func NewListener() Listener {
	return &Acme{
		hooks: make(map[AcmeOp]*Handler),
	}
}

func (a *Acme) Listen() err {
	l, err := acme.Log()
	if err != nil {
		return err
	}
	for {
		event, err := l.Read()
		if err != nil {
			return err
		}

		evt, err := tokenizeEvent(event)
		if err != nil {
			return err
		}
		a.runHooks(&evt)
	}
}

func (a *Acme) runHooks(event Event) {
	hooks := a.hooks[event.Op]
	if h == nil {
		return
	}
	w, err := acme.Open(id, nil)
	if err != nil {
		log.Print(err)
		return
	}
	defer w.CloseFiles()
	event.fid = w

	for _, hook := range hooks {
		fn := h.Handler
		fn(&event)
	}
}

func (a *Acme) RegisterHook(hook Hook) {
	hooks := a.hooks[h.Op]
	h.hooks = append(hooks, hook)
}


func tokenizeEvent(event acme.LogEvent) (*Event, error) {
	var op AcmeOp
	switch event.Op {
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
	default:
		return nil, fmt.Errorf("cannot handle '%s' event", event.Op)
	}

	return &Event{
		Op:   op,
		File: event.Name,
		log:  event,
	}, nil
}

// The messages have a fixed format:
//   - a character indicating the origin or cause of the action
//     - E: writes to the body or tag file
//     - F: for actions through the windows other files
//     - K: for the keyboard
//     - M: for the mouse
//   - a character indicating the type of the action
//     - D: text deleted from the body
//     - d: text deleted from the tag
//     - I: text inserted to the body
//     - i: text inserted to the tag
//     - L: for a button 3 action in the body
//     - l: for a button 3 action in the tax
//     - X for button 2 action in the body
//     - x for button 2 action in the tag
//   - four free-format blank-terminated decimal numbers
//     1. character address of action
//     2. character address of action
//     3. flag
//       - For D, d, I, and i the flag is always zero
//       - For X and x:
//         - 1 if the text indicated is an acme built-in command
//         - 2 if the text is a null string that has a non-null expansion;
//           if so, another complete message will follow describing the
//           expansion exactly as if it had been indicated explicitly (its
//           flag will always be 0)
//         - 8 if the command has an extra (chorded) argument; if so, two
//           more complete messages will follow reporting the argument (with
//           all numbers 0 except the character count) and where it originated,
//           in the form of a fully-qualified button 3 style address.
//       - for L and l:
//         - 1 if acme can interpret the action without loading a new file
//         - 2 if a second (post-expansion) message follows, analogous to that with X messages
//         - 4 if the text is a file or window name (perhaps with address) rather than plain literal text.
//       - For messages with the 1 bit on in the flag, writing the message
//         back to the event file, but with the flag, count, and text omitted,
//         will cause the action to be applied to the file exactly as it
//         would have been if the event file had not been open.
//     4. count of the characters in the optional text which may contain newlines
//   - optional text
//   - a newline
//
//
// If the relevant text has less than 256 characters,
// it is included in the message; otherwise it is elided, the fourth number
// is 0, and the program must read it from the data file if needed. No text
// is sent on a D or d message.
//
// This spec in 9fans/acme:
//
// type Event struct {
//     // The two event characters, indicating origin and type of action
//     C1, C2 rune
//
//     // The character addresses of the action.
//     // If the original event had an empty selection (OrigQ0=OrigQ1)
//     // and was accompanied by an expansion (the 2 bit is set in Flag),
//     // then Q0 and Q1 will indicate the expansion rather than the
//     // original event.
//     Q0, Q1 int
//
//     // The Q0 and Q1 of the original event, even if it was expanded.
//     // If there was no expansion, OrigQ0=Q0 and OrigQ1=Q1.
//     OrigQ0, OrigQ1 int
//
//     // The flag bits.
//     Flag int
//
//     // The number of bytes in the optional text.
//     Nb  int
//
//     // The number of characters (UTF-8 sequences) in the optional text.
//     Nr  int
//
//     // The optional text itself, encoded in UTF-8.
//     Text []byte
//
//     // The chorded argument, if present (the 8 bit is set in the flag).
//     Arg []byte
//
//     // The chorded location, if present (the 8 bit is set in the flag).
//     Loc []byte
// }
