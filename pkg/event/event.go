package event

import (
	"9fans.net/go/acme"
)

type Event struct {
	Op   AcmeOp
	File string
	log  acme.LogEvent
	Win  *Win
}
