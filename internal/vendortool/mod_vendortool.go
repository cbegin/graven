package vendortool

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cbegin/graven/internal/domain"
	"github.com/cbegin/graven/internal/util"

	"golang.org/x/mod/module"
)

type ModPackage struct {
	Path     string `json:"Path,omitempty"`
	Version  string `json:"Version,omitempty"`
	Info     string `json:"Info,omitempty"`
	GoMod    string `json:"GoMod,omitempty"`
	Zip      string `json:"Zip,omitempty"`
	Dir      string `json:"Dir,omitempty"`
	Sum      string `json:"Sum,omitempty"`
	GoModSum string `json:"GoModSum,omitempty"`
}

type ModVendorTool struct {
	Packages []ModPackage `json:"Packages,omitempty"`
}

func (g *ModVendorTool) Name() string {
	return "mod"
}

func (g *ModPackage) ArchiveFileName() string {
	unescaped := fmt.Sprintf("%v/@v/%v.zip", g.Path, g.Version)
	escapedPath, err := module.EscapePath(g.Path)
	if err != nil {
		fmt.Printf("Error, could not escape module path %v: %v\n", g.Path, err)
		return unescaped
	}
	escapedVersion, err := module.EscapeVersion(g.Version)
	if err != nil {
		fmt.Printf("Error, could not escape module version %v@%v: %v\n", g.Path, g.Version, err)
		return unescaped
	}
	return fmt.Sprintf("%v/@v/%v.zip", escapedPath, escapedVersion)
}

func (g *ModPackage) PackagePath() string {
	osIndependentPath := filepath.Join(strings.FieldsFunc(g.Path, func(r rune) bool { return r == '\\' || r == '/' })...)
	return osIndependentPath
}

func (g *ModPackage) Tag() string {
	return g.Version
}

func (g *ModVendorTool) VendorFileExists(project *domain.Project) bool {
	vendorFilePath := project.ProjectPath("go.mod")
	return util.PathExists(vendorFilePath)
}

func (g *ModVendorTool) LoadFile(project *domain.Project) error {
	c := exec.Command("go", "mod", "download", "-json")
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	c.Dir = project.ProjectPath()
	c.Env = project.Environment()
	bytes, err := c.Output()

	err = json.Unmarshal(fixJSONArray(bytes), g)
	if err != nil {
		return err
	}

	if err != nil {
		return fmt.Errorf("FAILED to download module information (%v)\n", err)
	}

	if !c.ProcessState.Success() {
		return fmt.Errorf("FAILED to download module information  (command exited in an error state. %v)\n", c)
	}

	return g.vendorModules(project)
}

func (g *ModVendorTool) vendorModules(project *domain.Project) error {
	c := exec.Command("go", "mod", "vendor")
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	c.Dir = project.ProjectPath()
	c.Env = project.Environment()
	_, err := c.Output()

	if err != nil {
		return fmt.Errorf("FAILED to vendor modules (%v)\n", err)
	}

	if !c.ProcessState.Success() {
		return fmt.Errorf("FAILED to vendor modules  (command exited in an error state. %v)\n", c)
	}
	return nil
}

func (g *ModVendorTool) Dependencies() []PackageDepencency {
	deps := make([]PackageDepencency, len(g.Packages))
	for i, dx := range g.Packages {
		d := dx
		deps[i] = &d
	}
	return deps
}

func fixJSONArray(notJSON []byte) []byte {
	var re = regexp.MustCompile(`}\s*{`)
	s := re.ReplaceAllString(string(notJSON), "},\n{")
	s = `{"Packages":[` + s + `]}`
	return []byte(s)
}
