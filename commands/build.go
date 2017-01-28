package commands

import (
	"github.com/urfave/cli"
	"github.com/cbegin/graven/domain"
	"os/exec"
	"fmt"
	"os"
)

var BuildCommand = cli.Command{
	Name: "build",
	//Flags: []cli.Flag{
	//	cli.StringFlag{
	//		Name: "",
	//	},
	//},
	Usage:       "build project",
	UsageText:   "build - build project",
	Description: "find the nearest project.yaml in the current directory tree and builds",
	Action: build,
}

func build(c *cli.Context) error {
	project := c.App.Metadata["project"].(*domain.Project)
	
	for _, pkg := range project.Packages {
		cmd := exec.Command("go", "build", "-o", pkg.Exec, pkg.Path)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}