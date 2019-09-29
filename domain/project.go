package domain

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	ProjectFileName = "project.yaml"
	TargetDirName   = "target"

	RepositoryRoleRelease    = "release"
	RepositoryRoleDependency = "dependency"
)

type Project struct {
	FilePath     string                `yaml:",omitempty"`
	Name         string                `yaml:"name"`
	Version      string                `yaml:"version"`
	GoVersion    string                `yaml:"go_version,omitempty"`
	Artifacts    []Artifact            `yaml:"artifacts,omitempty"`
	Repositories map[string]Repository `yaml:"repositories,omitempty"`
	Resources    []string              `yaml:"resources,omitempty"`
}

type Artifact struct {
	Classifier  string            `yaml:"classifier,omitempty"`
	Targets     []Target          `yaml:"targets,omitempty"`
	Archive     string            `yaml:"archive,omitempty"`
	Resources   []string          `yaml:"resources,omitempty"`
	Environment map[string]string `yaml:"env,omitempty"`
}

type Target struct {
	Executable  string            `yaml:"executable,omitempty"`
	Package     string            `yaml:"package,omitempty"`
	Flags       string            `yaml:"flags,omitempty"`
	Environment map[string]string `yaml:"env,omitempty"`
}

type Repository struct {
	URL      string   `yaml:"url,omitempty"`
	Group    string   `yaml:"group,omitempty"`
	Artifact string   `yaml:"artifact,omitempty"`
	File     string   `yaml:"file,omitempty"`
	Type     string   `yaml:"type,omitempty"`
	Roles    []string `yaml:"roles,omitempty"`
}

func (p *Project) TargetPath(subdirs ...string) string {
	targetPath := p.ProjectPath(TargetDirName)
	for _, s := range subdirs {
		targetPath = filepath.Join(targetPath, s)
	}
	return targetPath
}

func (p *Project) ProjectPath(subdirs ...string) string {
	projectPath := filepath.Dir(p.FilePath)
	for _, s := range subdirs {
		projectPath = filepath.Join(projectPath, s)
	}
	return projectPath
}

func (a *Artifact) ArtifactFile(project *Project) string {
	return fmt.Sprintf("%s-%s-%s.%s", project.Name, project.Version, a.Classifier, a.Archive)
}

func (r *Repository) HasRole(role string) bool {
	for _, r := range r.Roles {
		if r == role {
			return true
		}
	}
	return false
}

var FindProject = internalFindProject

func internalFindProject() (*Project, error) {
	wd, _ := os.Getwd()
	cwd := wd
	for len(cwd) > 1 {
		p := filepath.Join(cwd, ProjectFileName)
		_, err := os.Stat(p)
		if !os.IsNotExist(err) {
			return LoadProject(p)
		}
		cwd = filepath.Dir(cwd)
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

func (project *Project) Environment() []string {
	// This function used to do more with the existing environment,
	// but now we just forward the whole environment to avoid surprising
	// behavior. We'll keep it in case we ever want to add the ability to
	// unset existing env vars at some point.
	return os.Environ()
}
