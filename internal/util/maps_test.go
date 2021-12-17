package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldMergeMaps(t *testing.T) {
	a := map[string]string{
		"A": "1",
		"B": "1",
	}
	b := map[string]string{
		"B": "2",
	}
	expected := map[string]string{
		"A": "1",
		"B": "2",
	}
	merged := MergeMaps(a, b)
	assert.Equal(t, expected, merged)
}
