package nyne

import (
	"fmt"
	"log"
	"unicode/utf8"
)

// WinFunc retrieves the Win by its ID
type WinFunc func(int) (*Win, error)

// TabwidthFunc returns the tabwidth based on properties of the Event
type TabwidthFunc func(Event) int

// Tab constructs a tab character or the equivalent width of spaces
// depending on if expand is set
func Tab(width int, expand bool) []byte {
	if expand {
		var tab []byte
		for i := 0; i < width; i++ {
			tab = append(tab, ' ')
		}
		return tab
	}
	return []byte{0x09}
}

// Tabexpand expands tabs to spaces
func Tabexpand(condition Condition, win WinFunc, tabwidth TabwidthFunc) KeyHook {
	return KeyHook{
		Key:       '\t',
		Condition: condition,
		Handler: func(e Event) Event {
			if !condition(e) {
				return e
			}

			w, err := win(e.ID)
			if err != nil {
				log.Println(err)
				return e
			}

			tab := Tab(tabwidth(e), true)

			// select current character
			err = w.SetAddr(fmt.Sprintf("#%d;+#1", e.SelBegin))
			if err != nil {
				log.Println(err)
				w.WriteEvent(e)
			}

			// replace character with tab
			w.SetData(tab)

			// update the event to reflect the change
			rc := utf8.RuneCount(tab)
			selEnd := e.SelBegin + rc
			e.Origin = WindowFiles
			e.Type = BodyInsert
			e.SelEnd = selEnd
			e.OrigSelEnd = selEnd
			e.NumRunes = rc
			e.Text = tab

			return e
		},
	}
}