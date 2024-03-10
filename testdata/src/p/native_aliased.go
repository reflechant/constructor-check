package p

type aliasedT = T // want aliasedT:`{NewT \d* \d*}`

var (
	atNil       *aliasedT       // want `nil value of type p.T may be unsafe, use constructor NewT instead`
	atZero      = aliasedT{}    // want `zero value of type p.T may be unsafe, use constructor NewT instead`
	atZeroPtr   = &aliasedT{}   // want `zero value of type p.T may be unsafe, use constructor NewT instead`
	atNew       = new(aliasedT) // want `zero value of type p.T may be unsafe, use constructor NewT instead`
	atComposite = aliasedT{     // want `use constructor NewT for type p.T instead of a composite literal`
		x: 1,
		s: "abc",
	}
	atCompositePtr = &aliasedT{ // want `use constructor NewT for type p.T instead of a composite literal`
		x: 1,
		s: "abc",
	}
	atColl    = []aliasedT{aliasedT{x: 1}}   // want `use constructor NewT for type p.T instead of a composite literal`
	atPtrColl = []*aliasedT{&aliasedT{x: 1}} // want `use constructor NewT for type p.T instead of a composite literal`
)
