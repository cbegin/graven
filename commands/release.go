package commands

import (
	"github.com/cbegin/graven/domain"
	"github.com/cbegin/graven/repotool"
	"github.com/cbegin/graven/vcstool"
	"github.com/urfave/cli"
)

var ReleaseCommand = cli.Command{
	Name:   "release",
	Usage:  "Releases artifacts to repositories",
	Action: release,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "login",
			Usage: "Prompts for repo login credentials.",
		},
	},
}

func release(c *cli.Context) error {
	//TODO: Make this configurable
	var repoTool repotool.RepoTool = &repotool.GithubRepoTool{}

	if c.Bool("login") {
		return repoTool.Login()
	}

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

	return repoTool.Release(project)
}
