package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCommon_Roundup(t *testing.T) {
	assert.Equal(t, 8, roundup(1, 8))
	assert.Equal(t, 16, roundup(1, 16))
}
