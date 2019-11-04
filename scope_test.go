package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScope_wrapStar(t *testing.T) {
	typ := NewInt()
	wraped := wrapStar(1, typ)

	assert.Equal(t, wraped.Kind, Ptr)
	assert.Equal(t, wraped.PtrOf.Kind, Int)

	typ = NewInt()
	wraped = wrapStar(3, typ)

	assert.Equal(t, Ptr, wraped.Kind)
	assert.Equal(t, Ptr, wraped.PtrOf.Kind)
	assert.Equal(t, Ptr, wraped.PtrOf.PtrOf.Kind)
	assert.Equal(t, Int, wraped.PtrOf.PtrOf.PtrOf.Kind)
}

func TestScope_unwrapStar(t *testing.T) {
	n, unwrapType := unwrapStar(0, "*int")
	assert.Equal(t, 1, n)
	assert.Equal(t, "int", unwrapType)

	n, unwrapType = unwrapStar(0, "***int")
	assert.Equal(t, 3, n)
	assert.Equal(t, "int", unwrapType)

}
