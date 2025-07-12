package p

import (
	"subp"
)

var (
	subpT2Nil       *subp.T2                // want `nil value of type subp.T may be unsafe, use constructor NewT instead`
	subpT2Zero               = subp.T2{}    // want `zero value of type subp.T may be unsafe, use constructor NewT instead`
	subpT2ZeroPtr            = &subp.T2{}   // want `zero value of type subp.T may be unsafe, use constructor NewT instead`
	subpT2New                = new(subp.T2) // want `zero value of type subp.T may be unsafe, use constructor NewT instead`
	subpT2Composite          = subp.T2{     // want `use constructor NewT for type subp.T instead of a composite literal`
		X: 1,
	}
	subpT2CompositePtr = &subp.T2{ // want `use constructor NewT for type subp.T instead of a composite literal`
		X: 1,
	}
	subpT2Coll            = []subp.T2{subp.T2{X: 1}}   // want `use constructor NewT for type subp.T instead of a composite literal`
	subpT2PtrColl         = []*subp.T2{&subp.T2{X: 1}} // want `use constructor NewT for type subp.T instead of a composite literal`
	correctSubpT2 subp.T2 = subp.NewT()
)
