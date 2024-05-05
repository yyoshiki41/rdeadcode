package main

import "fmt"

var _ fmt.Stringer = myString{}

type myString struct {
	Value string
}

func (s myString) String() string {
	return s.Value
}

func (s myString) Reachable() {
	return
}
