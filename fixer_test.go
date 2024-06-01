package main

import (
	"errors"
	"testing"
)

func TestLookup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		filename, functionName string
		start, end             int
		err                    error
	}{
		{"testdata/lookup.txt", "Reachable", 7, 9, nil},
		{"testdata/lookup.txt", "myString.String", 15, 17, nil},
		{"testdata/lookup.txt", "myString.Unreachable", 19, 21, nil},
		{"testdata/lookup.txt", "NotFound", 0, 0, ErrFuncNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.filename+":"+tt.functionName, func(t *testing.T) {
			t.Parallel()
			start, end, err := lookup(tt.filename, tt.functionName)
			if start != tt.start || end != tt.end || !errors.Is(err, tt.err) {
				t.Errorf("lookup(%q, %q) = (%d, %d, %v); want (%d, %d, %v)",
					tt.filename, tt.functionName, start, end, err, tt.start, tt.end, tt.err)
			}
		})
	}
}
