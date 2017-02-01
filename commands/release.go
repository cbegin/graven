package commands

import (
	"github.com/urfave/cli"
	"github.com/cbegin/graven/domain"
)

var ReleaseCommand = cli.Command{
	Name: "release",
	Usage:       "project release",
	UsageText:   "release - release project",
	Description: "increments the revision and packages the release",
	Action: release,
}

func release(c *cli.Context) error {
	_ = c.App.Metadata["project"].(*domain.Project)

	return nil
}


