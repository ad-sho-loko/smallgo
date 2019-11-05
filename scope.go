package main

import (
	"fmt"
	"strconv"
)

const (
	IdentIsType      = "%s is defined as type"
	DuplicatedSymbol = "%s is already defined"
	UndefinedType    = "%s is undefined as type"
)

type Scope struct {
	Name     string
	Outer    *Scope
	Children []*Scope
	DeclType map[string]*Type
	Symbols  map[string]*Symbol
}

type Symbol struct {
	Type   *Type
	Offset int
}

func NewScope(name string, outer *Scope) *Scope {
	return &Scope{
		Name:     name,
		Outer:    outer,
		DeclType: make(map[string]*Type),
		Symbols:  make(map[string]*Symbol),
	}
}

func (s *Scope) LookUpSymbol(name string) (*Symbol, bool) {
	sym, found := s.Symbols[name]

	if found {
		return sym, true
	}

	if s.Outer == nil {
		return nil, false
	}

	return s.Outer.LookUpSymbol(name)
}

func (s *Scope) RegisterSymbol(name string, typ *Type) error {
	_, isTypeName := builtinTypes[name]
	if isTypeName {
		return fmt.Errorf(IdentIsType, name)
	}

	_, found := s.LookUpSymbol(name)
	if found {
		return fmt.Errorf(DuplicatedSymbol, name)
	}

	s.Symbols[name] = &Symbol{
		Type:   typ,
		Offset: typ.Size + s.frameSize(),
	}

	return nil
}

func (s *Scope) ResolveTop(expr Expr) (*Type, error) {
	switch e := expr.(type) {
	case *Type:
		return s.ResolveType(e)
	case *TypeName:
		return s.ResolveTypeName(e.Name)
	}

	panic("cannot resolve")
}

func (s *Scope) ResolveType(typ *Type) (*Type, error) {
	if typ.Kind == Array || typ.Kind == Ptr {
		var err error
		typ.PtrOf, err = s.ResolveTop(typ.PtrOf)
		return typ, err
	}

	panic("cannot reach here")
}

func (s *Scope) ResolveTypeName(typeName string) (*Type, error) {
	typ, found := s.DeclType[typeName]

	if !found {
		if s.Outer == nil {
			return nil, fmt.Errorf(UndefinedType, typeName)
		}

		return s.Outer.ResolveTypeName(typeName)
	}

	return typ, nil
}

func (s *Scope) frameSize() int {
	sum := 0

	var values []*Symbol
	for _, s := range s.Symbols {
		values = append(values, s)
	}

	for _, v := range values {
		if v.Type.Kind == Array {
			sizeLit := v.Type.ArraySize.(*Lit)
			sizeInt, _ := strconv.Atoi(sizeLit.Val)
			sum += sizeInt * v.Type.PtrOf.(*Type).Size
		} else {
			sum += v.Type.Size
		}
	}

	return sum
}
