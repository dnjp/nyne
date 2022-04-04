package nyne

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"

	"strings"
	"unicode/utf8"

	"9fans.net/go/acme"
	"9fans.net/go/draw"
	p9client "9fans.net/go/plan9/client"
)

// Win represents the active Acme window
type Win struct {
	ID        int
	File      string
	Lastpoint int
	w         *acme.Win
}

// NewWin constructs a Win object from acme window
func NewWin() (*Win, error) {
	w, err := acme.New()
	if err != nil {
		return nil, err
	}
	return &Win{w: w}, nil
}

// OpenWin opens an acme window
func OpenWin(id int, file string) (*Win, error) {
	w, err := acme.Open(id, nil)
	if err != nil {
		return nil, err
	}
	return &Win{
		ID:   id,
		File: file,
		w:    w,
	}, nil
}

// Windows returns all open acme windows
func Windows() (map[int]*Win, error) {
	ws, err := acme.Windows()
	if err != nil {
		return nil, err
	}
	wins := make(map[int]*Win)
	for _, wi := range ws {
		w, err := acme.Open(wi.ID, nil)
		if err != nil {
			return nil, err
		}
		wins[wi.ID] = &Win{
			ID:   wi.ID,
			File: wi.Name,
			w:    w,
		}
	}
	return wins, nil
}

// FocusedWinID returns the $winid if present, otherwise it connects
// to the given addr to find the ID
//
// Derived from https://github.com/fhs/acme-lsp/blob/623cb39c2e31bddda0ad7c216c2f3c2fcfcf237f/cmd/L/main.go#L256
func FocusedWinID(addr string) (int, error) {
	winid := os.Getenv("winid")
	if winid == "" {
		conn, err := net.Dial("unix", addr)
		if err != nil {
			return 0, fmt.Errorf("$winid is empty and could not dial acmefocused: %v", err)
		}
		defer conn.Close()
		b, err := ioutil.ReadAll(conn)
		if err != nil {
			return 0, fmt.Errorf("$winid is empty and could not read acmefocused: %v", err)
		}
		winid = string(bytes.TrimSpace(b))
	}
	return strconv.Atoi(winid)
}

// FocusedWinAddr returns the address of the active window using acmefocused
func FocusedWinAddr() string {
	return filepath.Join(p9client.Namespace(), "acmefocused")
}

// EventChan opens a channel to acme events
func (w *Win) EventChan(id int, filename string, stop <-chan struct{}) (<-chan Event, <-chan error) {
	errs := make(chan error)
	events := make(chan Event)
	go func() {
		ec := w.w.EventChan()
		for {
			select {
			case e := <-ec:
				event, err := NewEvent(e, id, filename)
				if err != nil {
					errs <- err
					continue
				}
				events <- event
			case <-stop:
				return
			}
		}
	}()
	return events, errs
}

// Name sets the name for the win
func (w *Win) Name(format string, args ...interface{}) error {
	return w.w.Name(format, args...)
}

// Close closes down the window with associated files
func (w *Win) Close() {
	w.w.CloseFiles()
}

// WriteEvent writes the acme event to the log
func (w *Win) WriteEvent(e Event) error {
	raw := e.Log()
	return w.w.WriteEvent(&raw)
}

// Exec executes the given command in the window tag
func (w *Win) Exec(exec string, args ...string) error {
	if w == nil || w.w == nil {
		return fmt.Errorf("window handle lost")
	}

	tag, err := w.Tag()
	if err != nil {
		return fmt.Errorf("could not read tag: %w", err)
	}

	var before string
	parts := strings.Split(string(tag), "|")
	if len(parts) >= 2 {
		before = parts[1]
	}

	cmd := fmt.Sprintf("%s %s", exec, strings.Join(args, " "))
	if err := w.AppendTag(cmd); err != nil {
		return fmt.Errorf("could not write tag: %w", err)
	}

	rc := utf8.RuneCount(tag)
	nr := utf8.RuneCountInString(cmd)
	evt := Event{
		Origin:   Mouse,
		Type:     B2Tag,
		SelBegin: rc,
		SelEnd:   rc + nr,
		NumRunes: nr,
		Text:     []byte(cmd),
		Flag:     HasChordedArg,
	}

	log := evt.Log()
	err = w.w.WriteEvent(&log)
	if err != nil {
		fmt.Fprintf(os.Stderr, "log=%+v\n", log)
		return fmt.Errorf("could not write event: %w", err)
	}

	if err = w.ClearTag(); err != nil {
		return fmt.Errorf("could not clear tag: %w", err)
	}

	if err := w.AppendTag(before); err != nil {
		return fmt.Errorf("could not write tag: %w", err)
	}

	return nil
}

