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
		{[]byte("-123"), Token{Kind:SUB, Val:""}},
		// {[]byte("int.Max"), Token{Kind:NUMBER, Val:"123"}},
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
		{[]byte("("), Token{Kind: LPAREN, Val: ""}},
		{[]byte(")"), Token{Kind: RPAREN, Val: ""}},

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
