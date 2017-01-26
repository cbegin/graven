package main

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
	Package  string `yaml:"package"`
	Version  string `yaml:"version"`
}

func FindProject() (*Project, error) {
	wd, _ := os.Getwd()

	fmt.Println(wd)

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
