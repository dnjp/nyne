package formatter

import (
	"fmt"
	"log"
	"unicode/utf8"

	"git.sr.ht/~danieljamespost/nyne/pkg/event"
)

// Tabexpand expands tabs to spaces
func (n *NFmt) Tabexpand() event.KeyCmdHook {
	return event.KeyCmdHook{
		Key: '\t',
		Condition: func(evt event.Event) bool {
			op, _ := n.getOp(evt.File)
			if op == nil {
				return false
			}
			return op.Fmt.tabexpand
		},
		Handler: func(evt event.Event) event.Event {	
			op, _ := n.getOp(evt.File)
			if op == nil {
				return evt
			}
			l := n.listener.GetEventLoopByID(evt.ID)
			if l == nil {
				log.Println("could not find event loop")
				return evt
			}
			win := l.GetWin()
			var tab []byte
			for i := 0; i < op.Fmt.indent; i++ {
				tab = append(tab, ' ')
			}
			err := win.SetAddr(fmt.Sprintf("#%d;+#1", evt.SelBegin))
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
