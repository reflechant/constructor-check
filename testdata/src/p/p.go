// Package p is a test package for constructor-check analyzer
package p

import (
	"fmt"

	"subp"
	alias "subp"
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
type T struct { // want T:`{NewT \d* \d*}`
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
	t     = T{}  // want `use constructor NewT for type p.T instead of a composite literal`
	t2    = &T{} // want `use constructor NewT for type p.T instead of a composite literal`
	t3    = new(T)
	justT = T{ // want `use constructor NewT for type p.T instead of a composite literal`
		x: 1,
		s: "abc",
	}
	ptrToT = &T{ // want `use constructor NewT for type p.T instead of a composite literal`
		x: 1,
		s: "abc",
	}
	tColl    = []T{T{x: 1}}   // want `use constructor NewT for type p.T instead of a composite literal`
	tPtrColl = []*T{&T{x: 1}} // want `use constructor NewT for type p.T instead of a composite literal`

)

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

func f() {
	x := T{}   // want `use constructor NewT for type p.T instead of a composite literal`
	x2 := &T{} // want `use constructor NewT for type p.T instead of a composite literal`
	// TODO: check nil values created with new
	x3 := new(T)
	fmt.Println(x, x2, x3)
}

func retT() T {
	return T{ // want `use constructor NewT for type p.T instead of a composite literal`
		x: 1,
	}
}

func retPtrT() *T {
	return &T{ // want `use constructor NewT for type p.T instead of a composite literal`
		x: 1,
	}
}

// imported package
var (
	u     = subp.T{}  // wwant `use constructor NewT for type T instead of a composite literal`
	u2    = &subp.T{} // wwant `use constructor NewT for type T instead of a composite literal`
	u3    = new(subp.T)
	justU = subp.T{ // wwant `use constructor NewT for type T instead of a composite literal`
		X: 1,
	}
	ptrToU = &subp.T{ // wwant `use constructor NewT for type T instead of a composite literal`
		X: 1,
	}
	uColl    = []subp.T{subp.T{X: 1}}   // wwant `use constructor NewT for type T instead of a composite literal`
	uPtrColl = []*subp.T{&subp.T{X: 1}} // wwant `use constructor NewT for type T instead of a composite literal`
	correctU = subp.NewT()
)

// aliased imported package
var (
	au     = alias.T{}  // wwant `use constructor NewT for type T instead of a composite literal`
	au2    = &alias.T{} // wwant `use constructor NewT for type T instead of a composite literal`
	au3    = new(alias.T)
	ajustU = alias.T{ // wwant `use constructor NewT for type T instead of a composite literal`
		X: 1,
	}
	aptrToU = &alias.T{ // wwant `use constructor NewT for type T instead of a composite literal`
		X: 1,
	}
	auColl    = []alias.T{alias.T{X: 1}}   // wwant `use constructor NewT for type T instead of a composite literal`
	auPtrColl = []*alias.T{&alias.T{X: 1}} // wwant `use constructor NewT for type T instead of a composite literal`
	acorrectU = alias.NewT()
)
