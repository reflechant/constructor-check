package subp

import "x"

// NewT is a valid constructor for type T. Here we check if it's called
// instead of constructing values of type T manually
func NewT() T {
	return T{
		M: make(map[int]int),
	}
}

type T2 = T

type T3 T

// T is a type whose zero values are supposedly invalid
// so a constructor NewU was created.
type T struct {
	X int
	M map[int]int
}

type TForeignConstructor struct {
	X int
}

type AliasT = x.X

type DerivedT x.X
