package p

import (
	"fmt"
	alias "subp"
)

// aliased imported package
var (
	aliasZero      = alias.T{}    // want `use constructor NewT for type subp.T instead of a composite literal`
	aliasPtr       = &alias.T{}   // want `use constructor NewT for type subp.T instead of a composite literal`
	aliasNew       = new(alias.T) // want `nil value of type subp.T may be unsafe to use, use constructor NewT instead`
	aliasComposite = alias.T{     // want `use constructor NewT for type subp.T instead of a composite literal`
		X: 1,
	}
	ptrToAliasComposite = &alias.T{ // want `use constructor NewT for type subp.T instead of a composite literal`
		X: 1,
	}
	aliasColl    = []alias.T{alias.T{X: 1}}   // want `use constructor NewT for type subp.T instead of a composite literal`
	aliasPtrColl = []*alias.T{&alias.T{X: 1}} // want `use constructor NewT for type subp.T instead of a composite literal`
	correctAlias = alias.NewT()
)

type structWithAliasField struct {
	i int
	t alias.T
}

var structWithAlias = structWithAliasField{
	i: 1,
	t: alias.T{X: 1}, // want `use constructor NewT for type subp.T instead of a composite literal`
}

type structWithAliasPtrField struct {
	i int
	t *alias.T
}

var structWithAliasPtr = structWithAliasPtrField{
	i: 1,
	t: &alias.T{X: 1}, // want `use constructor NewT for type subp.T instead of a composite literal`
}

func fAlias() {
	x := alias.T{}     // want `use constructor NewT for type subp.T instead of a composite literal`
	x2 := &alias.T{}   // want `use constructor NewT for type subp.T instead of a composite literal`
	x3 := new(alias.T) // want `nil value of type subp.T may be unsafe to use, use constructor NewT instead`
	fmt.Println(x, x2, x3)
}

func retAlias() alias.T {
	return alias.T{ // want `use constructor NewT for type subp.T instead of a composite literal`
		X: 1,
	}
}

func retAliasPtr() *alias.T {
	return &alias.T{ // want `use constructor NewT for type subp.T instead of a composite literal`
		X: 1,
	}
}
