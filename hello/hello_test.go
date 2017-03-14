package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldAddTwoNumbers(t *testing.T) {

	x := add(1, 2)
	assert.Equal(t, 3, x)
}
