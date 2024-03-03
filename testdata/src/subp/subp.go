package subp

// NewU is a valid constructor for type U. Here we check if it's called
// instead of constructing values of type U manually
func NewU() *U {
	return &U{
		M: make(map[int]int),
	}
}

// U is a type whose zero values are supposedly invalid
// so a constructor NewU was created.
type U struct {
	X int
	M map[int]int
}
