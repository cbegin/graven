package commands

import (
	"fmt"
	"path/filepath"
	"os"
	"strings"
	"os/exec"

	"github.com/urfave/cli"
	"github.com/cbegin/graven/domain"
)

var TestCommand = cli.Command{
	Name: "test",
	Usage:       "finds and runs tests in this project",
	Action: tester,
}

func tester(c *cli.Context) error {
	project := c.App.Metadata["project"].(*domain.Project)

	return filepath.Walk(project.ProjectPath(), getWalker(project))
}

func getWalker(project *domain.Project) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			subDir := path[len(project.ProjectPath()):]
			subDirParts := strings.Split(subDir, string(filepath.Separator))
			matches, _ := filepath.Glob(filepath.Join(path, "*_test*"));
			if len(matches) > 0 && !contains(subDirParts, map[string]struct{}{
				"vendor":struct{}{},
				"target":struct{}{},
				".git":struct{}{}}) {
				_ = runTestCommand(subDir, project)
			}
		}
		return nil
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

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Build error. %v", err)
	}

	if !cmd.ProcessState.Success() {
		return fmt.Errorf("Build command exited in an error state. %v", cmd)
	}
	return err
}
