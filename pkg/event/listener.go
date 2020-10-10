package event

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"unicode/utf8"
	
	"9fans.net/go/acme"
)

// Listener can listen for acme Event and Window hooks
type Listener interface {
	Listen() error
	RegisterPHook(hook EventHook)
	RegisterNHook(hook WinHook)
	SetTabexpand(w *Win, width int)
	GetEventLoopByID(id int) *FileLoop	
}

// Acme implements the Listener interface for acme events
type Acme struct {
	eventHooks map[AcmeOp][]EventHook
	winHooks   map[AcmeOp][]WinHook
	windows    map[int]string
	eventLoops map[int]*FileLoop	
	debug      bool
	mux        sync.Mutex
}

type EventLoop interface {
	GetWin() *Win
	Start() error
}

type FileLoop struct {
	ID int
	File string
	Win *Win
	debug bool
	eventHooks map[AcmeOp][]EventHook
	winHooks   map[AcmeOp][]WinHook
}

// NewListener constructs an Acme Listener
func NewListener() Listener {
	return &Acme{
		eventHooks: make(map[AcmeOp][]EventHook),
		winHooks:   make(map[AcmeOp][]WinHook),
		windows:    make(map[int]string),
		eventLoops: make(map[int]*FileLoop),
	}
}

// Listen watches the acme event log for events and executes hooks
// based on those events
func (a *Acme) Listen() error {
	if len(os.Getenv("DEBUG")) > 0 {
		a.debug = true
	}

	if a.debug {
		log.Println("opening acme log")
	}
	l, err := acme.Log()
	if err != nil {
		if a.debug {
			log.Printf("failed to read acme log: %v\n", err)
		}
		return err
	}
	for {
		if a.debug {
			log.Println("reading acme event")
		}
		event, err := l.Read()
		if err != nil {
			if a.debug {
				log.Printf("failed to read acme event: %v\n", err)
			}
			return err
		}
		// skip directory windows
		if strings.HasSuffix(event.Name, "/") {
			continue
		}
		// create listener on new window events
		if event.Op == "new" {
			
			err := a.mapWindows()
			if err != nil {
				if a.debug {
					log.Println("failed to map win IDs")
				}
				log.Println(err)
				continue
			}
			if a.isDisabled(event.ID) {
				continue
			}
			f := &FileLoop{
				ID: event.ID,
				File: a.windows[event.ID],
				debug: a.debug,
				eventHooks: a.eventHooks,
				winHooks: a.winHooks,
			}
			a.eventLoops[event.ID] = f
			go a.startEventLoop(f)		
		}
	}
}

func (a *Acme) startEventLoop(f *FileLoop) {
	log.Fatal(f.Start())
}

func (a *Acme) isDisabled(id int) bool {
	filename := a.windows[id]
	// TODO: this should be decerned in a more intelligent way
	disabledNames := []string{"/-", "Del", "xplor"}
	for _, name := range disabledNames {
		if strings.Contains(filename, name) {
			return true
		}
	}
	return false
}

func (a *Acme) mapWindows() error {
	if a.debug {
		log.Println("mapping win IDs to names")
	}
	ws, err := acme.Windows()
	if err != nil {
		return err
	}
	a.mux.Lock()
	defer a.mux.Unlock()
	a.windows = make(map[int]string)
	for _, w := range ws {
		a.windows[w.ID] = w.Name
	}
	return nil
}

func (a *Acme) GetEventLoopByID(id int) *FileLoop {
	return a.eventLoops[id]
}
func (f *FileLoop) GetWin() *Win {
	return f.Win
}
func (f *FileLoop) Start() error {
	if f.debug {
		log.Println("opening acme window")
	}
	// open window for modification
	w, err := OpenWin(f.ID, f.File)
	if err != nil {
		if f.debug {
			log.Println("failed to open acme window: %v", err)
		}
		return err
	}	
	f.Win = w

	// runs hooks for acme 'new' event
	f.runWinHooks(f.Win)

	for e := range f.Win.OpenEventChan() {
		if f.debug {
			log.Printf("RAW: %+v\n", *e)
		}

		event, err := TokenizeEvent(e, f.ID, f.File)
		if err != nil {
			return err
		}
		
		if event.Origin == DelOrigin && event.Type == DelType {
			f.Win.WriteEvent(event)
			f.Win.Close()
			return nil	
		}

		if f.debug {
			log.Printf("TOKEN: %+v\n", event)
		}

		newEvent := f.runEventHooks(event)
		if f.debug {
			log.Printf("NewEvent: %+v\n", newEvent)
		}
		f.Win.WriteEvent(newEvent)
	}
	return nil
}

func (a *Acme) ReadBodyForID(id int) ([]byte, error) {
	f := a.eventLoops[id]
	if f == nil {
		return []byte{}, fmt.Errorf("event loop not found")
	}
	return f.Win.ReadBody()
}

func (a *Acme) SetTabexpand(w *Win, width int) {
	var tab []byte
	for i := 0; i < width; i++ {
		tab = append(tab, ' ')
	}

	for e := range w.handle.EventChan() {
		if e.C1 == 0 && e.C2 == 0 {
			if a.debug {
				log.Println("received empty event: treating as del")
			}
			w.handle.WriteEvent(e)
			break
		}	
		event, err := TokenizeEvent(e, w.ID, "") // TODO fill file name
		if err != nil {
			if a.debug {
				log.Println(err)
			}
			w.handle.WriteEvent(e)
			break
		}	
	
		if event.Origin == Keyboard && event.Type == BodyInsert {
			evalKeyCmd(w, event, width)
			w.WriteEvent(event)	
		} else {
			w.WriteEvent(event)
		}
	}
}

func evalKeyCmd(w *Win, event Event, tabwidth int) {
	if len(event.Text) == 0 {
		return
	}
	r, _ := utf8.DecodeRune(event.Text)
	switch (r) {
	case '\t':
		var tab []byte
		for i := 0; i < tabwidth; i++ {
			tab = append(tab, ' ')
		}	
		err := w.SetAddr(fmt.Sprintf("#%d;+#1", event.SelBegin))
		if err != nil {
			log.Println(err)
			w.WriteEvent(event)
		}
		w.SetData(tab)
		runeCount := utf8.RuneCount(tab)
		selEnd := event.SelBegin + runeCount
		event.Origin = WindowFiles
		event.Type = BodyInsert
		event.SelEnd = selEnd
		event.OrigSelEnd = selEnd
		event.NumRunes = runeCount
		event.Text = tab
	}
}