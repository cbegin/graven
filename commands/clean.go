package commands

import (
	"github.com/urfave/cli"
	"github.com/cbegin/graven/domain"
	"os"
)

var CleanCommand = cli.Command{
	Name: "clean",
	Usage:       "project clean",
	UsageText:   "clean - project clean",
	Description: "cleans the target directory and its contents",
	Action: clean,
}

func clean(c *cli.Context) error {
	project := c.App.Metadata["project"].(*domain.Project)

	return os.RemoveAll(project.TargetPath())
}