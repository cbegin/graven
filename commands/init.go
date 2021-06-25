package commands

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/cbegin/graven/domain"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

var InitCommand = cli.Command{
	Name:   "init",
	Usage:  "Initializes a project directory",
	Action: initialize,
}

type ClassifierTemplate struct {
	Classifier   string
	Archive      string
	Extension    string
	OS           string
	Architecture string
}

type PackagePath struct {
	Package string
	Path    string
}

var (
	darwinTemplate = ClassifierTemplate{
		Classifier:   "darwin",
		Archive:      "tgz",
		Extension:    "",
		OS:           "darwin",
		Architecture: "amd64",
	}
	linuxTemplate = ClassifierTemplate{
		Classifier:   "linux",
		Archive:      "tar.gz",
		Extension:    "",
		OS:           "linux",
		Architecture: "amd64",
	}
	winTemplate = ClassifierTemplate{
		Classifier:   "win",
		Archive:      "zip",
		Extension:    ".exe",
		OS:           "windows",
		Architecture: "amd64",
	}
	templates = []ClassifierTemplate{
		darwinTemplate,
		linuxTemplate,
		winTemplate,
	}
)

var mainTemplate = `package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello!")
}`

func initialize(c *cli.Context) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	projectPath := filepath.Join(wd, "project.yaml")
	if _, err := os.Stat(projectPath); !os.IsNotExist(err) {
		return fmt.Errorf("%v already exists. No changes made.", projectPath)
	}

	packages := &[]PackagePath{}

	if err = filepath.Walk(wd, getInitializeWalkerFunc(wd, packages)); err != nil {
		return err
	}

	if len(*packages) < 1 {
		err := ioutil.WriteFile("main.go", []byte(mainTemplate), 0655)
		if err != nil {
			return err
		}
		*packages = append(*packages, PackagePath{
			Package: "main",
			Path:    "",
		})
	}

	err = writeVersionFile(&domain.Project{
		FilePath: projectPath,
		Version:  "0.0.1-DEV",
	})
	if err != nil {
		return err
	}

	artifacts := []domain.Artifact{}

	for _, template := range templates {
		targets := []domain.Target{}
		for _, p := range *packages {
			if p.Package == "main" {
				pkg := fmt.Sprintf(".%v", p.Path)
				executable := filepath.Base(filepath.Join(wd, pkg))
				targets = append(targets, domain.Target{
					Executable: fmt.Sprintf("%v%v", executable, template.Extension),
					Package:    pkg,
					Flags:      []string{},
				})
			}
		}

		artifacts = append(artifacts, domain.Artifact{
			Classifier: template.Classifier,
			Resources:  []string{},
			Archive:    template.Archive,
			Targets:    targets,
			Environment: map[string]string{
				"GOOS":   template.OS,
				"GOARCH": template.Architecture,
			},
		})
	}

	newProject := &domain.Project{}
	newProject.Name = filepath.Base(wd)
	newProject.Version = "0.0.1"
	newProject.Artifacts = artifacts

	bytes, err := yaml.Marshal(newProject)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(projectPath, bytes, 0655); err != nil {
		return err
	}

	return err
}

func getInitializeWalkerFunc(basePath string, packages *[]PackagePath) filepath.WalkFunc {
	fs := token.NewFileSet()
	return func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			subDir := path[len(basePath):]
			subDirParts := strings.Split(subDir, string(filepath.Separator))
			matches, _ := filepath.Glob(filepath.Join(path, "*.go"))
			if len(matches) > 0 && !contains(subDirParts, map[string]struct{}{
				"vendor": struct{}{},
				"target": struct{}{},
				".git":   struct{}{}}) {
				ast, err := parser.ParseDir(fs, path, nil, parser.PackageClauseOnly)
				if err != nil {
					fmt.Println(err)
					return err
				}
				for _, v := range ast {
					shortPath := path[len(basePath):]
					*packages = append(*packages, PackagePath{
						Package: v.Name,
						Path:    shortPath,
					})
				}
			}
		}
		return nil
	}
}
