package p

type derivedT T // want derivedT:`{NewT \d* \d*}`

var (
	dZero      = derivedT{}    // want `use constructor NewT for type p.derivedT instead of a composite literal`
	dZeroPtr   = &derivedT{}   // want `use constructor NewT for type p.derivedT instead of a composite literal`
	dNil       = new(derivedT) // want `nil value of type p.derivedT may be unsafe to use, use constructor NewT instead`
	dComposite = derivedT{     // want `use constructor NewT for type p.derivedT instead of a composite literal`
		x: 1,
		s: "abc",
	}
	dCompositePtr = &derivedT{ // want `use constructor NewT for type p.derivedT instead of a composite literal`
		x: 1,
		s: "abc",
	}
	dColl    = []derivedT{derivedT{x: 1}}   // want `use constructor NewT for type p.derivedT instead of a composite literal`
	dPtrColl = []*derivedT{&derivedT{x: 1}} // want `use constructor NewT for type p.derivedT instead of a composite literal`
)
