package commands

import (
	"github.com/urfave/cli"
	"github.com/cbegin/graven/domain"
	"os/exec"
	"os"
	"fmt"
	"path"
)

var BuildCommand = cli.Command{
	Name: "build",
	//Flags: []cli.Flag{
	//	cli.StringFlag{
	//		Name: "",
	//	},
	//},
	Usage:       "build project",
	UsageText:   "build - build project",
	Description: "find the nearest project.yaml in the current directory tree and builds",
	Action: build,
}

func build(c *cli.Context) error {
	project := c.App.Metadata["project"].(*domain.Project)
	
	for _, artifact := range project.Artifacts {
		for _, target := range artifact.Targets {
			outpath := path.Join(project.TargetPath(artifact.Classifier))
			if _, err := os.Stat(outpath); os.IsNotExist(err) {
				os.Mkdir(outpath, 0755)
			}
			cmd := exec.Command("go", "build", "-o", path.Join(outpath, target.Executable), target.Flags, target.Package)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin

			environment := []string{}
			for k, v := range target.Environment {
				environment = append(environment, fmt.Sprintf("%s=%s", k, v))
			}
			gopath, _ := os.LookupEnv("GOPATH")
			environment = append(environment, fmt.Sprintf("%s=%s", "GOPATH", gopath))
			cmd.Env = environment

			err := cmd.Run()
			if err != nil {
				return err
			}
		}
	}

	return nil
}