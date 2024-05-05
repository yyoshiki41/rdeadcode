package main

import (
	"context"
	"os"
	"testing"

	"rsc.io/script"
)

func TestGolden(t *testing.T) {
	ctx := context.Background()
	state, err := script.NewState(ctx, "./", nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = script.Cp().Run(state, "testdata/main.go", "testdata/main.go.backup")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_, err = script.Mv().Run(state, "testdata/main.go.backup", "testdata/main.go")
		if err != nil {
			t.Fatal(err)
		}
	}()

	f, err := os.Open("testdata/golden.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	os.Stdin = f
	main()

	got, err := os.Stat("testdata/main.go")
	expected, err := os.Stat("testdata/main.go.golden")
	if !os.SameFile(got, expected) {
		t.Errorf("unexpected file content")
	}
}
