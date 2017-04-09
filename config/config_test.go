package config

import (
	"testing"
	"fmt"
)

func TestShouldFoo(t *testing.T) {
	c := NewConfig()
	c.Set("foo", []map[string]interface{}{
		{"a":1},
		{"b":map[string]int{"c":2}},
	})
	c.Read()
	fmt.Printf("\n%+v\n\n", c.Get("foo"))
	fmt.Printf("\n%+v\n\n", c.Get("bar"))
}