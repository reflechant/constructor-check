// Package p is a test package for constructor-check analyzer
package p

import (
	"bytes"
	"fmt"
)

var buf = bytes.Buffer{} // standard library is excluded from analysis

// T is a type whose zero values are supposedly invalid
// so a constructor NewT was created.
type T struct { // want T:`{NewT \d* \d*}`
	x int
	s string
	m map[int]int
}

var (
	tZero      = T{}    // want `use constructor NewT for type p.T instead of a composite literal`
	tZeroPtr   = &T{}   // want `use constructor NewT for type p.T instead of a composite literal`
	tNil       = new(T) // want `nil value of type p.T may be unsafe to use, use constructor NewT instead`
	tComposite = T{     // want `use constructor NewT for type p.T instead of a composite literal`
		x: 1,
		s: "abc",
	}
	tCompositePtr = &T{ // want `use constructor NewT for type p.T instead of a composite literal`
		x: 1,
		s: "abc",
	}
	tColl    = []T{T{x: 1}}   // want `use constructor NewT for type p.T instead of a composite literal`
	tPtrColl = []*T{&T{x: 1}} // want `use constructor NewT for type p.T instead of a composite literal`

)

// NewT is a valid constructor for type T. Here we check if it's called
// instead of constructing values of type T manually
func NewT() *T {
	return &T{
		m: make(map[int]int),
	}
}

type structWithTField struct {
	i int
	t T
}

var structWithT = structWithTField{
	i: 1,
	t: T{x: 1}, // want `use constructor NewT for type p.T instead of a composite literal`
}

type structWithTPtrField struct {
	i int
	t *T
}

var structWithTPtr = structWithTPtrField{
	i: 1,
	t: &T{x: 1}, // want `use constructor NewT for type p.T instead of a composite literal`
}

func fnWithT() {
	x := T{}     // want `use constructor NewT for type p.T instead of a composite literal`
	x2 := &T{}   // want `use constructor NewT for type p.T instead of a composite literal`
	x3 := new(T) // want `nil value of type p.T may be unsafe to use, use constructor NewT instead`
	fmt.Println(x, x2, x3)
}

func retT() T {
	return T{ // want `use constructor NewT for type p.T instead of a composite literal`
		x: 1,
	}
}

func retTPtr() *T {
	return &T{ // want `use constructor NewT for type p.T instead of a composite literal`
		x: 1,
	}
}

type T2 struct { // want T2:`{NewT2 \d* \d*}`
	x int
}

func NewT2() *T2 {
	// new(T) inside T's constructor is permitted
	return new(T2)
}