// Get is the equivalent to the Get interactive command with no
// arguments; accepts no arguments.
func (w *Win) Get() error {
	return w.write("ctl", []byte("get"))
}

// Del is the equivalent to the Del interactive command.
func (w *Win) Del() error {
	return w.write("ctl", []byte("del"))
}

// Put is the equivalent to the Put interactive command with no
// arguments; accepts no arguments.
func (w *Win) Put() error {
	return w.write("ctl", []byte("put"))
}

// Dump sets the command string to recreate the window from a
// dump file.
func (w *Win) Dump(file string) error {
	return w.write("ctl", []byte(fmt.Sprintf("dump %s", file)))
}

// Dumpdir sets the directory in which to run the command to recreate
// the window from a dump file.
func (w *Win) Dumpdir(dir string) error {
	return w.write("ctl", []byte(fmt.Sprintf("dumpdir %s", dir)))
}

// Show guarantees at least some of the selected text is visible on
// the display.
func (w *Win) Show() error {
	return w.write("ctl", []byte("show"))
}

// NoMark turns off automatic ‘marking’ of changes, so a set of
// related changes may be undone in a single Undo interactive command.
func (w *Win) NoMark() error {
	return w.write("ctl", []byte("nomark"))
}

// DisableNoMark cancels nomark, returning the window to the usual state
// wherein each modification to the body must be undone individually.
func (w *Win) DisableNoMark() error {
	return w.write("ctl", []byte("mark"))
}

// Clean marks the window clean as though it has just been written.
func (w *Win) Clean() error {
	return w.write("ctl", []byte("clean"))
}

// Dirty marks the window dirty, the opposite of clean.
func (w *Win) Dirty() error {
	return w.write("ctl", []byte("dirty"))
}

// Tag returns the tag contents
func (w *Win) Tag() ([]byte, error) {
	if w == nil || w.w == nil {
		return []byte{}, fmt.Errorf("window handle lost")
	}
	return w.w.ReadAll("tag")
}

// ClearTag removes all text in the tag after the vertical bar.
func (w *Win) ClearTag() error {
	return w.write("ctl", []byte("cleartag"))
}

// AppendTag writes to the windows tag
func (w *Win) AppendTag(text string) error {
	if w == nil || w.w == nil {
		return fmt.Errorf("window handle lost")
	}
	return w.w.Fprintf("tag", "%s", text)
}

// Body returns the window body
func (w *Win) Body() ([]byte, error) {
	if w == nil || w.w == nil {
		return []byte{}, fmt.Errorf("window handle lost")
	}
	return w.w.ReadAll("body")
}

// ClearBody clears the text from the body
func (w *Win) ClearBody() error {
	if err := w.SetAddr(","); err != nil {
		return fmt.Errorf("could not set addr: %w", err)
	}
	if err := w.SetData(nil); err != nil {
		return fmt.Errorf("could not set data: %w", err)
	}
	return nil
}

// AppendBody appends the given text to the body
func (w *Win) AppendBody(data []byte) error {
	if w == nil || w.w == nil {
		return fmt.Errorf("window handle lost")
	}
	return w.write("body", data)
}

// Char reads the character at q0
func (w *Win) Char(q0 int) (c byte, err error) {
	var dat []byte
	err = w.SetAddr("#%d;#%d", q0, q0+1)
	if err != nil {
		if err.Error() == "address out of range" {
			c = 0
			err = io.EOF
			return
		}
		return
	}
	dat, err = w.Data(q0, q0+1)
	if err != nil {
		return
	}
	if len(dat) == 0 {
		err = fmt.Errorf("no data")
		return
	}
	c = dat[0]
	return
}

