package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cbegin/graven/domain"
	"github.com/cbegin/graven/repotool"
	"github.com/cbegin/graven/util"
	"github.com/cbegin/graven/vendortool"
	"github.com/urfave/cli"
)

var supportedVendorTools = []vendortool.VendorTool{
	&vendortool.ModVendorTool{},
	&vendortool.GovendorVendorTool{},
	&vendortool.GlideVendorTool{},
	&vendortool.DepVendorTool{},
}

var FreezeCommand = cli.Command{
	Name:   "freeze",
	Usage:  "Freezes vendor dependencies to avoid having to check in source",
	Action: freeze,
}

func selectVendorTool(project *domain.Project) (vendortool.VendorTool, error) {
	for _, vt := range supportedVendorTools {
		if vt.VendorFileExists(project) {
			fmt.Printf("Using %v vendor tool.\n", vt.Name())
			return vt, nil
		}
	}
	return nil, fmt.Errorf("Could not find supported vendor file.")
}

func freeze(c *cli.Context) error {
	project, err := domain.FindProject()
	if err != nil {
		return err
	}

	vendorTool, err := selectVendorTool(project)
	if err != nil {
		return err
	}

	err = vendorTool.LoadFile(project)
	if err != nil {
		return err
	}

	freezerPath := project.ProjectPath(".modules")

	if _, err := os.Stat(freezerPath); !os.IsNotExist(err) {
		if err := os.RemoveAll(freezerPath); err != nil {
			return fmt.Errorf("Could not clean .modules: %v", err)
		}
	}

	if err := os.Mkdir(freezerPath, 0755); err != nil {
		return fmt.Errorf("Could not make .modules: %v", err)
	}

	for _, p := range vendorTool.Dependencies() {
		sourcePath := project.ProjectPath("vendor", p.PackagePath())
		targetFile := project.ProjectPath(".modules", p.ArchiveFileName())

		if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
			fmt.Printf("Skipping %v wasn't found in vendor folder (may be an empty dependency).\n", p.PackagePath())
			continue
		}

		dirString, _ := filepath.Split(targetFile)
		if dirString != "" {
			os.MkdirAll(dirString, 0755)
		}

		err = util.ZipDir(sourcePath, targetFile, false)
		if err != nil {
			return err
		}

		for repoName, repo := range project.Repositories {
			if repo.HasRole(domain.RepositoryRoleDependency) {
				if repoTool, ok := repotool.RepoRegistry[repo.Type]; ok {
					if err := repoTool.UploadDependency(project, repoName, targetFile, vendortool.Coordinates(p)); err != nil {
						fmt.Println(err)
					}

				} else {
					fmt.Printf("Unkown repository type %v for %v\n", repo.Type, repoName)
					break
				}
			}
		}
	}

	return err
}
