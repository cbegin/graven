package main

import (
	"github.com/cbegin/graven/hello/version"
	"github.com/fatih/color"
)

func main() {
	color.Magenta("Hello %v\n", version.Version)
}
