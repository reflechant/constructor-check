package p

import "subp"

// this is not considered a constructor as it returns a type from another package
func NewTForeignConstructor() subp.TForeignConstructor {
	return subp.TForeignConstructor{
		X: 42,
	}
}

type T3 struct {
	x int
}

// this is not considered a constructor as it returns multiple values
func NewMultipleT3() (T3, T3) {
	return T3{}, T3{}
}

type T4 struct {
	x int
}

// this is not considered a constructor as it doesn't follow `NewT` naming pattern
func NewZeroT4() T4 {
	return T4{}
}
