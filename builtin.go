package main

type Type struct {
	Kind TypeKind
	Size int
	PtrOf     Expr
	ArraySize Expr
	String    string
}

type TypeKind uint

const (
	Int TypeKind = iota + 1
	Byte
	Ptr
	Array
	String
	Function
)

var builtinTypes = map[string]*Type{
	"int":   NewInt(),
	"int64": NewInt(),
	"byte":  NewByte(),
	"string":  NewString(),
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

func NewPointer(typ Expr) *Type {
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

func NewString() *Type{
	return &Type{
		Kind: String,
		Size: 8,
	}
}