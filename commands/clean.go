package commands

import (
	"os"

	"github.com/urfave/cli"
	"github.com/cbegin/graven/domain"
)

var CleanCommand = cli.Command{
	Name: "clean",
	Usage:       "cleans the target directory and its contents",
	Action: clean,
}

func clean(c *cli.Context) error {
	project := c.App.Metadata["project"].(*domain.Project)
	return os.RemoveAll(project.TargetPath())
}