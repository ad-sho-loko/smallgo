package main

type Node interface {
}

type Expr interface {
	Node
	exprNode()
}

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
)

func (l *Lit) exprNode()    {}
func (o *Binary) exprNode() {}
