package commands

import (
	"os"
	"io/ioutil"
	"fmt"
	"text/template"

	"github.com/urfave/cli"
	"github.com/cbegin/graven/domain"
	"gopkg.in/yaml.v2"
)

const (
	versionPackage = "version"
	versionFileName = "version.go"
	versionTemplate = `// graven - This file was generated. It will be overwritten. Do not modify.
package {{.Package}}
var Version="{{.Version}}"`
)

var ReleaseCommand = cli.Command{
	Name: "release",
	Usage:       "Increments the revision and packages the release",
	Action: release,
}

func release(c *cli.Context) error {
	project, err := domain.FindProject()
	if err != nil {
		return err
	}

	arg := c.Args().First()
	if err := bumpVersion(project, arg); err != nil {
		return err
	}
	if err := writeVersionFile(project); err != nil {
		return err
	}
	return pkg(c)
}

func writeVersionFile(project *domain.Project) error {
	versionPath := project.ProjectPath(versionPackage)
	versionFile := project.ProjectPath(versionPackage, versionFileName)

	if err := validateHeader(versionFile); err != nil {
		return err
	}

	_ = os.Mkdir(versionPath, 0755) // ignore error, we'll catch file errors

	file, err := os.Create(versionFile);
	defer file.Close()
	if err != nil {
		return err
	}
	tmpl, err := template.New("version").Parse(versionTemplate)
	if err != nil {
		return err
	}

	tmpl.Execute(file, struct {
		Version string
		Package string
	}{
		Version: project.Version,
		Package: versionPackage,
	})
	return nil
}

func validateHeader(versionFile string) error {
	const headerLength = 10
	file, err := os.Open(versionFile)
	defer file.Close()
	if err != nil {
		return err
	}
	var buffer = make([]byte, headerLength)
	_, err = file.Read(buffer)
	if err != nil {
		return err
	}
	if string(buffer) != versionTemplate[:headerLength] {
		return fmt.Errorf("Header in %v doesn't match, so graven won't overwrite it.", versionFile)
	}
	return nil
}

func bumpVersion(project *domain.Project, arg string) error {
	version := domain.Version{}

	err := version.Parse(project.Version)
	if err != nil {
		return fmt.Errorf("Error parsing version: %v", err)
	}

	switch arg {
	case "major":
		version.Major++
		version.Minor = 0
		version.Patch = 0
		version.Qualifier = ""
	case "minor":
		version.Minor++
		version.Patch = 0
		version.Qualifier = ""
	case "patch":
		version.Patch++
		version.Qualifier = ""
	case "":
	default:
		version.Qualifier = arg
	}

	project.Version = version.ToString()

	// Silly workaround for YAML libs inability to ignore fields.
	projectFilePath := project.FilePath
	project.FilePath = ""
	defer func() {
		project.FilePath = projectFilePath
	}()

	bytes, err := yaml.Marshal(project)
	if err != nil {
		return err
	}

	fileInfo, err := os.Stat(projectFilePath)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(projectFilePath, bytes, fileInfo.Mode())
	if err != nil {
		return err
	}

	return nil
}

