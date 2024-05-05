package main

import (
	"context"
	"os"
	"testing"

	"rsc.io/script"
)

func TestMainGolden(t *testing.T) {
	ctx := context.Background()
	state, err := script.NewState(ctx, "./", nil)
	if err != nil {
		t.Fatal(err)
	}
	mainFile := "testdata/main.go"
	if _, err := script.Cp().Run(state, mainFile, mainFile+".backup"); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if _, err := script.Mv().Run(state, mainFile+".backup", mainFile); err != nil {
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

	checkContent(t, mainFile, mainFile+".golden")
}

func checkContent(t *testing.T, f1, f2 string) {
	t.Helper()

	got, err := os.ReadFile(f1)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(f2)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(expected) {
		t.Errorf("got %s, expected %s", got, expected)
	}
}
