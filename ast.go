package main

import "strconv"

type Ast struct {
	Nodes          []Node
	TopScope       *Scope
	CurrentScope   *Scope
	semanticErrors []error
	labelCount     int
}

func (ast *Ast) L() string {
	ast.labelCount++
	return strconv.Itoa(ast.labelCount)
}

func (ast *Ast) createScope(name string) {
	scope := NewScope(name, ast.CurrentScope)
	ast.CurrentScope.Children = append(ast.CurrentScope.Children, scope)
	ast.CurrentScope = scope
}

func (ast *Ast) exitScope() {
	ast.CurrentScope = ast.CurrentScope.Outer
}

func (ast *Ast) scopeDown() {
	ast.CurrentScope = ast.CurrentScope.Children[0]
}

func (ast *Ast) scopeUp() {
	ast.CurrentScope = ast.CurrentScope.Outer
	ast.CurrentScope.Children = ast.CurrentScope.Children[1:]
}

type Node interface {
}

type Stmt interface {
	Node
	stmtNode()
}

type Expr interface {
	Node
	exprNode()
}

type Spec interface {
	Node
	specNode()
}

type Decl interface {
	Node
	declNode()
}

type (
	AssignStmt struct {
		Lhs []Expr
		Rhs []Expr
	}

	ReturnStmt struct {
		Exprs []Expr
	}

	DeclStmt struct {
		Decl Decl
	}

	IfStmt struct {
		Cond Expr
		Then *BlockStmt
		Else Stmt
	}

	ForStmt struct {
		Init Stmt
		Cond Expr
		Post Stmt
		Body *BlockStmt
	}

	BlockStmt struct {
		List []Stmt
	}

	ExprStmt struct {
		Exprs []Expr
	}
)

type (
	Lit struct {
		Kind TokenKind
		Val  string
	}

	Binary struct {
		Kind  TokenKind
		Left  Expr
		Right Expr
	}

	Ident struct {
		Name    string
		_Size   int
		_Offset int
	}

	CallFunc struct {
		FuncName string
		Args     []Expr
	}

	StarExpr struct {
		X Expr
	}

	UnaryExpr struct {
		X Expr
	}
)

type (
	GenDecl struct {
		Kind  TokenKind
		Specs []Spec
	}

	FuncDecl struct {
		FuncName   *Ident
		FuncType   *FuncType
		Body       *BlockStmt
		_FrameSize int
	}
)

type (
	ValueSpec struct {
		Type       *Type
		TypeIdent  *Ident
		Names      []*Ident
		InitValues []Expr
	}
)

type (
	FuncType struct {
		Args    []*Field
		Returns []*Field
	}
)

type Field struct {
	Names []*Ident
	Type  Expr
}

func (r *ReturnStmt) stmtNode() {}
func (a *AssignStmt) stmtNode() {}
func (d *DeclStmt) stmtNode()   {}
func (b *BlockStmt) stmtNode()  {}
func (e *ExprStmt) stmtNode()   {}
func (i *IfStmt) stmtNode()     {}
func (f *ForStmt) stmtNode()    {}

func (l *Lit) exprNode()       {}
func (b *Binary) exprNode()    {}
func (i *Ident) exprNode()     {}
func (c *CallFunc) exprNode()  {}
func (s *StarExpr) exprNode()  {}
func (u *UnaryExpr) exprNode() {}
func (t *Type) exprNode()      {}

func (v *ValueSpec) specNode() {}
func (g *GenDecl) declNode()   {}
func (g *FuncDecl) declNode()  {}
