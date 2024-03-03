// Package p is a test package for constructor-check analyzer
package p

import "fmt"

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

var (
	// TODO: check for zero values
	t = T{}
	// TODO: check for nil values
	t2         = &T{}
	t3         = new(T)
	tComposite = T{
		x: 1,
		s: "abc",
	}
	tColl = []T{T{x: 1}}
	n     = nested{
		i: 1,
		t: T{x: 1},
	}
)

type nested struct {
	i int
	t T
}

func f() {
	x := T{}
	x2 := &T{}
	x3 := new(T)
	fmt.Println(x, x2, x3)
}

func retT() T {
	return T{ // want `use constructor NewT for type T instead of a composite literal`
		x: 1,
	}
}
