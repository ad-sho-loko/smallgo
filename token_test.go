package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNumeric(t *testing.T) {
	tests := []struct {
		b    []byte
		want Token
	}{
		{[]byte("0"), Token{Kind: NUMBER, Val: "0"}},
		{[]byte("123456789"), Token{Kind: NUMBER, Val: "123456789"}},
		// {[]byte("-123"), Token{Kind:NUMBER, Val:"123"}},
		// {[]byte("int.Max"), Token{Kind:NUMBER, Val:"123"}},
	}

	for _, tt := range tests {
		tk := NewTokenizer(tt.b)
		got := tk.Tokenize()
		assert.Equal(t, tt.want, *got[0])
	}
}

func TestOperator(t *testing.T) {
	tests := []struct {
		b    []byte
		want Token
	}{
		{[]byte("+"), Token{Kind: PLUS, Val: ""}},
	}

	for _, tt := range tests {
		tk := NewTokenizer(tt.b)
		got := tk.Tokenize()
		assert.Equal(t, tt.want, *got[0])
	}
}

func TestSpecial(t *testing.T) {
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
