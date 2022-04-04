package nyne

import (
	"fmt"
	"testing"

	"9fans.net/go/acme"
)

func TestBuiltin(t *testing.T) {
	testCases := []struct {
		given    []byte
		expected Text
	}{
		{[]byte("New"), New},
		{[]byte("Zerox"), Zerox},
		{[]byte("Get"), Get},
		{[]byte("Put"), Put},
		{[]byte("Del"), Del},
		{[]byte("bla"), "bla"},
	}
	for _, tc := range testCases {
		t.Run(string(tc.given), func(t *testing.T) {
			text := NewText(tc.given)
			if text != tc.expected {
				t.Fatalf("expected text %s, got %s",
					tc.expected, text)
			}
		})
	}
}

func TestOrigin(t *testing.T) {
	testCases := []struct {
		given    rune
		expected Origin
	}{
		{0x0, Delete},
		{'E', BodyOrTag},
		{'F', WindowFiles},
		{'K', Keyboard},
		{'M', Mouse},
	}
	for _, tc := range testCases {
		t.Run(string(tc.given), func(t *testing.T) {
			origin := NewOrigin(tc.given)
			if origin != tc.expected {
				t.Fatalf("expected origin %d, got %d",
					tc.expected, origin)
			}
		})
	}
}

func TestAction(t *testing.T) {
	testCases := []struct {
		given    rune
		expected Action
	}{
		{'D', BodyDelete},
		{'d', TagDelete},
		{'I', BodyInsert},
		{'i', TagInsert},
		{'L', B3Body},
		{'l', B3Tag},
		{'X', B2Body},
		{'x', B2Tag},
		{0x0, DelType},
	}
	for _, tc := range testCases {
		t.Run(string(tc.given), func(t *testing.T) {
			action := NewAction(tc.given)
			if action != tc.expected {
				t.Fatalf("expected action '%c', got '%c'",
					tc.expected, action)
			}
		})
	}
}

func TestFlag(t *testing.T) {
	testCases := []struct {
		given    int
		action   Action
		expected Flag
	}{
		{1, B2Body | B2Tag, IsBuiltin},
		{2, B2Body | B2Tag, IsNull},
		{8, B2Body | B2Tag, HasChordedArg},
		{1, B3Body | B3Tag, NoReloadNeeded},
		{2, B3Body | B3Tag, PostExpandFollows},
		{4, B3Body | B3Tag, IsFileOrWindow},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("action=%d flag=%d", tc.action, tc.given), func(t *testing.T) {
			flag := NewFlag(tc.action, tc.given)
			if flag != tc.expected {
				t.Fatalf("expected flag %d, got %d",
					tc.expected, flag)
			}
		})
	}
}

func TestNewEvent(t *testing.T) {
	rawEvent := &acme.Event{
		C1:     'E',
		C2:     'x',
		Q0:     0,
		Q1:     8,
		OrigQ0: 0,
		OrigQ1: 0,
		Flag:   8,
		Nb:     8,
		Nr:     0,
		Text:   []byte("echo"),
		Arg:    []byte("hello"),
		Loc:    []byte{},
	}
	event, err := NewEvent(rawEvent, 1, "/tmp/testfile")
	if err != nil {
		t.Fatal(err)
	}
	if event.Origin != BodyOrTag {
		t.Fatal(event.Origin)
	}
	if event.Action != B2Tag {
		t.Fatal(event.Action)
	}
	if event.Flag != HasChordedArg {
		t.Fatal(event.Flag)
	}
}
