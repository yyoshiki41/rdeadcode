//go:build ignore

package main

import "fmt"

func Reachable() {
	fmt.Println("reachable")
}

type myString struct {
	Value string
}

func (s myString) String() string {
	return s.Value
}

func (s *myString) Unreachable() {
	return
}
