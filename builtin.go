package main

type Type struct {
	Kind TypeKind
	Size int
}

type TypeKind uint

const (
	INT TypeKind = iota + 1
	FUNCTION
)

var builtinTypes = map[string]*Type{
	"int":   NewInt(),
	"int64": NewInt(),
}

func NewInt() *Type {
	return &Type{
		Kind: INT,
		Size: 8,
	}
}

func NewFunc() *Type {
	return &Type{
		Kind: FUNCTION,
		Size: 0,
	}
}
