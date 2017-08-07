package vendortool

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"path/filepath"

	"github.com/cbegin/graven/domain"
	"gopkg.in/yaml.v2"
	"github.com/cbegin/graven/util"
)

type GlidePackage struct {
	Name    string `yaml:"name,omitempty"`
	Version string `yaml:"version,omitempty"`
}

type GlideVendorTool struct {
	Imports []GlidePackage `yaml:"imports,omitempty"`
	TestImports []GlidePackage `yaml:"testImports,omitempty"`
}

func (g *GlideVendorTool) Name() string {
	return "glide"
}

func (g *GlidePackage) ArchiveFileName() string {
	chars := []string{"/", "\\", "."}
	name := g.Name
	for _, c := range chars {
		name = strings.Replace(name, c, "-", -1)
	}

	return fmt.Sprintf("%v-%v.zip", name, g.Version)
}

func (g *GlidePackage) PackagePath() string {
	osIndependentPath := filepath.Join(strings.FieldsFunc(g.Name, func(r rune) bool { return r == '\\' || r == '/' })...)
	return osIndependentPath
}

func (g *GlidePackage) Tag() string {
	return g.Version
}

func (g *GlideVendorTool) VendorFileExists(project *domain.Project) bool {
	vendorFilePath := project.ProjectPath("glide.lock")
	return util.PathExists(vendorFilePath)
}

func (g *GlideVendorTool) LoadFile(project *domain.Project) error {
	vendorFilePath := project.ProjectPath("glide.lock")
	vendorFile, err := os.Open(vendorFilePath)
	if err != nil {
		return err
	}
	vendorFileBytes, err := ioutil.ReadAll(vendorFile)
	err = yaml.Unmarshal(vendorFileBytes, g)
	if err != nil {
		return err
	}
	return nil
}

func (g *GlideVendorTool) Dependencies() []PackageDepencency {
	deps := make([]PackageDepencency, len(g.Imports) + len(g.TestImports))
	for i, dx := range g.Imports {
		d := dx
		deps[i] = &d
	}
	for i, dx := range g.TestImports {
		d := dx
		deps[len(g.Imports) + i] = &d
	}
	return deps
}
