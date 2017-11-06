package commands

import (
	"fmt"
	"os"

	"github.com/cbegin/graven/buildtool"
	"github.com/cbegin/graven/domain"
	"github.com/cbegin/graven/util"
	"github.com/hashicorp/go-multierror"
	"github.com/urfave/cli"
)

var BuildCommand = cli.Command{
	Name:   "build",
	Usage:  "Builds the current project",
	Action: build,
}

func build(c *cli.Context) error {
	if err := clean(c); err != nil {
		return err
	}

	project, err := domain.FindProject()
	if err != nil {
		return err
	}

	if err := writeVersionFile(project); err != nil {
		fmt.Println(err)
	}

	var merr error
	for _, artifact := range project.Artifacts {
		a := artifact
		for _, target := range artifact.Targets {
			t := target
			err := buildTarget(project, &a, &t)
			if err != nil {
				merr = multierror.Append(merr, err)
			}
		}
	}

	return merr
}

func buildTarget(project *domain.Project, artifact *domain.Artifact, target *domain.Target) error {
	classifiedPath := project.TargetPath(artifact.Classifier)
	if _, err := os.Stat(classifiedPath); os.IsNotExist(err) {
		os.Mkdir(classifiedPath, 0755)
	}

	for _, resource := range append(project.Resources, artifact.Resources...) {
		resourcePath := project.ProjectPath(resource)
		baseProjectPath := project.ProjectPath()
		if len(resourcePath[len(baseProjectPath):]) < 1 {
			return fmt.Errorf("Resource path cannot be the entire project folder: %s", resourcePath)
		}
		err := util.CopyDir(resourcePath, classifiedPath)
		if err != nil {
			return err
		}
	}

	// TODO - Make this configurable
	var buildTool builder.BuildTool = &builder.GoBuildTool{}
	return buildTool.Build(classifiedPath, project, artifact, target)
}
