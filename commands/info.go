package commands

import (
	"github.com/urfave/cli"
	"github.com/cbegin/graven/domain"
	"gopkg.in/yaml.v2"
	"fmt"
)

var InfoCommand = cli.Command{
	Name: "info",
	Usage:       "project info",
	UsageText:   "info - project info",
	Description: "prints the known information about a project",
	Action: info,
}

func info(c *cli.Context) error {
	project := c.App.Metadata["project"].(*domain.Project)

	bytes, err := yaml.Marshal(project)

	fmt.Println(string(bytes))
	return err
}