package commands

import (
	"fmt"
	"os"

	"github.com/cbegin/graven/domain"
	"github.com/cbegin/graven/repotool"
	"github.com/cbegin/graven/vcstool"
	"github.com/urfave/cli"
)

var ReleaseCommand = cli.Command{
	Name:   "release",
	Usage:  "Releases artifacts to repositories",
	Action: release,
}

func release(c *cli.Context) error {
	project, err := domain.FindProject()
	if err != nil {
		return err
	}

	//TODO: Make this configurable
	var vcsTool vcstool.VCSTool = &vcstool.GitVCSTool{}
	if os.Getenv("TESTRELEASE") == "" {
		if err := vcsTool.VerifyRepoState(project); err != nil {
			return err
		}
	}
	if err := pkg(c); err != nil {
		return err
	}
	if os.Getenv("TESTRELEASE") == "" {
		tagName := fmt.Sprintf("v%s", project.Version)
		if err := vcsTool.Tag(project, tagName); err != nil {
			return err
		}
	}

	for repoName, repo := range project.Repositories {
		if repo.HasRole(domain.RepositoryRoleRelease) {
			if repoTool, ok := repotool.RepoRegistry[repo.Type]; ok {
				if err := repoTool.Release(project, repoName); err != nil {
					return err
				}
			} else {
				fmt.Printf("Unkown repository type %v for %v\n", repo.Type, repoName)
			}
		}
	}

	return nil
}
