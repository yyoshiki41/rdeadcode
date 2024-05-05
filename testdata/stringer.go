package main

import "fmt"

var _ fmt.Stringer = myString{}

type myString struct {
	Value string
}

// NOTE: This function is unreachable but it is neccessary to implement the fmt.Stringer interface
func (s myString) String() string {
	return s.Value
}

func (s myString) Reachable() {
	return
}

func (s *myString) Unreachable() {
	return
}
