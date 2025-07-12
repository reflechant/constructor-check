package testpackage

// declaration blocks are not supported, only individual type declarations

// constructors: NewName, NewID
type (
	Name string
	ID   int
)

func NewName() Name {
	return ""
}

func NewID() ID {
	return 0
}
