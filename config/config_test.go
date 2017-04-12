package config

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestShouldWriteAndReadProperty(t *testing.T) {
	c := NewConfig()
	c.configFileName = ".graven_test.yml"
	c.Set("foo", "bar", "baz")
	c.Write()
	c2 := NewConfig()
	c2.configFileName = ".graven_test.yml"
	c2.Read()
	assert.Equal(t, "baz", c2.Get("foo", "bar"))
}