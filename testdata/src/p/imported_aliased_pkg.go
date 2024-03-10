package p

import (
	"fmt"
	alias "subp"
)

var (
	aliasTZero      = alias.T{}    // want `use constructor NewT for type subp.T instead of a composite literal`
	aliasTZeroPtr   = &alias.T{}   // want `use constructor NewT for type subp.T instead of a composite literal`
	aliasTNil       = new(alias.T) // want `nil value of type subp.T may be unsafe to use, use constructor NewT instead`
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
	x := alias.T{}     // want `use constructor NewT for type subp.T instead of a composite literal`
	x2 := &alias.T{}   // want `use constructor NewT for type subp.T instead of a composite literal`
	x3 := new(alias.T) // want `nil value of type subp.T may be unsafe to use, use constructor NewT instead`
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
