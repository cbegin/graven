package main

import (
	"fmt"
	commands2 "github.com/cbegin/graven/internal/commands"
	"os"

	"github.com/cbegin/graven/version"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Version = version.Version
	app.Name = "graven"
	app.Usage = "A build automation tool for Go."

	app.Commands = []cli.Command{
		commands2.BuildCommand,
		commands2.InfoCommand,
		commands2.CleanCommand,
		commands2.PackageCommand,
		commands2.BumpCommand,
		commands2.TestCommand,
		commands2.InitCommand,
		commands2.ReleaseCommand,
		commands2.RepoCommand,
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		os.Exit(1)
	}
}
