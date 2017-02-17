package commands

import (
	"github.com/urfave/cli"
	"github.com/cbegin/graven/domain"
	"encoding/json"
	"io/ioutil"
	"os"
	"fmt"
	"github.com/cbegin/graven/util"
)

var FreezeCommand = cli.Command{
	Name:        "freeze",
	Usage:       "freeze",
	UsageText:   "freeze - freezes dependencies",
	Description: "freezes dependencies to avoid having to check in source",
	Action: freeze,
}

func freeze(c *cli.Context) error {
	project := c.App.Metadata["project"].(*domain.Project)

	govendorFile, err := domain.ReadGovendorFile(project)

	os.Mkdir(project.ProjectPath(".freezer"), 0755) // ignore errors, we'll catch subsequent ones
	if err != nil {
		return err
	}
	for _, p := range govendorFile.Packages {
		sourcePath := project.ProjectPath("vendor", p.Path)
		targetFile := project.ProjectPath(".freezer", p.ArchiveFileName())
		frozenFile := project.ProjectPath("vendor", p.Path, ".frozen")

		j, err := json.Marshal(p)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(frozenFile, j, 0644)
		if err != nil {
			return err
		}

		err = util.ZipDir(sourcePath, targetFile)
		if err != nil {
			return err
		}

		_ = os.Remove(frozenFile)

		fmt.Printf("%s => %s\n", p.Path, p.ArchiveFileName())
	}

	return err
}

