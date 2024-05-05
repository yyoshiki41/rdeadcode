package main

import (
	"testing"
)

func TestLookup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		filename, functionName string
		start, end             int
		err                    error
	}{
		{"testdata/reachable.go", "Reachable", 5, 7, nil},
		{"testdata/stringer.go", "myString.String", 12, 14, nil},
		{"testdata/stringer.go", "myString.Unreachable", 20, 22, nil},
	}

	for _, tt := range tests {
		t.Run(tt.filename+":"+tt.functionName, func(t *testing.T) {
			t.Parallel()
			start, end, err := lookup(tt.filename, tt.functionName)
			if start != tt.start || end != tt.end || err != tt.err {
				t.Errorf("lookup(%q, %q) = (%d, %d, %v); want (%d, %d, %v)",
					tt.filename, tt.functionName, start, end, err, tt.start, tt.end, tt.err)
			}
		})
	}
}
