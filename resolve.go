package main

import (
	"fmt"
	"strconv"
)

type Resolver struct {
	ast *Ast
	TopScope       *Scope
	CurrentScope   *Scope
}

func (r *Resolver) createScope(name string) {
	scope := NewScope(name, r.CurrentScope)
	r.CurrentScope.Children = append(r.CurrentScope.Children, scope)
	r.CurrentScope = scope
}

func (r *Resolver) exitScope() {
	r.CurrentScope = r.CurrentScope.Outer
}

/*
func (r *Resolver) scopeDown() {
	r.CurrentScope = r.CurrentScope.Children[0]
}

func (r *Resolver) scopeUp() {
	r.CurrentScope = r.CurrentScope.Outer
	r.CurrentScope.Children = r.CurrentScope.Children[1:]
}
*/

func (r *Resolver) VisitNode(node Node){
	switch n := node.(type) {

	case *FuncDecl:
		r.createScope(n.FuncName.Name)

		err := r.CurrentScope.RegisterSymbol(n.FuncName.Name, NewFunc())
		if err != nil {
			r.ast.semanticErrors = append(r.ast.semanticErrors, err)
		}

		for _, arg := range n.FuncType.Args{
			t, err := r.CurrentScope.ResolveTop(arg.Type)
			if err != nil{
				r.ast.semanticErrors = append(r.ast.semanticErrors, err)
			}

			for _, ident := range arg.Names{
				err = r.CurrentScope.RegisterSymbol(ident.Name, t)
				if err != nil {
					r.ast.semanticErrors = append(r.ast.semanticErrors, err)
				}
				ident._Offset = r.CurrentScope.frameSize()
				ident._Size = t.Size
			}
		}

		for _, ret := range n.FuncType.Returns {
			ret.Type, err = r.CurrentScope.ResolveTop(ret.Type)
			if err != nil {
				r.ast.semanticErrors = append(r.ast.semanticErrors, err)
			}
		}

	case *ValueSpec:
		var err error
		n.Type, err = r.CurrentScope.ResolveTop(n.Type)
		if err != nil {
			r.ast.semanticErrors = append(r.ast.semanticErrors, err)
		}

		for _, ident := range n.Names {
			_assert(n.Type != nil, "type must not be nil here")

			t := n.Type.(*Type)
			err := r.CurrentScope.RegisterSymbol(ident.Name, t)
			if err != nil {
				r.ast.semanticErrors = append(r.ast.semanticErrors, err)
			}

			if t.Kind == Array {
				sizeLit := t.ArraySize.(*Lit)
				sizeInt, _ := strconv.Atoi(sizeLit.Val)
				ident._Size = sizeInt * t.PtrOf.(*Type).Size
				ident._Offset = r.CurrentScope.frameSize()
			} else {
				ident._Size = t.Size
				ident._Offset = r.CurrentScope.frameSize()
			}
		}

	default:
		// nop
	}
}

func (r *Resolver) LeaveNode(node Node){
	switch n := node.(type) {
	case *FuncDecl:
		// FrameSize = the size of arguments + the frame size of blockstmt
		n._FrameSize = r.CurrentScope.frameSize() + r.CurrentScope.Children[0].frameSize()
		r.exitScope()
	}
}

func (r *Resolver) VisitStmt(stmt Stmt){
	switch stmt.(type) {
	case *BlockStmt:
		r.createScope("__blockStmt")
	}
}

func (r *Resolver) LeaveStmt(stmt Stmt){
	switch stmt.(type) {
	case *BlockStmt:
		r.exitScope()
	}
}

func (r *Resolver) VisitExpr(expr Expr) {
	switch e := expr.(type) {
	case *Ident:
		sym, found := r.CurrentScope.LookUpSymbol(e.Name)
		if !found {
			r.ast.semanticErrors = append(r.ast.semanticErrors,
				fmt.Errorf("undefined variable : %s", e.Name))
			return
		}

		if sym.Type.Kind == Array {
			// e._Size = sym.Type.ArraySize
			e._Offset = sym.Offset
		} else if sym.Type.Kind == String{
			e._Offset = sym.Offset
			e._Label = sym.Type.String
		}else{
			e._Size = sym.Type.Size
			e._Offset = sym.Offset
		}

	case *IndexExpr:
		ident := e.X.(*Ident)
		indexStr := e.Index.(*Lit)
		indexNum, _ := strconv.Atoi(indexStr.Val)

		sym, found := r.CurrentScope.LookUpSymbol(ident.Name)
		if !found {
			r.ast.semanticErrors = append(r.ast.semanticErrors,
				fmt.Errorf("undefined variable : %s", ident.Name))
			return
		}

		arraySize, _ := strconv.Atoi(sym.Type.ArraySize.(*Lit).Val)

		ident._Size =  sym.Type.PtrOf.(*Type).Size
		ident._Offset = sym.Offset + (arraySize * ident._Size - ident._Size * indexNum)

	case *Lit:
		if e.Kind == STRING {
			label := ".LC" + r.ast.L()
			r.ast.stringLabels = append(r.ast.stringLabels, label + ":")
			r.ast.strings = append(r.ast.strings, e.Val)
			e.Val = label
		}
	default:
		// nop
	}
}