// SetAddr takes an addr which may be written with any textual address
// in the format understood by button 3 but without the initial colon
func (w *Win) SetAddr(fmtstr string, args ...interface{}) error {
	addr := fmtstr
	if len(args) > 0 {
		addr = fmt.Sprintf(fmtstr, args...)
	}
	return w.w.Addr(addr)
}

// Addr returns the current address of the window
//
// Derived from https://github.com/fhs/acme-lsp/blob/623cb39c2e31bddda0ad7c216c2f3c2fcfcf237f/internal/acme/acme.go#L366
func (w *Win) Addr() (q0, q1 int, err error) {
	if w == nil || w.w == nil {
		return 0, 0, fmt.Errorf("window handle lost")
	}
	buf, err := w.w.ReadAll("addr")
	if err != nil {
		return 0, 0, err
	}
	a := strings.Fields(string(buf))
	if len(a) < 2 {
		return 0, 0, errors.New("short read from acme addr")
	}
	q0, err0 := strconv.Atoi(a[0])
	q1, err1 := strconv.Atoi(a[1])
	if err0 != nil || err1 != nil {
		return 0, 0, errors.New("invalid read from acme addr")
	}
	return q0, q1, nil
}

// CurrentAddr sets the addr to dot and reads the addr
func (w *Win) CurrentAddr() (q0, q1 int, err error) {
	_, _, err = w.Addr() // open addr file
	if err != nil {
		return 0, 0, fmt.Errorf("read addr: %v", err)
	}
	err = w.AddrFromSelection()
	if err != nil {
		return 0, 0, fmt.Errorf("setting addr=dot: %v", err)
	}
	return w.Addr()
}

// AddrFromSelection sets the addr address to that of the user’s selected
// text in the window.
func (w *Win) AddrFromSelection() error {
	return w.write("ctl", []byte("addr=dot"))
}

// SelectionFromAddr sets the user’s selected text in the window to the text
// addressed by the addr address.
func (w *Win) SelectionFromAddr() error {
	return w.write("ctl", []byte("dot=addr"))
}

// LimitSearchToAddr restricts subsequent searches to the current addr
// address.
func (w *Win) LimitSearchToAddr() error {
	return w.write("ctl", []byte("limit=addr"))
}

// SetData is used in conjunction with addr for random access to the
// contents of the body. The file offset is ignored when writing the
// data file; instead the location of the data to be read or written is
// determined by the state of the addr file. Text, which must contain only
// whole characters (no ‘partial runes’), written to data replaces the
// characters addressed by the addr file and sets the address to the null
// string at the end of the written text. A read from data returns as many
// whole characters as the read count will permit starting at the beginning
// of the addr address (the end of the address has no effect) and sets the
// address to the null string at the end of the returned characters.
func (w *Win) SetData(data []byte) error {
	return w.write("data", data)
}

// Data reads the data in the body between q0 and q1. It is assumed
// that CurrentAddr() or similar has been called to properly set the addr
// and retrieve valid q0 and q1 points.
func (w *Win) Data(q0, q1 int) ([]byte, error) {
	n := q1 - q0
	buf := make([]byte, n)
	n2, err := w.w.Read("data", buf)
	if err != nil {
		return buf, err
	}
	if n2 != n {
		return buf, fmt.Errorf("read %d bytes, expected %d", n2, n)
	}
	return buf, nil
}

// SetFont sets the font for the win
func (w *Win) SetFont(font string) error {
	return w.w.Ctl("font %s", font)
}

// Font returns the font for the current win
func (w *Win) Font() (tab int, font *draw.Font, err error) {
	return w.w.Font()
}

func (w *Win) write(file string, data []byte) error {
	if w == nil || w.w == nil {
		return fmt.Errorf("window handle lost")
	}
	_, err := w.w.Write(file, data)
	return err
}
