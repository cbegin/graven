package vendortool

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"path/filepath"

	"github.com/cbegin/graven/domain"
	"github.com/cbegin/graven/util"
)

type GovendorPackage struct {
	Origin       string `json:"oriin,omitempty"`
	Path         string `json:"path,omitempty"`
	Revision     string `json:"revision,omitempty"`
	RevisionTime string `json:"revisionTime,omitempty"`
	Comment      string `json:"comment,omitempty"`
}

type GovendorVendorTool struct {
	Packages []GovendorPackage `json:"package,omitempty"`
}

func (g *GovendorVendorTool) Name() string {
	return "govendor"
}

func (g *GovendorPackage) ArchiveFileName() string {
	chars := []string{"/", "\\", "."}
	name := g.Path
	for _, c := range chars {
		name = strings.Replace(name, c, "-", -1)
	}

	var revision = ""
	if g.Revision == "" {
		t, _ := time.Parse(time.RFC3339, g.RevisionTime)
		revision = t.Format("2006.01.02-150405")
	} else {
		revision = g.Revision
	}
	return fmt.Sprintf("%v-%v.zip", name, revision)
}

func (g *GovendorPackage) PackagePath() string {
	osIndependentPath := filepath.Join(strings.FieldsFunc(g.Path, func(r rune) bool { return r == '\\' || r == '/' })...)
	return osIndependentPath
}

func (g *GovendorPackage) Tag() string {
	return g.Revision
}


func (g *GovendorVendorTool) VendorFileExists(project *domain.Project) bool {
	vendorFilePath := project.ProjectPath("vendor", "vendor.json")
	return util.PathExists(vendorFilePath)
}

func (g *GovendorVendorTool) LoadFile(project *domain.Project) error {
	vendorFilePath := project.ProjectPath("vendor", "vendor.json")
	vendorFile, err := os.Open(vendorFilePath)
	if err != nil {
		return err
	}
	vendorFileBytes, err := ioutil.ReadAll(vendorFile)
	err = json.Unmarshal(vendorFileBytes, g)
	if err != nil {
		return err
	}
	return nil
}

func (g *GovendorVendorTool) Dependencies() []PackageDepencency {
	deps := make([]PackageDepencency, len(g.Packages))
	for i, dx := range g.Packages {
		d := dx
		deps[i] = &d
	}
	return deps
}
