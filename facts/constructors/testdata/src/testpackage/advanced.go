package testpackage

// multiple constructor declarations are supported
// we don't prevent from declaring deprecated functions as constructors
//
// constructors: NewT2, EmptyT2,OldConstructor
type T2 struct { // want T2:`&\[{NewT2 \d+} {EmptyT2 \d+} {OldConstructor \d+}\]`
	x int
}

func NewT2() T2 {
	return T2{
		x: 42,
	}
}

func EmptyT2() T2 {
	return T2{}
}

// Deprecated: it's old and rusty
func OldConstructor(_ int) T2 {
	return T2{}
}
