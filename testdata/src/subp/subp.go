package subp

// NewT is a valid constructor for type T. Here we check if it's called
// instead of constructing values of type T manually
func NewT() T {
	return T{
		M: make(map[int]int),
	}
}

// T is a type whose zero values are supposedly invalid
// so a constructor NewU was created.
type T struct {
	X int
	M map[int]int
}

type TForeignConstructor struct {
	X int
}
