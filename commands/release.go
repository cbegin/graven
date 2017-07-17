package commands

import (
	"fmt"

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
	if err := vcsTool.VerifyRepoState(project); err != nil {
		return err
	}
	if err := pkg(c); err != nil {
		return err
	}
	tagName := fmt.Sprintf("v%s", project.Version)
	if err := vcsTool.Tag(project, tagName); err != nil {
		return err
	}

	//TODO: Make this configurable

	repoToolMap := map[string]repotool.RepoTool{}
	repoToolMap["github"] = &repotool.GithubRepoTool{}
	repoToolMap["maven"] = &repotool.MavenRepoTool{}
	for repoName, repo := range project.Repositories {
		if repoTool, ok := repoToolMap[repo["type"]]; ok {
			if err := repoTool.Release(project, repoName); err != nil {
				return err
			}
		} else {
			fmt.Println("Unkown repository type %v: ", repo["type"])
		}
	}

	if bumpVersion(project, "patch"); err != nil {
		return err
	}
	if bumpVersion(project, "DEV"); err != nil {
		return err
	}

	return nil
}
