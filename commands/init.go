package commands

import (
	"go/parser"
	"go/token"
	"path/filepath"
	"os"
	"strings"
	"fmt"
	"io/ioutil"

	"github.com/urfave/cli"
	"github.com/cbegin/graven/domain"
	"gopkg.in/yaml.v2"
	"path"
)

var InitCommand = cli.Command{
	Name: "init",
	Usage:       "Initializes a project directory",
	Action: initialize,
}

func initialize(c *cli.Context) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	projectPath := path.Join(wd, "project.yaml")

	if err = filepath.Walk(wd, getInitializeWalkerFunc(wd)); err != nil {
		return err
	}

	newProject := &domain.Project{}
	newProject.Name = "github.com/org/myProject"
	newProject.Version = "0.0.1"
	newProject.Artifacts = []domain.Artifact{
		domain.Artifact{

		},
	}

	bytes, err := yaml.Marshal(newProject)
	if err != nil {
		return err
	}
	if _, err := os.Stat(projectPath); !os.IsNotExist(err) {
		return fmt.Errorf("%v already exists. No changes made.", projectPath)
	}

	if err := ioutil.WriteFile(projectPath, bytes, 0655); err != nil {
		return err
	}

	return err
}

func getInitializeWalkerFunc(basePath string) filepath.WalkFunc {
	fs := token.NewFileSet()
	return func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			subDir := path[len(basePath):]
			subDirParts := strings.Split(subDir, string(filepath.Separator))
			matches, _ := filepath.Glob(filepath.Join(path, "*.go"));
			if len(matches) > 0 && !contains(subDirParts, map[string]struct{}{
				"vendor":struct{}{},
				"target":struct{}{},
				".git":struct{}{}}) {
				ast, err := parser.ParseDir(fs, path, nil, parser.PackageClauseOnly)
				if err != nil {
					fmt.Println(err)
					return err
				}
				for _, v := range ast {
					shortPath := path[len(basePath):]
					fmt.Printf("%v => %v\n", v.Name, shortPath)
				}
			}
		}
		return nil
	}
}
