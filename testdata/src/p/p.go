// Package p is a test package for constructor-check analyzer
package p

import (
	"fmt"

	"subp"
)

// NewT is a valid constructor for type T. Here we check if it's called
// instead of constructing values of type T manually
func NewT() *T {
	return &T{
		m: make(map[int]int),
	}
}

// T is a type whose zero values are supposedly invalid
// so a constructor NewT was created.
type T struct {
	x int
	s string
	m map[int]int
}

// TODO: check for derived type literals
type tDerived T

// TODO: check for type alias literals
type tAlias = T

// same package
var (
	t     = T{}  // want `use constructor NewT for type T instead of a composite literal`
	t2    = &T{} // want `use constructor NewT for type T instead of a composite literal`
	t3    = new(T)
	justT = T{ // want `use constructor NewT for type T instead of a composite literal`
		x: 1,
		s: "abc",
	}
	ptrToT = &T{ // want `use constructor NewT for type T instead of a composite literal`
		x: 1,
		s: "abc",
	}
	tColl    = []T{T{x: 1}}   // want `use constructor NewT for type T instead of a composite literal`
	tPtrColl = []*T{&T{x: 1}} // want `use constructor NewT for type T instead of a composite literal`

)

// imported package
var (
	u     = subp.U{}  // want `use constructor NewU for type U instead of a composite literal`
	u2    = &subp.U{} // want `use constructor NewU for type U instead of a composite literal`
	u3    = new(subp.U)
	justU = subp.U{ // want `use constructor NewU for type U instead of a composite literal`
		X: 1,
	}
	ptrToU = &subp.U{ // want `use constructor NewU for type U instead of a composite literal`
		X: 1,
	}
	uColl    = []subp.U{subp.U{X: 1}}   // want `use constructor NewU for type U instead of a composite literal`
	uPtrColl = []*subp.U{&subp.U{X: 1}} // want `use constructor NewU for type U instead of a composite literal`

)

type structWithTField struct {
	i int
	t T
}

var structWithT = structWithTField{
	i: 1,
	t: T{x: 1}, // want `use constructor NewT for type T instead of a composite literal`
}

type structWithTPtrField struct {
	i int
	t *T
}

var structWithTPtr = structWithTPtrField{
	i: 1,
	t: &T{x: 1}, // want `use constructor NewT for type T instead of a composite literal`
}

func f() {
	x := T{}   // want `use constructor NewT for type T instead of a composite literal`
	x2 := &T{} // want `use constructor NewT for type T instead of a composite literal`
	// TODO: check nil values created with new
	x3 := new(T)
	fmt.Println(x, x2, x3)
}

func retT() T {
	return T{ // want `use constructor NewT for type T instead of a composite literal`
		x: 1,
	}
}

func retPtrT() *T {
	return &T{ // want `use constructor NewT for type T instead of a composite literal`
		x: 1,
	}
}
