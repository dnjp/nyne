package main

import (
	"testing"
)

func TestMovedown(t *testing.T) {
	body := []byte(`	printf("hello world\n");
printf("hello world\n");
}
`)
	testCases := []struct {
		currentQ0 int
		startQ0   int
		endQ0     int
		tabwidth  int
	}{
		//	printf("hello world\n");
		// ^
		{
			currentQ0: 88,
			startQ0:   88,
			endQ0:     114,
			tabwidth:  8,
		},
		//	printf("hello world\n");
		// 	^
		{
			currentQ0: 89,
			startQ0:   88,
			endQ0:     122,
			tabwidth:  8,
		},
		//	printf("hello world\n");
		// 	                       ^
		{
			currentQ0: 113,
			startQ0:   88,
			endQ0:     138,
			tabwidth:  8,
		},
	}
	for _, tc := range testCases {
		endQ0c := body[tc.endQ0-tc.startQ0]
		nQ0 := down(body, tc.tabwidth, tc.startQ0, tc.currentQ0)
		nQ0c := body[nQ0-tc.currentQ0]
		if nQ0 != tc.endQ0 {
			t.Fatalf("expected nq0=%d(%c), got nq0=%d(%c)",
				tc.endQ0, endQ0c, nQ0, nQ0c)
		}
	}
}
