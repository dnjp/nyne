package event

import (
	"fmt"
	"strings"

	"9fans.net/go/acme"
	"git.sr.ht/~danieljamespost/nyne/util/io"
)

type AcmeOp int

const (
	NEW   AcmeOp = iota // window creation
	ZEROX               // window creation via zerox
	GET                 // load file into window
	PUT                 // write window to the named file
	DEL                 // window deletion
	FOCUS                 // window focus
)


type Handler func(*Event)

type Hook struct {
	Op      AcmeOp
	Handler Handler
}

type Listener interface {
	Listen() error
	RegisterHook(hook Hook)
}
type Acme struct {
	hooks map[AcmeOp][]Hook
	windows map[int]string
}

func NewListener() Listener {
	return &Acme{
		hooks: make(map[AcmeOp][]Hook),
		windows: make(map[int]string),
	}
}

func (a *Acme) Listen() error {
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
		if evt.Op == NEW {
			err := a.mapWindows()
			if err != nil {
				fmt.Println(err)
			}
			a.startEventListener(event.ID)
		}

		if evt != nil {
			// a.runHooks(evt) // TODO re-enable
		}
	}
}

func (a *Acme) mapWindows() error {
	ws, err := acme.Windows()
	if err != nil {
		return err
	}
	for _, w := range ws {
		a.windows[w.ID] = w.Name
	}
	return nil
}

func (a *Acme) runHooks(event *Event) {
	hooks := a.hooks[event.Op]
	if len(hooks) == 0 {
		return
	}
	if err := event.ConnectWin(); err != nil {
		io.PrintErr(err)
		return
	}

	for _, hook := range hooks {
		fn := hook.Handler
		fn(event)
	}
	event.CloseFilesForWin()
}

func (a *Acme) RegisterHook(hook Hook) {
	hooks := a.hooks[hook.Op]
	hooks = append(hooks, hook)
	a.hooks[hook.Op] = hooks
}

func parseOp(in string) (*AcmeOp, error) {
	var op AcmeOp
	switch in {
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
	case "focus":
		op = FOCUS
	default:
		return nil, fmt.Errorf("cannot handle '%s' event", in)
	}
	return &op, nil
}

func tokenizeEvent(event acme.LogEvent) (*Event, error) {
	op, err := parseOp(event.Op)
	if err != nil {
		return nil, err
	}
	if op == nil {
		return nil, fmt.Errorf("could not find matching op")
	}
	return &Event{
		ID: event.ID,
		Op:   *op,
		File: event.Name,
		log:  event,
	}, nil
}

type ActionOrigin int
const (
	BodyOrTag ActionOrigin = iota
	WindowFiles
	Keyboard
	Mouse
)
func parseActionOrigin(event *acme.Event) (*ActionOrigin, error) {
	var origin ActionOrigin
	c := rune(event.C1)
	switch c {
	case 'E':
		origin = BodyOrTag
	case 'F':
		origin = WindowFiles
	case 'K':
		origin = Keyboard
	case 'M':
		origin = Mouse
	default:
		return nil, fmt.Errorf("%c is not a known ActionOrigin", c)
	}
	return &origin, nil
}

type ActionType int
const (
	BodyDelete ActionType = iota
	TagDelete
	BodyInsert
	TagInsert
	B3Body
	B3Tag
	B2Body
	B2Tag
)
func parseActionType(event *acme.Event) (*ActionType, error) {
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
	default:
		return nil, fmt.Errorf("'%c' is not a known ActionType", c)
	}
	return &action, nil
}

// For BodyDelete, TagDelete, BodyInsert, and TagInsert the flag is
// always zero. For messages with the 1 bit on in the flag, writing the message
// back to the event file, but with the flag, count, and text omitted,
// will cause the action to be applied to the file exactly as it
// would have been if the event file had not been open.
type Flag int
const (
	// when action is B2Body and B2Tag

	// acme built-in command
	IsBuiltin Flag = iota

	// if the text is a null string that has a non-null expansion;
	// if so, another complete message will follow describing the
	// expansion exactly as if it had been indicated explicitly (its
	// flag will always be 0)
	IsNull

	// if the command has an extra (chorded) argument; if so, two
	// more complete messages will follow reporting the argument (with
	// all numbers 0 except the character count) and where it originated,
	// in the form of a fully-qualified button 3 style address.
	HasChordedArg


	// when action is B3Body or B3Tag

	// if acme can interpret the action without loading a new file
	NoReloadNeeded

	// if a second (post-expansion) message follows, analogous to that with
	// X messages
	PostExpandFollows

	// If the text is a file or window name (perhaps with address)
	// rather than plain literal text.
	IsFileOrWindow

	// TODO: determine what flag with value '3' means
)
func parseFlag(actionType ActionType, event *acme.Event) *Flag {
	var flag Flag
	if actionType == B2Body || actionType == B2Tag {
		switch event.Flag {
		case 1:
			flag = IsBuiltin
		case 2:
			flag = IsNull
		case 8:
			flag = HasChordedArg
		}
	}

	if actionType == B3Body || actionType == B3Tag {
		switch event.Flag {
		case 1:
			flag = NoReloadNeeded
		case 2:
			flag = PostExpandFollows
		case 4:
			flag = IsFileOrWindow
		}
	}
	return &flag
}

func parseBuiltin(event *acme.Event) *AcmeOp {
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
	default:
		return nil
	}
	return &op
}


type RawEvent struct {
	Origin ActionOrigin
	Type ActionType
	Text []byte
	Builtin *AcmeOp
	Flag Flag
	File string
	WinID int
}

// If the relevant text has less than 256 characters,
// it is included in the message; otherwise it is elided, the fourth number
// is 0, and the program must read it from the data file if needed. No text
// is sent on a D or d message.
func (a *Acme) tokenizeRawEvent(event *acme.Event, id int) (*RawEvent, error) {
	actionOrigin, err := parseActionOrigin(event)
	if err != nil {
		return nil, err
	}
	actionType, err := parseActionType(event)
	if err != nil {
		return nil, err
	}
	builtin := parseBuiltin(event)

	// keep file names in sync
	if builtin != nil && *builtin == PUT {
		a.mapWindows()
	}

	flag := parseFlag(*actionType, event)
	return &RawEvent{
		Origin: *actionOrigin,
		Type: *actionType,
		Text: event.Text,
		Builtin: builtin,
		Flag: *flag,
		File: a.windows[id],
		WinID: id,
	}, nil
}

func (a *Acme) startEventListener(id int) {
	fmt.Println("starting event listener")
	w, err := acme.Open(id, nil)
	if err != nil {
		fmt.Println(err)
	}

	for e := range w.EventChan() {

		fmt.Printf("RAW: %+v\n", *e)

		// empty event received on delete
		if e.C1 == 0 && e.C2 == 0 {
			w.CloseFiles()
			break
		}

		event, err := a.tokenizeRawEvent(e, id)
		if err != nil {
			w.WriteEvent(e)
			fmt.Println(err)
			return
		}
		switch (event.Builtin) {
		default:
			fmt.Printf("TOKEN: %+v\n", *event)
			w.WriteEvent(e)
		}
	}
}