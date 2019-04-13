package builder

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/blang/semver"

	"github.com/cbegin/graven/domain"
	"github.com/cbegin/graven/util"
)

type GoBuildTool struct{}

func (g *GoBuildTool) Build(outputPath string, project *domain.Project, artifact *domain.Artifact, target *domain.Target) error {
	v, err := getGoVersion()
	if err != nil {
		return err
	}

	if err := ensureVersion(project.GoVersion, v); err != nil {
		return err
	}

	fmt.Printf("Building %v/%v:%v\n", artifact.Classifier, target.Executable, project.Version)
	var c *exec.Cmd
	if target.Flags == "" {
		c = exec.Command("go", "build", "-o", filepath.Join(outputPath, target.Executable), target.Package)
	} else {
		c = exec.Command("go", "build", "-o", filepath.Join(outputPath, target.Executable), target.Flags, target.Package)
	}
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	c.Dir = project.ProjectPath()

	environment := project.Environment()
	for k, v := range util.MergeMaps(artifact.Environment, target.Environment) {
		environment = append(environment, fmt.Sprintf("%s=%s", k, v))
	}
	c.Env = environment

	if err := c.Run(); err != nil {
		fmt.Printf("FAILED to build %v/%v:%v\n", artifact.Classifier, target.Executable, project.Version)
		return fmt.Errorf("FAILED to build %v/%v:%v (%v)\n", artifact.Classifier, target.Executable, project.Version, err)
	}

	if !c.ProcessState.Success() {
		fmt.Printf("FAILED to build %v/%v:%v\n", artifact.Classifier, target.Executable, project.Version)
		return fmt.Errorf("FAILED to build %v/%v:%v (Build command exited in an error state. %v)\n", artifact.Classifier, target.Executable, project.Version, c)
	}

	fmt.Printf("Built %v/%v:%v\n", artifact.Classifier, target.Executable, project.Version)

	return err
}

func (g *GoBuildTool) Test(testPackage string, project *domain.Project) error {
	v, err := getGoVersion()
	if err != nil {
		return err
	}

	if err := ensureVersion(project.GoVersion, v); err != nil {
		return err
	}

	if err := runTestCommand(testPackage, project); err != nil {
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

	cmd := exec.Command("go", "test", "-v", "-covermode=atomic", coverOut, relativePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = project.ProjectPath()

	cmd.Env = project.Environment()

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

func runCoverageCommand(testPackage string, project *domain.Project) error {
	coverOutPath := fmt.Sprintf("%s.out", project.TargetPath("reports", testPackage))

	if _, err := os.Stat(coverOutPath); os.IsNotExist(err) {
		return nil
	}

	coverOut := fmt.Sprintf("-html=%s", coverOutPath)
	coverHtml := fmt.Sprintf("-o=%s.html", project.TargetPath("reports", testPackage))

	cmd := exec.Command("go", "tool", "cover", coverOut, coverHtml)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = project.ProjectPath()

	cmd.Env = project.Environment()

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

func ensureVersion(requiredVersion, actualVersion string) error {
	if requiredVersion != "" {
		r, err := semver.ParseRange(requiredVersion)
		if err != nil {
			return fmt.Errorf("Error parsing version range %v: %v.", requiredVersion, err)
		}
		//Needed for every new minor version of Go, naming convention is 1.10, 1.10.1 instead of 1.10.0, 1.10.1
		for len(strings.Split(actualVersion, ".")) < 3 {
			actualVersion = fmt.Sprintf("%s.0", actualVersion)
		}
		v, err := semver.Parse(actualVersion)
		if err != nil {
			return fmt.Errorf("Error parsing version %v: %v.", actualVersion, err)
		}
		if !r(v) {
			return fmt.Errorf("Go version %v does not satisfy %v.", actualVersion, requiredVersion)
		}
	}
	return nil
}

func getGoVersion() (string, error) {
	c := exec.Command("go", "version")
	buffer := bytes.NewBufferString("")
	c.Stdout = buffer
	if err := c.Run(); err != nil {
		return "", err
	}
	versionString := string(buffer.Bytes())

	parts := strings.Split(versionString, " ")
	if len(parts) < 4 {
		return "", fmt.Errorf("Version %v is invalid.", versionString)
	}

	versionPart := parts[2]
	if !strings.HasPrefix(versionPart, "go") {
		return "", fmt.Errorf("Version %v is invalid.", versionString)
	}

	version := versionPart[2:]
	return version, nil
}
