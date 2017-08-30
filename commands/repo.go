package commands

import (
	"fmt"

	"github.com/cbegin/graven/domain"
	"github.com/cbegin/graven/repotool"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

var RepoCommand = cli.Command{
	Name:   "repo",
	Usage:  "Manages repository connections",
	Action: repo,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "login",
			Usage: "Prompts for repo login credentials.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "Name of the repo to manage.",
		},
	},
}

func repo(c *cli.Context) error {
	project, err := domain.FindProject()
	if err != nil {
		return err
	}

	repoName := c.String("name")
	if repoName == "" {
		return fmt.Errorf("No repo name specified.")
	}

	repository, found := project.Repositories[repoName]
	if !found {
		return fmt.Errorf("No repo named %v is found in project.", repoName)
	}

	if c.Bool("login") {
		if repoTool, ok := repotool.RepoRegistry[repository.Type]; ok {
			err := repoTool.Login(project, repoName)
			if err != nil {
				return err
			}
		}
	} else {
		s, err := yaml.Marshal(repository)
		if err != nil {
			return err
		}
		fmt.Printf("%v\n", string(s))
	}

	return nil
}
