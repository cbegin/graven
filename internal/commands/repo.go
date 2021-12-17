package commands

import (
	"fmt"
	"github.com/cbegin/graven/internal/domain"
	"github.com/cbegin/graven/internal/repotool"
	"reflect"

	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

var RepoCommand = cli.Command{
	Name:  "repo",
	Usage: "Manages repository connections",
	Subcommands: []cli.Command{
		{
			Name:   "login",
			Usage:  "Logs into a repository.",
			Action: repoLogin,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name",
					Usage: "Repository name.",
				},
			},
		},
		{
			Name:   "validate",
			Usage:  "Tests repo settings and authentication credentials.",
			Action: repoValidate,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name",
					Usage: "Repository name.",
				},
			},
		},
	},
}

func repoLogin(c *cli.Context) error {
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

	if repoTool, ok := repotool.RepoRegistry[repository.Type]; ok {
		err := repoTool.Login(project, repoName)
		if err != nil {
			return err
		}
	} else {
		fmt.Printf("Unknown repo type: %v. Expected one of %v", repository.Type, reflect.ValueOf(repotool.RepoRegistry).MapKeys())
	}

	return nil
}

func repoValidate(c *cli.Context) error {
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

	if repoTool, ok := repotool.RepoRegistry[repository.Type]; ok {
		err := repoTool.LoginTest(project, repoName)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Unknown repo type: %v. Expected one of %v", repository.Type, reflect.ValueOf(repotool.RepoRegistry).MapKeys())
	}

	s, err := yaml.Marshal(repository)
	if err != nil {
		return err
	}
	fmt.Printf("# Validation successful.\n%v", string(s))

	return nil
}
