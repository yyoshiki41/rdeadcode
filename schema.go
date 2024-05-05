package main

// This file contains the JSON schema for the deadcode output.
// https://pkg.go.dev/golang.org/x/tools@v0.20.0/cmd/deadcode#hdr-JSON_schema

type Package struct {
	Name  string     // declared name
	Path  string     // full import path
	Funcs []Function // list of dead functions within it
}

type Function struct {
	Name      string   // name (sans package qualifier)
	Position  Position // file/line/column of function declaration
	Generated bool     // function is declared in a generated .go file
}

type Edge struct {
	Initial  string   // initial entrypoint (main or init); first edge only
	Kind     string   // = static | dynamic
	Position Position // file/line/column of call site
	Callee   string   // target of the call
}

type Position struct {
	File      string // name of file
	Line, Col int    // line and byte index, both 1-based
}
