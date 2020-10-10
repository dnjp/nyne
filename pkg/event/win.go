package event

import (
	"fmt"
	// "log"
	"strings"
	"unicode/utf8"

	"9fans.net/go/acme"
)

// Win represents the active Acme window
type Win struct {
	ID     int
	File   string
	handle *acme.Win
}

// ExecInTag executes the given command in the window tag
func (w *Win) ExecInTag(exec string, args ...string) error {
	if w == nil || w.handle == nil {
		return fmt.Errorf("window handle lost")
	}
	cmd := fmt.Sprintf("%s %s", exec, strings.Join(args, " "))

	tag, err := w.ReadTag()
	if err != nil {
		panic(err)
		// log.Print(err)
	}
	offset := utf8.RuneCount(tag)
	cmdlen := utf8.RuneCountInString(cmd)
	if err := w.WriteToTag(cmd); err != nil {
		panic(err)
		// log.Print(err)
	}
	evt := new(acme.Event)
	evt.C1 = 'M'
	evt.C2 = 'x'
	evt.Q0 = offset
	evt.Q1 = offset + cmdlen
	return w.handle.WriteEvent(evt)
}

// ExecGet is the equivalent to the Get interactive command with no arguments; accepts no arguments.
func (w *Win) ExecGet() error {
	return w.write("ctl", []byte("get"))
}

// ExecDelete is the equivalent to the Del interactive command.
func (w *Win) ExecDelete() error {
	return w.write("ctl", []byte("del"))
}

// ExecDumpToFile sets the command string to recreate the window from a dump file.
func (w *Win) ExecDumpToFile(file string) error {
	return w.write("ctl", []byte(fmt.Sprintf("dump %s", file)))
}

// ExecPut is the equivalent to the Put interactive command with no arguments; accepts no arguments.
func (w *Win) ExecPut() error {
	return w.write("ctl", []byte("put"))
}

// ExecShow guarantees at least some of the selected text is visible on the display.
func (w *Win) ExecShow() error {
	return w.write("ctl", []byte("show"))
}

// SetAddrToSelText sets the addr address to that of the user’s selected text in the window.
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
	if w == nil || w.handle == nil {
		return []byte{}, fmt.Errorf("window handle lost")
	}
	return w.handle.ReadAll("tag")
}

// ReadBody returns the window body
func (w *Win) ReadBody() ([]byte, error) {
	if w == nil || w.handle == nil {
		return []byte{}, fmt.Errorf("window handle lost")
	}
	return w.handle.ReadAll("body")
}

// ReadAddr returns the current address of the window
func (w *Win) ReadAddr() ([]byte, error) {
	if w == nil || w.handle == nil {
		return []byte{}, fmt.Errorf("window handle lost")
	}
	return w.handle.ReadAll("addr")
}

// SetTextToAddr sets the user’s selected text in the window to the text addressed by the addr address.
func (w *Win) SetTextToAddr() error {
	return w.write("ctl", []byte("dot=addr"))
}

// SetDumpDir sets the directory in which to run the command to recreate the window from a dump file.
func (w *Win) SetDumpDir(dir string) error {
	return w.write("ctl", []byte(fmt.Sprintf("dumpdir %s", dir)))
}

// RestrictSearchToAddr restricts subsequent searches to the current addr address.
func (w *Win) RestrictSearchToAddr() error {
	return w.write("ctl", []byte("limit=addr"))
}

// EnableNoMark turns off automatic ‘marking’ of changes, so a set of related changes may be undone in a single Undo interactive command.
func (w *Win) EnableNoMark() error {
	return w.write("ctl", []byte("nomark"))
}

// DisableNoMark cancels nomark, returning the window to the usual state wherein each modification to the body must be undone individually.
func (w *Win) DisableNoMark() error {
	return w.write("ctl", []byte("mark"))
}

// SetWinName sets the name of the window to name.
func (w *Win) SetWinName(name string) error {
	return w.write("ctl", []byte(fmt.Sprintf("name %s", name)))
}

// SetAddr takes an addr which may be written with any textual address
// in the format understood by button 3 but without the initial colon
func (w *Win) SetAddr(addr string) error {
	return w.write("addr", []byte(addr))
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

// WriteToTag writes to the windows tag
func (w *Win) WriteToTag(text string) error {
	if w == nil || w.handle == nil {
		return fmt.Errorf("window handle lost")
	}
	return w.handle.Fprintf("tag", "%s", text)
}

func (w *Win) write(file string, data []byte) error {
	if w == nil || w.handle == nil {
		return fmt.Errorf("window handle lost")
	}
	_, err := w.handle.Write(file, data)
	return err
}
