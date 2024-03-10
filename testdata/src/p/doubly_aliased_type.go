package p

import (
	"subp"
	"x"
)

var (
	subpAliasTZero      = subp.AliasT{}    // want `use constructor NewX for type x.X instead of a composite literal`
	subpAliasTZeroPtr   = &subp.AliasT{}   // want `use constructor NewX for type x.X instead of a composite literal`
	subpAliasTNil       = new(subp.AliasT) // want `nil value of type x.X may be unsafe to use, use constructor NewX instead`
	subpAliasTComposite = subp.AliasT{     // want `use constructor NewX for type x.X instead of a composite literal`
		X: 1,
	}
	subpAliasTCompositePtr = &subp.AliasT{ // want `use constructor NewX for type x.X instead of a composite literal`
		X: 1,
	}
	subpAliasTColl                = []subp.AliasT{subp.AliasT{X: 1}}   // want `use constructor NewX for type x.X instead of a composite literal`
	subpAliasTPtrColl             = []*subp.AliasT{&subp.AliasT{X: 1}} // want `use constructor NewX for type x.X instead of a composite literal`
	correctSubpAliasT subp.AliasT = x.NewX()
)

var (
	subpDerivedTZero      = subp.DerivedT{}    // want `use constructor NewX for type subp.DerivedT instead of a composite literal`
	subpDerivedTZeroPtr   = &subp.DerivedT{}   // want `use constructor NewX for type subp.DerivedT instead of a composite literal`
	subpDerivedTNil       = new(subp.DerivedT) // want `nil value of type subp.DerivedT may be unsafe to use, use constructor NewX instead`
	subpDerivedTComposite = subp.DerivedT{     // want `use constructor NewX for type subp.DerivedT instead of a composite literal`
		X: 1,
	}
	subpDerivedTCompositePtr = &subp.DerivedT{ // want `use constructor NewX for type subp.DerivedT instead of a composite literal`
		X: 1,
	}
	subpDerivedTColl                  = []subp.DerivedT{subp.DerivedT{X: 1}}   // want `use constructor NewX for type subp.DerivedT instead of a composite literal`
	subpDerivedTPtrColl               = []*subp.DerivedT{&subp.DerivedT{X: 1}} // want `use constructor NewX for type subp.DerivedT instead of a composite literal`
	correctSubpDerivedT subp.DerivedT = subp.DerivedT(x.NewX())
)
