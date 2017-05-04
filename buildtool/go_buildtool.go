package builder

import (
	"os/exec"
	"path"
	"path/filepath"
	"fmt"
	"os"

	"github.com/cbegin/graven/domain"
	"github.com/cbegin/graven/util"
)

type GoBuildTool struct {}

func (g *GoBuildTool) Build(outputPath string, project *domain.Project, artifact *domain.Artifact, target *domain.Target) error {
	fmt.Printf("Building %v/%v:%v\n", artifact.Classifier, target.Executable, project.Version)
	defer fmt.Printf("Done %v/%v:%v\n", artifact.Classifier, target.Executable, project.Version)
	var c *exec.Cmd
	if target.Flags == "" {
		c = exec.Command("go", "build", "-o", path.Join(outputPath, target.Executable), target.Package)
	} else {
		c = exec.Command("go", "build", "-o", path.Join(outputPath, target.Executable), target.Flags, target.Package)
	}
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	c.Dir = project.ProjectPath()

	environment := []string{}
	for k, v := range util.MergeMaps(artifact.Environment, target.Environment) {
		environment = append(environment, fmt.Sprintf("%s=%s", k, v))
	}
	gopath, _ := os.LookupEnv("GOPATH")
	environment = append(environment, fmt.Sprintf("%s=%s", "GOPATH", gopath))
	c.Env = environment

	err := c.Run()
	if err != nil {
		return fmt.Errorf("Build error. %v", err)
	}

	if !c.ProcessState.Success() {
		return fmt.Errorf("Build command exited in an error state. %v", c)
	}
	return err

}

func  (g *GoBuildTool) Test(testPackage string, project *domain.Project) error {
	err := runTestCommand(testPackage, project)
	if err != nil {
		return err
	}
	return runCoverageCommand(testPackage, project)
}

func runTestCommand(testPackage string, project *domain.Project) error {
	relativePath := "." + testPackage

	coverPath := project.TargetPath("reports", testPackage)

	if err := os.MkdirAll(filepath.Dir(coverPath), 0777); err != nil {
		return err
	}


	coverOut := fmt.Sprintf("-coverprofile=%s.out", coverPath)

	cmd := exec.Command("go", "test", "-v", "-parallel=4", "-p=4", "-covermode=atomic", coverOut, relativePath)
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
