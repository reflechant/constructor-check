package p

import (
	"fmt"
	alias "subp"
)

var (
	aliasTNil       *alias.T       // want `nil value of type subp.T may be unsafe, use constructor NewT instead`
	aliasTZero      = alias.T{}    // want `zero value of type subp.T may be unsafe, use constructor NewT instead`
	aliasTZeroPtr   = &alias.T{}   // want `zero value of type subp.T may be unsafe, use constructor NewT instead`
	aliasTNew       = new(alias.T) // want `zero value of type subp.T may be unsafe, use constructor NewT instead`
	aliasTComposite = alias.T{     // want `use constructor NewT for type subp.T instead of a composite literal`
		X: 1,
	}
	aliasTCompositePtr = &alias.T{ // want `use constructor NewT for type subp.T instead of a composite literal`
		X: 1,
	}
	aliasTColl    = []alias.T{alias.T{X: 1}}   // want `use constructor NewT for type subp.T instead of a composite literal`
	aliasTPtrColl = []*alias.T{&alias.T{X: 1}} // want `use constructor NewT for type subp.T instead of a composite literal`
	correctAliasT = alias.NewT()
)

type structWithAliasTField struct {
	i int
	t alias.T
}

var structWithAliasT = structWithAliasTField{
	i: 1,
	t: alias.T{X: 1}, // want `use constructor NewT for type subp.T instead of a composite literal`
}

type structWithAliasTPtrField struct {
	i int
	t *alias.T
}

var structWithAliasTPtr = structWithAliasTPtrField{
	i: 1,
	t: &alias.T{X: 1}, // want `use constructor NewT for type subp.T instead of a composite literal`
}

func fnWithAliasT() {
	x := alias.T{}     // want `zero value of type subp.T may be unsafe, use constructor NewT instead`
	x2 := &alias.T{}   // want `zero value of type subp.T may be unsafe, use constructor NewT instead`
	x3 := new(alias.T) // want `zero value of type subp.T may be unsafe, use constructor NewT instead`
	fmt.Println(x, x2, x3)
}

func retAliasT() alias.T {
	return alias.T{ // want `use constructor NewT for type subp.T instead of a composite literal`
		X: 1,
	}
}

func retAliasTPtr() *alias.T {
	return &alias.T{ // want `use constructor NewT for type subp.T instead of a composite literal`
		X: 1,
	}
}
