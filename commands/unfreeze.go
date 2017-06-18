package commands

import (
	"fmt"
	"os"

	"github.com/cbegin/graven/domain"
	"github.com/cbegin/graven/util"
	"github.com/cbegin/graven/vendortool"
	"github.com/urfave/cli"
)

var UnfreezeCommand = cli.Command{
	Name:   "unfreeze",
	Usage:  "Unfreezes vendor dependencies",
	Action: unfreeze,
}

func unfreeze(c *cli.Context) error {
	project, err := domain.FindProject()
	if err != nil {
		return err
	}

	// TODO: Make this configurable
	var vendorTool vendortool.VendorTool = &vendortool.GovendorVendorTool{}

	err = vendorTool.LoadFile(project)
	if err != nil {
		return err
	}

	for _, p := range vendorTool.Dependencies() {
		sourceFile := project.ProjectPath(".freezer", p.ArchiveFileName())
		targetDir := project.ProjectPath("vendor", p.PackagePath())

		_, err := os.Stat(sourceFile)
		if !os.IsNotExist(err) {
			err = util.UnzipDir(sourceFile, targetDir)
			if err != nil {
				return err
			}
			fmt.Printf("%s => %s\n", p.ArchiveFileName(), p.PackagePath())
		} else {
			fmt.Printf("MISSING frozen dependency: %s => %s\n", p.ArchiveFileName(), p.PackagePath())
		}
	}

	return nil
}
