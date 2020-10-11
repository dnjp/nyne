package formatter

import (
	"fmt"
	"log"
	"unicode/utf8"

	"git.sr.ht/~danieljamespost/nyne/pkg/event"
)

// Keymap contains parameters for constructing a custom keymap
type Keymap struct {
	GetWinFn func(int) (*event.Win, error)
	GetIndentFn func(event.Event) int
}

// Tabexpand expands tabs to spaces
func (k *Keymap) Tabexpand(condition event.Condition) event.KeyCmdHook {
	return event.KeyCmdHook{
		Key: '\t',
		Condition: condition,
		Handler: func(evt event.Event) event.Event {
			if !condition(evt) {
				return evt
			}	
			win, err := k.GetWinFn(evt.ID)
			if err != nil {
				log.Println(err)
				return evt
			}
			indent := k.GetIndentFn(evt)
			var tab []byte
			for i := 0; i < indent; i++ {
				tab = append(tab, ' ')
			}
			err = win.SetAddr(fmt.Sprintf("#%d;+#1", evt.SelBegin))
			if err != nil {
				log.Println(err)
				win.WriteEvent(evt)
			}
			win.SetData(tab)
			runeCount := utf8.RuneCount(tab)
			selEnd := evt.SelBegin + runeCount
			evt.Origin = event.WindowFiles
			evt.Type = event.BodyInsert
			evt.SelEnd = selEnd
			evt.OrigSelEnd = selEnd
			evt.NumRunes = runeCount
			evt.Text = tab
			return evt
		},
	}
}
