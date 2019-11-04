package main

import (
	"fmt"
	"strings"
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

func wrapStar(n int, inner *Type) *Type {
	if n == 0 {
		return inner
	}

	return wrapStar(n-1, NewPointer(inner))
}

func unwrapStar(n int, typeName string) (int, string) {
	if strings.HasPrefix(typeName, "*") {
		return unwrapStar(n+1, typeName[1:])
	}

	return n, typeName
}

func (s *Scope) ResolveType(typeName string) (*Type, error) {
	numOfStar := 0
	unwrapTypeName := typeName
	if strings.HasPrefix(typeName, "*") {
		numOfStar, typeName = unwrapStar(0, typeName)
	}

	typ, found := s.DeclType[typeName]

	if !found {
		if s.Outer == nil {
			return nil, fmt.Errorf(UndefinedType, typeName)
		}

		return s.Outer.ResolveType(unwrapTypeName)
	}

	if numOfStar > 0 {
		return wrapStar(numOfStar, typ), nil
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
		sum += v.Type.Size
	}

	return sum
}
