package vendortool

import (
	"strings"
	"time"
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"

	"github.com/cbegin/graven/domain"
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

func (g *GovendorPackage) ArchiveFileName() string {
	chars := []string{"/","."}
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
	return g.Path
}

func (g *GovendorVendorTool) LoadFile(project *domain.Project) error {
	vendorFilePath := project.ProjectPath("vendor", "vendor.json")
	vendorFile, err := os.Open(vendorFilePath)
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
