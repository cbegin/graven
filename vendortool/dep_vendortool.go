package vendortool

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"path/filepath"

	"github.com/cbegin/graven/domain"
	"github.com/pelletier/go-toml"
	"github.com/cbegin/graven/util"
)

type DepPackage struct {
	Name     string `toml:"name,omitempty"`
	Revision string `toml:"revision,omitempty"`
}

type DepVendorTool struct {
	Imports []DepPackage `toml:"projects,omitempty"`
}

func (g *DepVendorTool) Name() string {
	return "dep"
}

func (g *DepPackage) ArchiveFileName() string {
	chars := []string{"/", "\\", "."}
	name := g.Name
	for _, c := range chars {
		name = strings.Replace(name, c, "-", -1)
	}

	return fmt.Sprintf("%v-%v.zip", name, g.Revision)
}

func (g *DepPackage) PackagePath() string {
	osIndependentPath := filepath.Join(strings.FieldsFunc(g.Name, func(r rune) bool { return r == '\\' || r == '/' })...)
	return osIndependentPath
}

func (g *DepPackage) Tag() string {
	return g.Revision
}

func (g *DepVendorTool) VendorFileExists(project *domain.Project) bool {
	vendorFilePath := project.ProjectPath("Gopkg.lock")
	return util.PathExists(vendorFilePath)
}

func (g *DepVendorTool) LoadFile(project *domain.Project) error {
	vendorFilePath := project.ProjectPath("Gopkg.lock")
	vendorFile, err := os.Open(vendorFilePath)
	if err != nil {
		return err
	}
	vendorFileBytes, err := ioutil.ReadAll(vendorFile)
	err = toml.Unmarshal(vendorFileBytes, g)
	if err != nil {
		return err
	}
	return nil
}

func (g *DepVendorTool) Dependencies() []PackageDepencency {
	deps := make([]PackageDepencency, len(g.Imports))
	for i, dx := range g.Imports {
		d := dx
		deps[i] = &d
	}
	return deps
}
