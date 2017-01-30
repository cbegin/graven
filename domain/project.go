package domain

import (
	"os"
	"gopkg.in/yaml.v2"
	"fmt"
	"path"
)

const (
	ProjectFileName = "project.yaml"
)

type Project struct {
	FilePath string
	Name  string `yaml:"name"`
	Version  string `yaml:"version"`
	Artifacts []Artifact `yaml:"artifacts"`
	Repositories []Repository `yaml:"repositories"`
}

type Artifact struct {
	Classifier string `yaml:"classifier"`
	Resources []string `yaml:"resources"`
	Targets []Target `yaml:"targets"`
	Archive string `yaml:"archive"`
}

type Target struct {
	Executable string `yaml:"executable"`
	Package string `yaml:"package"`
	Flags string `yaml:"flags"`
	Environment map[string]string `yaml:"env"`
}

type Repository struct {
	Name string `yaml:"name"`
	URL string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Token string `yaml:"token"`
	Type string `yaml:"type"`
}

func (p *Project) TargetPath(subdirs ...string) string {
	target := path.Join("target")
	for _, s := range subdirs {
		target = path.Join(target, s)
	}
	return target
}

func FindProject() (*Project, error) {
	wd, _ := os.Getwd()
	cwd := wd
	for len(cwd) > 1 {
		p := path.Join(cwd, ProjectFileName)
		_, err := os.Stat(p)
		if !os.IsNotExist(err) {
			return LoadProject(p)
		}
		cwd = path.Dir(cwd)
	}
	return nil, fmt.Errorf("Could not find project file in (or in any parent of) %v", wd)
}

func LoadProject(filepath string) (*Project, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	buffer := make([]byte, info.Size())

	n, err := file.Read(buffer)
	if err != nil {
		return nil, err
	}
	if int64(n) != info.Size() {
		return nil, fmt.Errorf("Expected to read %v but read %v.", info.Size(), n)
	}

	project := &Project{}
	err = yaml.Unmarshal(buffer, project)
	if err != nil {
		return nil, err
	}

	project.FilePath = filepath
	return project, nil
}
