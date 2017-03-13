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
	if err := clean(c); err != nil {
		return err
	}

	project, err := domain.FindProject()
	if err != nil {
		return err
	}

	err = os.MkdirAll(project.TargetPath("reports"), 0755)
	if err != nil {
		return fmt.Errorf("Could not create reports directory. %v", err)
	}

	var merr error
	if err := filepath.Walk(project.ProjectPath(), getTestWalkerFunc(project, &merr)); err != nil {
		merr = multierror.Append(merr, err)
	}
	return merr
}

func getTestWalkerFunc(project *domain.Project, merr *error) filepath.WalkFunc {
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
					*merr = multierror.Append(*merr, err)
				}
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

	coverOut := fmt.Sprintf("-coverprofile=%s.out", project.TargetPath("reports", testPackage))

	cmd := exec.Command("go", "test", "-covermode=atomic", coverOut, relativePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = project.ProjectPath()

	environment := []string{}
	gopath, _ := os.LookupEnv("GOPATH")
	path, _ := os.LookupEnv("PATH")
	environment = append(environment, fmt.Sprintf("%s=%s", "GOPATH", gopath))
	environment = append(environment, fmt.Sprintf("%s=%s", "PATH", path))
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

	return runCoverageCommand(testPackage, project)
}

func runCoverageCommand(testPackage string, project *domain.Project) error {
	coverOut := fmt.Sprintf("-html=%s.out", project.TargetPath("reports", testPackage))
	coverHtml := fmt.Sprintf("-o=%s.html", project.TargetPath("reports", testPackage))

	cmd := exec.Command("go", "tool", "cover", coverOut, coverHtml)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = project.ProjectPath()

	environment := []string{}
	gopath, _ := os.LookupEnv("GOPATH")
	path, _ := os.LookupEnv("PATH")
	environment = append(environment, fmt.Sprintf("%s=%s", "GOPATH", gopath))
	environment = append(environment, fmt.Sprintf("%s=%s", "PATH", path))
	cmd.Env = environment


	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Error generating coverage HTML. %v", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("Error state returned after generating coverage HTML. %v", err)
	}

	if !cmd.ProcessState.Success() {
		return fmt.Errorf("Build command exited in an error state. %v", cmd)
	}

	return nil
}
