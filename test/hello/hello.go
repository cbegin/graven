package main

import (
	"github.com/cbegin/graven/test/hello/version"
	"github.com/fatih/color"
)

func main() {
	color.Magenta("Hello %v\n", version.Version)
}

func add(a, b int) int {
	return a + b
}
