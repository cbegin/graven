package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cbegin/graven/domain"
	"github.com/cbegin/graven/util"
	"github.com/cbegin/graven/vendortool"
	"github.com/urfave/cli"
)

var FreezeCommand = cli.Command{
	Name:   "freeze",
	Usage:  "Freezes vendor dependencies to avoid having to check in source",
	Action: freeze,
}

func freeze(c *cli.Context) error {
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

	freezerPath := project.ProjectPath(".freezer")

	if _, err := os.Stat(freezerPath); !os.IsNotExist(err) {
		if err := os.RemoveAll(freezerPath); err != nil {
			return fmt.Errorf("Could not clean .freezer: %v", err)
		}
	}

	if err := os.Mkdir(freezerPath, 0755); err != nil {
		return fmt.Errorf("Could not make .freezer: %v", err)
	}

	for _, p := range vendorTool.Dependencies() {
		sourcePath := project.ProjectPath("vendor", p.PackagePath())
		targetFile := project.ProjectPath(".freezer", p.ArchiveFileName())
		frozenFile := project.ProjectPath("vendor", p.PackagePath(), ".frozen")

		if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
			fmt.Printf("MISSING dependency %v\n", p.PackagePath())
			continue
		}

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

		fmt.Printf("%s => %s\n", p.PackagePath(), p.ArchiveFileName())
	}

	return err
}
