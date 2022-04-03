package nyne

import (
	"bytes"
	"errors"
	"fmt"
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
func NewWin(w *acme.Win) *Win {
	return &Win{w: w}
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

// FindFocusedWinID finds the active window ID using acmefocused
func FindFocusedWinID() (int, error) {
	return FocusedWinID(filepath.Join(p9client.Namespace(), "acmefocused"))
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

// OpenEventChan opens a channel to raw acme events
func (w *Win) OpenEventChan() <-chan *acme.Event {
	return w.w.EventChan()
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

// ExecInTag executes the given command in the window tag
func (w *Win) ExecInTag(exec string, args ...string) error {
	if w == nil || w.w == nil {
		return fmt.Errorf("window handle lost")
	}

	tag, err := w.ReadTag()
	if err != nil {
		return err
	}
	offset := utf8.RuneCount(tag)

	cmd := fmt.Sprintf("%s %s", exec, strings.Join(args, " "))
	if err := w.WriteToTag(cmd); err != nil {
		return err
	}

	evt := Event{
		Origin:   Mouse,
		Type:     B2Tag,
		SelBegin: offset,
		SelEnd:   offset + utf8.RuneCountInString(cmd),
	}

	log := evt.Log()

	err = w.w.WriteEvent(&log)
	if err != nil {
		return err
	}

	return w.ClearTagText()
}

// ExecGet is the equivalent to the Get interactive command with no
// arguments; accepts no arguments.
func (w *Win) ExecGet() error {
	return w.write("ctl", []byte("get"))
}

// ExecDelete is the equivalent to the Del interactive command.
func (w *Win) ExecDelete() error {
	return w.write("ctl", []byte("del"))
}

// ExecDumpToFile sets the command string to recreate the window from a
// dump file.
func (w *Win) ExecDumpToFile(file string) error {
	return w.write("ctl", []byte(fmt.Sprintf("dump %s", file)))
}

// ExecPut is the equivalent to the Put interactive command with no
// arguments; accepts no arguments.
func (w *Win) ExecPut() error {
	return w.write("ctl", []byte("put"))
}

// ExecShow guarantees at least some of the selected text is visible on
// the display.
func (w *Win) ExecShow() error {
	return w.write("ctl", []byte("show"))
}

// SetAddrToSelText sets the addr address to that of the user’s selected
// text in the window.
func (w *Win) SetAddrToSelText() error {
	return w.write("ctl", []byte("addr=dot"))
}

// MarkWinClean marks the window clean as though it has just been written.
func (w *Win) MarkWinClean() error {
	return w.write("ctl", []byte("clean"))
}

// MarkWinDirty marks the window dirty, the opposite of clean.
func (w *Win) MarkWinDirty() error {
	return w.write("ctl", []byte("dirty"))
}

// ClearTagText removes all text in the tag after the vertical bar.
func (w *Win) ClearTagText() error {
	return w.write("ctl", []byte("cleartag"))
}

// ReadTag returns the tag contents
func (w *Win) ReadTag() ([]byte, error) {
	if w == nil || w.w == nil {
		return []byte{}, fmt.Errorf("window handle lost")
	}
	return w.w.ReadAll("tag")
}

// ReadBody returns the window body
func (w *Win) ReadBody() ([]byte, error) {
	if w == nil || w.w == nil {
		return []byte{}, fmt.Errorf("window handle lost")
	}
	return w.w.ReadAll("body")
}

// ReadAddr returns the current address of the window
//
// Derived from https://github.com/fhs/acme-lsp/blob/623cb39c2e31bddda0ad7c216c2f3c2fcfcf237f/internal/acme/acme.go#L366
func (w *Win) ReadAddr() (q0, q1 int, err error) {
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
	_, _, err = w.ReadAddr() // open addr file
	if err != nil {
		return 0, 0, fmt.Errorf("read addr: %v", err)
	}
	err = w.SetAddrToSelText()
	if err != nil {
		return 0, 0, fmt.Errorf("setting addr=dot: %v", err)
	}
	return w.ReadAddr()
}

// SetTextToAddr sets the user’s selected text in the window to the text
// addressed by the addr address.
func (w *Win) SetTextToAddr() error {
	return w.write("ctl", []byte("dot=addr"))
}

// SetDumpDir sets the directory in which to run the command to recreate
// the window from a dump file.
func (w *Win) SetDumpDir(dir string) error {
	return w.write("ctl", []byte(fmt.Sprintf("dumpdir %s", dir)))
}

// RestrictSearchToAddr restricts subsequent searches to the current addr
// address.
func (w *Win) RestrictSearchToAddr() error {
	return w.write("ctl", []byte("limit=addr"))
}

// EnableNoMark turns off automatic ‘marking’ of changes, so a set of
// related changes may be undone in a single Undo interactive command.
func (w *Win) EnableNoMark() error {
	return w.write("ctl", []byte("nomark"))
}

// DisableNoMark cancels nomark, returning the window to the usual state
// wherein each modification to the body must be undone individually.
func (w *Win) DisableNoMark() error {
	return w.write("ctl", []byte("mark"))
}

// SetWinName sets the name of the window to name.
func (w *Win) SetWinName(name string) error {
	return w.write("ctl", []byte(fmt.Sprintf("name %s", name)))
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

// Font returns the font for the current win
func (w *Win) Font() (tab int, font *draw.Font, err error) {
	return w.w.Font()
}

// SetFont sets the font for the win
func (w *Win) SetFont(font string) error {
	return w.w.Ctl("font %s", font)
}

// ReadData reads the data in the body between q0 and q1. It is assumed
// that CurrentAddr() or similar has been called to properly set the addr
// and retrieve valid q0 and q1 points.
func (w *Win) ReadData(q0, q1 int) ([]byte, error) {
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

// WriteToTag writes to the windows tag
func (w *Win) WriteToTag(text string) error {
	if w == nil || w.w == nil {
		return fmt.Errorf("window handle lost")
	}
	return w.w.Fprintf("tag", "%s", text)
}

// WriteMenu writes the specified menu options to the Acme buffer
func (w *Win) WriteMenu(menu []string) error {
	if w == nil {
		return fmt.Errorf("state has drifted: *Win is nil")
	}
	for _, opt := range menu {
		if err := w.WriteToTag(opt); err != nil {
			return err
		}
	}
	return nil
}

func (w *Win) write(file string, data []byte) error {
	if w == nil || w.w == nil {
		return fmt.Errorf("window handle lost")
	}
	_, err := w.w.Write(file, data)
	return err
}
