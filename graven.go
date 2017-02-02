package main

import (
	"github.com/urfave/cli"
	"os"
	"fmt"
	"github.com/cbegin/graven/commands"
	"github.com/cbegin/graven/domain"
	"github.com/cbegin/graven/version"
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
		commands.ReleaseCommand,
	}

	p, err := domain.FindProject()
	if err != nil {
		fmt.Println("Could not find project.yaml in current or parent path.")
		return
	}

	app.Metadata = map[string]interface{}{"project":p}

	// new -- initializes new directory and project.yaml
	// test -- runs tests with flags, coverage etc.
	// deploy -- deploy one artifact to one repository
	// release [major|minor|patch] package, deploy each archive

	fmt.Printf("Project Path: %s\n", p.ProjectPath())

	app.Run(os.Args)
}