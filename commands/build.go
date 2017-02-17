package commands

import (
	"os/exec"
	"os"
	"fmt"
	"path"
	"sync"

	"github.com/hashicorp/go-multierror"
	"github.com/urfave/cli"
	"github.com/cbegin/graven/domain"
)

var BuildCommand = cli.Command{
	Name:        "build",
	Usage:       "builds the current project",
	Action: build,
}

func build(c *cli.Context) error {
	if err := clean(c); err != nil {
		return err
	}

	var merr error
	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}
	project := c.App.Metadata["project"].(*domain.Project)
	for _, artifact := range project.Artifacts {
		for _, target := range artifact.Targets {
			wg.Add(len(artifact.Targets))
			go func() {
				defer wg.Done()
				err := buildTarget(project, &artifact, &target)
				if err != nil {
					mutex.Lock()
					merr = multierror.Append(merr, err)
					mutex.Unlock()
				}

			}()
		}
	}
	wg.Wait()

	return merr
}

func buildTarget(project *domain.Project, artifact *domain.Artifact, target *domain.Target) error {
	classifiedPath := project.TargetPath(artifact.Classifier)
	if _, err := os.Stat(classifiedPath); os.IsNotExist(err) {
		os.Mkdir(classifiedPath, 0755)
	}

	for _, resource := range artifact.Resources {
		resourcePath := project.ProjectPath(resource)
		baseProjectPath := project.ProjectPath()
		if len(resourcePath[len(baseProjectPath):]) < 1 {
			return fmt.Errorf("Resource path cannot be the entire project folder: %s", resourcePath)
		}
		err := domain.CopyDir(resourcePath, classifiedPath)
		if err != nil {
			return err
		}
	}

	return runBuildCommand(classifiedPath, project, target)
}

func runBuildCommand(classifiedPath string, project *domain.Project, target *domain.Target) error {
	cmd := exec.Command("go", "build", "-o", path.Join(classifiedPath, target.Executable), target.Flags, target.Package)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = project.ProjectPath()

	environment := []string{}
	for k, v := range target.Environment {
		environment = append(environment, fmt.Sprintf("%s=%s", k, v))
	}
	gopath, _ := os.LookupEnv("GOPATH")
	environment = append(environment, fmt.Sprintf("%s=%s", "GOPATH", gopath))
	cmd.Env = environment

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Build error. %v", err)
	}

	if !cmd.ProcessState.Success() {
		return fmt.Errorf("Build command exited in an error state. %v", cmd)
	}
	return err
}