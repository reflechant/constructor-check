package testpackage

// constructors: NewT
type T struct { // want T:`{NewT \d*}`
	x int
}

func NewT() T {
	return T{
		x: 42,
	}
}
