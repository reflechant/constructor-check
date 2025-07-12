package p

import (
	"fmt"
	"subp"
)

var (
	subpTNil       *subp.T       // want `nil value of type subp.T may be unsafe, use constructor NewT instead`
	subpTZero      = subp.T{}    // want `zero value of type subp.T may be unsafe, use constructor NewT instead`
	subpTZeroPtr   = &subp.T{}   // want `zero value of type subp.T may be unsafe, use constructor NewT instead`
	subpTNew       = new(subp.T) // want `zero value of type subp.T may be unsafe, use constructor NewT instead`
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
	x := subp.T{}     // want `zero value of type subp.T may be unsafe, use constructor NewT instead`
	x2 := &subp.T{}   // want `zero value of type subp.T may be unsafe, use constructor NewT instead`
	x3 := new(subp.T) // want `zero value of type subp.T may be unsafe, use constructor NewT instead`
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
