package main

import "fmt"

const (
	IdentIsType   = "sema.go : %s is defined as type"
	UndefinedType = "sema.go : %s is undefined as type"
)

func assertIdentIsNotTypeName(ident string) error {
	_, found := builtinTypes[ident]
	if found {
		return fmt.Errorf(IdentIsType, ident)
	}

	return nil
}

func resolveType(typeName string) (*Type, error) {
	t, found := builtinTypes[typeName]
	if !found {
		return nil, fmt.Errorf(UndefinedType, typeName)
	}
	return t, nil
}

//noinspection ALL
func walkNode(ast *Ast, n Node) error {
	switch typ := n.(type) {
	case *FuncDecl:
		if typ.ReturnTypeIdent != nil && typ.ReturnType == nil {
			t, err := resolveType(typ.ReturnTypeIdent.Name)
			if err != nil {
				return err
			}
			typ.ReturnType = t
		}
		walkNode(ast, typ.FuncName)
		walkNode(ast, typ.ReturnType)

		for _, b := range typ.Body.List {
			walkNode(ast, b)
		}

	case *DeclStmt:
		walkNode(ast, typ.Decl)

	case *GenDecl:
		for _, spec := range typ.Specs {
			walkNode(ast, spec)
		}

	case *ValueSpec:
		if typ.Type == nil {
			t, err := resolveType(typ.TypeIdent.Name)
			if err != nil {
				return err
			}
			typ.Type = t
		}

		for _, ident := range typ.Names {
			ast.Symbols[*ident] = &Symbol{
				Type:   typ.Type,
				Offset: typ.Type.Size + ast.FrameSize(),
			}
		}

		for _, ident := range typ.Names {
			err := assertIdentIsNotTypeName(ident.Name)
			if err != nil {
				return err
			}

			walkNode(ast, ident)
		}

		walkNode(ast, typ.InitValues)

	case *Ident:
	default:
	}

	return nil
}

func WalkAst(ast *Ast) error {
	for _, n := range ast.Nodes {
		err := walkNode(ast, n)
		if err != nil {
			return err
		}
	}

	return nil
}
