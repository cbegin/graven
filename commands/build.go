package commands

import "github.com/urfave/cli"

var BuildCommand = cli.Command{
	Name: "build",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name: "",
		},
	},
}
