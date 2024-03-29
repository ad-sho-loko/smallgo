package main

import "strconv"

type Ast struct {
	Nodes          []Node
	semanticErrors []error
	labelCount     int
	stringLabels   []string
	strings        []string
}

func (ast *Ast) L() string {
	ast.labelCount++
	return strconv.Itoa(ast.labelCount)
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
		_Label string
	}

	TypeName struct {
		Name string
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

	IndexExpr struct {
		X Expr
		Index Expr
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
		Type       Expr
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
func (t *TypeName) exprNode()  {}
func (c *CallFunc) exprNode()  {}
func (s *StarExpr) exprNode()  {}
func (u *UnaryExpr) exprNode() {}
func (i *IndexExpr) exprNode() {}
func (t *Type) exprNode()      {}
func (f *FuncType) exprNode()  {}

func (v *ValueSpec) specNode() {}
func (g *GenDecl) declNode()   {}
func (g *FuncDecl) declNode()  {}
