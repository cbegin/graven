package commands

import (
	"fmt"

	"github.com/cbegin/graven/internal/domain"

	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

var InfoCommand = cli.Command{
	Name:   "info",
	Usage:  "Prints the known information about a project",
	Action: info,
}

func info(*cli.Context) error {
	project, err := domain.FindProject()
	if err != nil {
		return err
	}

	bytes, err := yaml.Marshal(project)

	fmt.Println(string(bytes))
	return err
}
