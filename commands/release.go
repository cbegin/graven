package commands

import (
	"github.com/urfave/cli"
	"github.com/cbegin/graven/domain"
	"gopkg.in/yaml.v2"
	"os"
	"io/ioutil"
	"fmt"
	"text/template"
)

var versionTemplate = `// This file was generated. It will be overwritten. Do not modify.
package version
var Version="{{.}}"`

var ReleaseCommand = cli.Command{
	Name: "release",
	Usage:       "project release",
	UsageText:   "release - release project",
	Description: "increments the revision and packages the release",
	Action: release,
}

func release(c *cli.Context) error {
	project := c.App.Metadata["project"].(*domain.Project)
	arg := c.Args().First()
	if err := bumpVersion(project, arg); err != nil {
		return err
	}

	versionPath := project.ProjectPath("version")
	versionFile := project.ProjectPath("version","version.go")
	_ = os.Mkdir(versionPath, 0755)
	file, _ := os.Create(versionFile);
	tmpl, _ := template.New("test").Parse(versionTemplate)
	tmpl.Execute(file, project.Version)

	return pkg(c)
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

	bytes, err := yaml.Marshal(project)
	if err != nil {
		return err
	}

	fileInfo, err := os.Stat(project.FilePath)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(project.FilePath, bytes, fileInfo.Mode())
	if err != nil {
		return err
	}

	return nil
}

