package main

import (
	"os"
	"fmt"

	"github.com/cbegin/graven/commands"
	"github.com/cbegin/graven/version"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Version = version.Version
	app.Name = "graven"
	app.Usage = "A build automation tool for Go."

	// TODO:
	// new -- initializes new directory and project.yaml

	app.Commands = []cli.Command{
		commands.BuildCommand,
		commands.InfoCommand,
		commands.CleanCommand,
		commands.PackageCommand,
		commands.ReleaseCommand,
		commands.TestCommand,
		commands.FreezeCommand,
		commands.UnfreezeCommand,
		commands.InitCommand,
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		os.Exit(1)
	}
}
