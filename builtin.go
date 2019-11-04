package main

type Type struct {
	Kind  TypeKind
	Size  int
	PtrOf *Type
}

type TypeKind uint

const (
	Int TypeKind = iota + 1
	Byte
	Ptr
	Function
)

var builtinTypes = map[string]*Type{
	"int":   NewInt(),
	"int64": NewInt(),
	"byte":  NewByte(),
}

func NewInt() *Type {
	return &Type{
		Kind: Int,
		Size: 8,
	}
}

func NewByte() *Type {
	return &Type{
		Kind: Byte,
		Size: 1,
	}
}

func NewPointer(typ *Type) *Type {
	return &Type{
		Kind:  Ptr,
		Size:  8,
		PtrOf: typ,
	}
}

func NewFunc() *Type {
	return &Type{
		Kind: Function,
		Size: 0,
	}
}
