package domain

import (
	"fmt"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

const (
	ProjectFileName = "project.yaml"
	TargetDirName   = "target"
)

type Project struct {
	FilePath     string       `yaml:",omitempty"`
	Name         string       `yaml:"name"`
	Version      string       `yaml:"version"`
	Artifacts    []Artifact   `yaml:"artifacts"`
	Repositories map[string]map[string]string `yaml:"repositories"`
	Resources    []string     `yaml:"resources"`
}

type Artifact struct {
	Classifier  string            `yaml:"classifier"`
	Targets     []Target          `yaml:"targets"`
	Archive     string            `yaml:"archive"`
	Resources   []string          `yaml:"resources"`
	Environment map[string]string `yaml:"env"`
}

type Target struct {
	Executable  string            `yaml:"executable"`
	Package     string            `yaml:"package"`
	Flags       string            `yaml:"flags"`
	Environment map[string]string `yaml:"env"`
}

type Repository struct {
	Name     string `yaml:"name"`
	URL      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Token    string `yaml:"token"`
	Type     string `yaml:"type"`
}

func (p *Project) TargetPath(subdirs ...string) string {
	targetPath := p.ProjectPath(TargetDirName)
	for _, s := range subdirs {
		targetPath = path.Join(targetPath, s)
	}
	return targetPath
}

func (p *Project) ProjectPath(subdirs ...string) string {
	projectPath := path.Dir(p.FilePath)
	for _, s := range subdirs {
		projectPath = path.Join(projectPath, s)
	}
	return projectPath
}

func (a *Artifact) ArtifactFile(project *Project) string {
	return fmt.Sprintf("%s-%s-%s.%s", project.Name, project.Version, a.Classifier, a.Archive)
}

var FindProject = internalFindProject
func internalFindProject() (*Project, error) {
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
