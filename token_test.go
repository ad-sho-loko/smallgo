package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTokenize_Numeric(t *testing.T) {
	tests := []struct {
		b    []byte
		want Token
	}{
		{[]byte("0"), Token{Kind: NUMBER, Val: "0"}},
		{[]byte("123456789"), Token{Kind: NUMBER, Val: "123456789"}},
		{[]byte("-123"), Token{Kind: SUB, Val: ""}},
		// {[]byte("int.Max"), Token{Kind:NUMBER, Val:"123"}},
	}

	for _, tt := range tests {
		tk := NewTokenizer(tt.b)
		got := tk.Tokenize()
		assert.Equal(t, tt.want, *got[0])
	}
}

func TestTokenize_String(t *testing.T) {
	tests := []struct {
		b    []byte
		want Token
	}{
		{[]byte("var"), Token{Kind: VAR, Val: ""}},
		{[]byte("return"), Token{Kind: RETURN, Val: ""}},
		{[]byte("x"), Token{Kind: IDENT, Val: "x"}},
		{[]byte("int64"), Token{Kind: IDENT, Val: "int64"}},
		{[]byte("func"), Token{Kind: FUNC, Val: ""}},
		{[]byte("if"), Token{Kind: IF, Val: ""}},
		{[]byte("else"), Token{Kind: ELSE, Val: ""}},
		{[]byte("for"), Token{Kind: FOR, Val: ""}},
		// {[]byte("_abc"), Token{Kind: IDENT, Val: "_abc"}},
		// {[]byte("a_b_c_"), Token{Kind: IDENT, Val: "_abc"}},
	}

	for _, tt := range tests {
		tk := NewTokenizer(tt.b)
		got := tk.Tokenize()
		assert.Equal(t, tt.want, *got[0])
	}
}

func TestTokenize_Operator(t *testing.T) {
	tests := []struct {
		b    []byte
		want Token
	}{
		{[]byte("+"), Token{Kind: ADD, Val: ""}},
		{[]byte("-"), Token{Kind: SUB, Val: ""}},
		{[]byte("*"), Token{Kind: MUL, Val: ""}},
		{[]byte("/"), Token{Kind: DIV, Val: ""}},
		{[]byte("%"), Token{Kind: REM, Val: ""}},
		{[]byte("|"), Token{Kind: OR, Val: ""}},
		{[]byte("&"), Token{Kind: AND, Val: ""}},
		{[]byte("&^"), Token{Kind: AND_NOT, Val: ""}},
		{[]byte("|="), Token{Kind: OR_ASSIGN, Val: ""}},
		{[]byte("&="), Token{Kind: AND_ASSIGN, Val: ""}},
		{[]byte("^"), Token{Kind: XOR, Val: ""}},
		{[]byte("^="), Token{Kind: XOR_ASSIGN, Val: ""}},
		{[]byte("+="), Token{Kind: ADD_ASSIGN, Val: ""}},
		{[]byte("-="), Token{Kind: SUB_ASSIGN, Val: ""}},
		{[]byte("*="), Token{Kind: MUL_ASSIGN, Val: ""}},
		{[]byte("/="), Token{Kind: DIV_ASSIGN, Val: ""}},
		{[]byte("%="), Token{Kind: REM_ASSIGN, Val: ""}},
		{[]byte("++"), Token{Kind: INC, Val: ""}},
		{[]byte("--"), Token{Kind: DEC, Val: ""}},
		{[]byte("("), Token{Kind: LPAREN, Val: ""}},
		{[]byte(")"), Token{Kind: RPAREN, Val: ""}},
		{[]byte("{"), Token{Kind: LBRACE, Val: ""}},
		{[]byte("}"), Token{Kind: RBRACE, Val: ""}},
		{[]byte("["), Token{Kind: LBRACK, Val: ""}},
		{[]byte("]"), Token{Kind: RBRACK, Val: ""}},
		{[]byte(","), Token{Kind: COMMA, Val: ""}},
		{[]byte(";"), Token{Kind: SEMICOLON, Val: ""}},
		{[]byte("=="), Token{Kind: EQL, Val: ""}},
		{[]byte("!="), Token{Kind: NEQ, Val: ""}},
		{[]byte("&&"), Token{Kind: LAND, Val: ""}},
		{[]byte("||"), Token{Kind: LOR, Val: ""}},
		{[]byte("<"), Token{Kind: LSS, Val: ""}},
		{[]byte("<="), Token{Kind: LEQ, Val: ""}},
		{[]byte("<<"), Token{Kind: SHL, Val: ""}},
		{[]byte("<<="), Token{Kind: SHL_ASSIGN, Val: ""}},
		{[]byte(">"), Token{Kind: GTR, Val: ""}},
		{[]byte(">="), Token{Kind: GEQ, Val: ""}},
		{[]byte(">>"), Token{Kind: SHR, Val: ""}},
		{[]byte(">>="), Token{Kind: SHR_ASSIGN, Val: ""}},
	}

	for _, tt := range tests {
		tk := NewTokenizer(tt.b)
		got := tk.Tokenize()
		assert.Equal(t, tt.want, *got[0])
	}
}

func TestTokenize_Special(t *testing.T) {
	tests := []struct {
		b    []byte
		want Token
	}{
		{[]byte("  0  "), Token{Kind: NUMBER, Val: "0"}},
	}

	for _, tt := range tests {
		tk := NewTokenizer(tt.b)
		got := tk.Tokenize()
		assert.Equal(t, tt.want, *got[0])
	}
}

func TestTokenize_Calc(t *testing.T) {
	tests := struct {
		b    []byte
		want []*Token
	}{
		[]byte("1+2*3"),
		[]*Token{
			{Kind: NUMBER, Val: "1"},
			{Kind: ADD, Val: ""},
			{Kind: NUMBER, Val: "2"},
			{Kind: MUL, Val: ""},
			{Kind: NUMBER, Val: "3"},
			{Kind: EOF, Val: ""},
		},
	}

	tk := NewTokenizer(tests.b)
	got := tk.Tokenize()

	for i, want := range tests.want {
		assert.Equal(t, want, got[i])
	}
}

func TestTokenize_ReturnStmt(t *testing.T) {
	tests := struct {
		b    []byte
		want []*Token
	}{
		[]byte("return 5"),
		[]*Token{
			{Kind: RETURN, Val: ""},
			{Kind: NUMBER, Val: "5"},
			{Kind: EOF, Val: ""},
		},
	}

	tk := NewTokenizer(tests.b)
	got := tk.Tokenize()

	for i, want := range tests.want {
		assert.Equal(t, want, got[i])
	}
}

func TestReadString(t *testing.T){
	test1 := []byte("\"abcde\"")
	tk := NewTokenizer(test1)
	token := tk.Tokenize()
	assert.Equal(t, token[0].Val, "abcde")
}