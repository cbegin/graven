package commands

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
	"github.com/cbegin/graven/domain"
	"github.com/cbegin/graven/util"
)

var UnfreezeCommand = cli.Command{
	Name: "unfreeze",
	Usage:       "Unfreezes vendor dependencies",
	Action: unfreeze,
}

func unfreeze(c *cli.Context) error {
	project, err := domain.FindProject()
	if err != nil {
		return err
	}

	govendorFile, err := domain.ReadGovendorFile(project)
	if err != nil {
		return err
	}

	for _, p := range govendorFile.Packages {
		sourceFile := project.ProjectPath(".freezer", p.ArchiveFileName())
		targetDir := project.ProjectPath("vendor", p.Path)
		frozenFile := project.ProjectPath("vendor", p.Path, ".frozen")

		_, err := os.Stat(sourceFile)
		if !os.IsNotExist(err) {
			err = util.Unzip(sourceFile, targetDir)
			if err != nil {
				return err
			}
			fmt.Printf("%s => %s\n", p.ArchiveFileName(), p.Path)
		}

		_ = os.Remove(frozenFile)
	}

	return nil
}