package main

import (
	"fmt"
	"os"

	"github.com/cbegin/graven/internal/commands"
	"github.com/cbegin/graven/version"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Version = version.Version
	app.Name = "graven"
	app.Usage = "A build automation tool for Go."

	app.Commands = []cli.Command{
		commands.BuildCommand,
		commands.InfoCommand,
		commands.CleanCommand,
		commands.PackageCommand,
		commands.BumpCommand,
		commands.TestCommand,
		commands.InitCommand,
		commands.ReleaseCommand,
		commands.RepoCommand,
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		os.Exit(1)
	}
}
