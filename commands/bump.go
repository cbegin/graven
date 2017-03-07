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

var BumpCommand = cli.Command{
	Name: "bump",
	Usage:       "Manage the version (major, minor, patch) and clear or set qualifier (e.g. DEV)",
	Description: `
	Bump manages incremental version updates using simple semantic versioning practices.

	The valid arguments are:

	major - bumps the major version X._._
	minor - bumps the minor version _.X._
	patch - bumps the patch version _._.X
	clear - clears the qualifier if any
	_____ - Anything else, sets the qualifier (e.g. SNAPSHOT, DEV, ALPHA)`,
	Action: bump,
}

func bump(c *cli.Context) error {
	project, err := domain.FindProject()
	if err != nil {
		return err
	}

	oldVersion := project.Version
	arg := c.Args().First()
	if arg == "" {
		fmt.Printf("%v\n", oldVersion)
		return nil
	}
	if err := bumpVersion(project, arg); err != nil {
		return err
	}
	if err := writeVersionFile(project); err != nil {
		return err
	}
	newVersion := project.Version
	fmt.Printf("Version changed from %v to %v\n", oldVersion, newVersion)
	return nil
}

func writeVersionFile(project *domain.Project) error {
	versionPath := project.ProjectPath(versionPackage)
	versionFile := project.ProjectPath(versionPackage, versionFileName)

	if err := validateHeader(versionFile); err != nil {
		return err
	}

	_ = os.Mkdir(versionPath, 0755) // ignore error. we'll catch file errors later

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
	if os.IsNotExist(err) {
		// it's okay if the file doesn't exist
		return nil
	} else if err != nil {
		// but fail for any other reason
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
	case "clear":
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

