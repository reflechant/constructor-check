package p

import (
	"fmt"
	"subp"
)

// imported package
var (
	u     = subp.T{}    // want `use constructor NewT for type subp.T instead of a composite literal`
	u2    = &subp.T{}   // want `use constructor NewT for type subp.T instead of a composite literal`
	u3    = new(subp.T) // want `nil value of type subp.T may be unsafe to use, use constructor NewT instead`
	justU = subp.T{     // want `use constructor NewT for type subp.T instead of a composite literal`
		X: 1,
	}
	ptrToU = &subp.T{ // want `use constructor NewT for type subp.T instead of a composite literal`
		X: 1,
	}
	uColl    = []subp.T{subp.T{X: 1}}   // want `use constructor NewT for type subp.T instead of a composite literal`
	uPtrColl = []*subp.T{&subp.T{X: 1}} // want `use constructor NewT for type subp.T instead of a composite literal`
	correctU = subp.NewT()
)

type structWithUField struct {
	i int
	t subp.T
}

var structWithU = structWithUField{
	i: 1,
	t: subp.T{X: 1}, // want `use constructor NewT for type subp.T instead of a composite literal`
}

type structWithUPtrField struct {
	i int
	t *subp.T
}

var structWithUPtr = structWithUPtrField{
	i: 1,
	t: &subp.T{X: 1}, // want `use constructor NewT for type subp.T instead of a composite literal`
}

func fu() {
	x := subp.T{}     // want `use constructor NewT for type subp.T instead of a composite literal`
	x2 := &subp.T{}   // want `use constructor NewT for type subp.T instead of a composite literal`
	x3 := new(subp.T) // want `nil value of type subp.T may be unsafe to use, use constructor NewT instead`
	fmt.Println(x, x2, x3)
}

func retU() subp.T {
	return subp.T{ // want `use constructor NewT for type subp.T instead of a composite literal`
		X: 1,
	}
}

func retPtrU() *subp.T {
	return &subp.T{ // want `use constructor NewT for type subp.T instead of a composite literal`
		X: 1,
	}
}
