package domain

import (
	"strings"
	"time"
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
)

type GovendorPackage struct {
	Origin       string `json:"oriin,omitempty"`
	Path         string `json:"path,omitempty"`
	Revision     string `json:"revision,omitempty"`
	RevisionTime string `json:"revisionTime,omitempty"`
	Comment      string `json:"comment,omitempty"`
}

type GovendorFile struct {
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

func ReadGovendorFile (project *Project) (*GovendorFile, error) {
	govendorFile := &GovendorFile{}

	vendorFilePath := project.ProjectPath("vendor", "vendor.json")
	vendorFile, err := os.Open(vendorFilePath)
	vendorFileBytes, err := ioutil.ReadAll(vendorFile)
	err = json.Unmarshal(vendorFileBytes, govendorFile)
	if err != nil {
		return nil, err
	}
	return govendorFile, nil
}
