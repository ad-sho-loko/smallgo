package main

type Ast struct {
	Nodes   []Node
	Symbols map[Ident]*Symbol
}

type Symbol struct {
	Type   *Type
	Offset int
}

func (a *Ast) FrameSize() int {
	sum := 0
	for _, t := range a.Symbols {
		sum += t.Type.Size
	}
	return sum
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
		Name string
	}
)

type (
	GenDecl struct {
		Kind  TokenKind
		Specs []Spec
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

func (r *ReturnStmt) stmtNode() {}
func (a *AssignStmt) stmtNode() {}
func (d *DeclStmt) stmtNode()   {}

func (l *Lit) exprNode()    {}
func (b *Binary) exprNode() {}
func (i *Ident) exprNode()  {}

func (v *ValueSpec) specNode() {}
func (g *GenDecl) declNode()   {}
