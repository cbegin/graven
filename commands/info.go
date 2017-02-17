package commands

import (
	"fmt"

	"github.com/urfave/cli"
	"github.com/cbegin/graven/domain"
	"gopkg.in/yaml.v2"
)

var InfoCommand = cli.Command{
	Name: "info",
	Usage:       "prints the known information about a project",
	Action: info,
}

func info(c *cli.Context) error {
	project := c.App.Metadata["project"].(*domain.Project)

	bytes, err := yaml.Marshal(project)

	fmt.Println(string(bytes))
	return err
}