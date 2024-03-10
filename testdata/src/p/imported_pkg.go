package p

import (
	"fmt"
	"subp"
)

var (
	subpTZero      = subp.T{}    // want `use constructor NewT for type subp.T instead of a composite literal`
	subpTZeroPtr   = &subp.T{}   // want `use constructor NewT for type subp.T instead of a composite literal`
	subpTNil       = new(subp.T) // want `nil value of type subp.T may be unsafe to use, use constructor NewT instead`
	subpTComposite = subp.T{     // want `use constructor NewT for type subp.T instead of a composite literal`
		X: 1,
	}
	subpTCompositePtr = &subp.T{ // want `use constructor NewT for type subp.T instead of a composite literal`
		X: 1,
	}
	subpTColl    = []subp.T{subp.T{X: 1}}   // want `use constructor NewT for type subp.T instead of a composite literal`
	subpTPtrColl = []*subp.T{&subp.T{X: 1}} // want `use constructor NewT for type subp.T instead of a composite literal`
	correctSubpT = subp.NewT()
)

type structWithSubpTField struct {
	i int
	t subp.T
}

var structWithSubpT = structWithSubpTField{
	i: 1,
	t: subp.T{X: 1}, // want `use constructor NewT for type subp.T instead of a composite literal`
}

type structWithSubpTPtrField struct {
	i int
	t *subp.T
}

var structWithSubpTPtr = structWithSubpTPtrField{
	i: 1,
	t: &subp.T{X: 1}, // want `use constructor NewT for type subp.T instead of a composite literal`
}

func fnWithSubpT() {
	x := subp.T{}     // want `use constructor NewT for type subp.T instead of a composite literal`
	x2 := &subp.T{}   // want `use constructor NewT for type subp.T instead of a composite literal`
	x3 := new(subp.T) // want `nil value of type subp.T may be unsafe to use, use constructor NewT instead`
	fmt.Println(x, x2, x3)
}

func retSubpT() subp.T {
	return subp.T{ // want `use constructor NewT for type subp.T instead of a composite literal`
		X: 1,
	}
}

func retSubpTPtr() *subp.T {
	return &subp.T{ // want `use constructor NewT for type subp.T instead of a composite literal`
		X: 1,
	}
}
