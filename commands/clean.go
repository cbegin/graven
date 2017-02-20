package commands

import (
	"os"

	"github.com/urfave/cli"
	"github.com/cbegin/graven/domain"
)

var CleanCommand = cli.Command{
	Name: "clean",
	Usage:       "Cleans the target directory and its contents",
	Action: clean,
}

func clean(c *cli.Context) error {
	project, err := domain.FindProject()
	if err != nil {
		return err
	}
	return os.RemoveAll(project.TargetPath())
}