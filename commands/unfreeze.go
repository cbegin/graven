package commands

import (
	"fmt"
	"os"

	"github.com/cbegin/graven/domain"
	"github.com/cbegin/graven/util"
	"github.com/urfave/cli"
	"github.com/cbegin/graven/repotool"
	"github.com/cbegin/graven/vendortool"
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

	vendorTool, err := selectVendorTool(project)
	if err != nil {
		return err
	}

	err = vendorTool.LoadFile(project)
	if err != nil {
		return err
	}

	for _, p := range vendorTool.Dependencies() {
		sourceFile := project.ProjectPath(".freezer", p.ArchiveFileName())
		targetDir := project.ProjectPath("vendor", p.PackagePath())

		_, err := os.Stat(sourceFile)
		if os.IsNotExist(err) {
			for repoName, repo := range project.Repositories {
				if repo.HasRole(domain.RepositoryRoleDependency) {
					if repoTool, ok := repotool.RepoRegistry[repo.Type]; ok {
						if err := repoTool.DownloadDependency(project, repoName, sourceFile, vendortool.Coordinates(p)); err != nil {
							fmt.Println(err)
						} else {
							err = util.UnzipDir(sourceFile, targetDir)
							if err != nil {
								fmt.Println(err)
							}
							break
						}
					} else {
						fmt.Printf("Unkown repository type %v for %v\n", repo.Type, repoName)
					}
				}
			}
		} else {
			err = util.UnzipDir(sourceFile, targetDir)
			if err != nil {
				return err
			}
		}

	}

	return nil
}
