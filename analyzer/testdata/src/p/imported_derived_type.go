package p

import (
	"subp"
)

var (
	subpT3Nil       *subp.T3                // want `nil value of type subp.T3 may be unsafe, use constructor NewT instead`
	subpT3Zero               = subp.T3{}    // want `zero value of type subp.T3 may be unsafe, use constructor NewT instead`
	subpT3ZeroPtr            = &subp.T3{}   // want `zero value of type subp.T3 may be unsafe, use constructor NewT instead`
	subpT3New                = new(subp.T3) // want `zero value of type subp.T3 may be unsafe, use constructor NewT instead`
	subpT3Composite          = subp.T3{     // want `use constructor NewT for type subp.T3 instead of a composite literal`
		X: 1,
	}
	subpT3CompositePtr = &subp.T3{ // want `use constructor NewT for type subp.T3 instead of a composite literal`
		X: 1,
	}
	subpT3Coll            = []subp.T3{subp.T3{X: 1}}   // want `use constructor NewT for type subp.T3 instead of a composite literal`
	subpT3PtrColl         = []*subp.T3{&subp.T3{X: 1}} // want `use constructor NewT for type subp.T3 instead of a composite literal`
	correctSubpT3 subp.T3 = subp.T3(subp.NewT())
)
