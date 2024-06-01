package main

import (
	"context"
	"fmt"
)

func main() {
	Reachable()
}

func init() {
	s := myString{Value: "hello"}
	s.Reachable()
}

func Reachable() {
	fmt.Println("reachable")
}

func Unreachable() {
	fmt.Println("unreachable")
}

func ReachableByTest() {
	fmt.Println("reachableByTest")
}

func UnusedImportStatementRemoval() context.Context {
	return context.TODO() // Test unused import statement removal
}

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
