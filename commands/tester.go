package commands

import (
	"fmt"
	"path/filepath"
	"os"
	"strings"
	"os/exec"

	"github.com/urfave/cli"
	"github.com/cbegin/graven/domain"
	"github.com/hashicorp/go-multierror"
)

var TestCommand = cli.Command{
	Name: "test",
	Usage:       "Finds and runs tests in this project",
	Action: tester,
}

func tester(c *cli.Context) error {
	project, err := domain.FindProject()
	if err != nil {
		return err
	}

	var merr error
	if err := filepath.Walk(project.ProjectPath(), getTestWalkerFunc(project, merr)); err != nil {
		merr = multierror.Append(merr, err)
	}
	return merr
}

func getTestWalkerFunc(project *domain.Project, merr error) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			subDir := path[len(project.ProjectPath()):]
			subDirParts := strings.Split(subDir, string(filepath.Separator))
			matches, _ := filepath.Glob(filepath.Join(path, "*_test*"));
			if len(matches) > 0 && !contains(subDirParts, map[string]struct{}{
				"vendor":struct{}{},
				"target":struct{}{},
				".git":struct{}{}}) {
				if err := runTestCommand(subDir, project); err != nil {
					merr = multierror.Append(merr, err)
				}
			}
		}
		return merr
	}
}

func contains(strings []string, exclusions map[string]struct{}) bool {
	for _, a := range strings {
		if _, ok := exclusions[a]; ok {
			return true
		}
	}
	return false
}

func runTestCommand(testPackage string, project *domain.Project) error {
	relativePath := "." + testPackage
	fmt.Println(relativePath)
	cmd := exec.Command("go", "test", "-p", "4", "-cover", "-parallel", "4", relativePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = project.ProjectPath()

	environment := []string{}
	gopath, _ := os.LookupEnv("GOPATH")
	environment = append(environment, fmt.Sprintf("%s=%s", "GOPATH", gopath))
	cmd.Env = environment


	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Error starting tests. %v", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("Error state returned after tests. %v", err)
	}

	if !cmd.ProcessState.Success() {
		return fmt.Errorf("Build command exited in an error state. %v", cmd)
	}

	return nil
}
