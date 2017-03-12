package main

import (
	"go.riotgames.com/pipe/bt/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"testing"
)

func TestShouldAddTwoNumbers(t *testing.T) {

	x := add(1, 2)
	assert.Equal(t, 3, x)
}
