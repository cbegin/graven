package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/cbegin/graven/buildtool"
	"github.com/cbegin/graven/domain"
	"github.com/hashicorp/go-multierror"
	"github.com/urfave/cli"
)

var TestCommand = cli.Command{
	Name:   "test",
	Usage:  "Finds and runs tests in this project",
	Action: tester,
}

func tester(c *cli.Context) error {
	if err := clean(c); err != nil {
		return err
	}

	project, err := domain.FindProject()
	if err != nil {
		return err
	}

	err = os.MkdirAll(project.TargetPath("reports"), 0755)
	if err != nil {
		return fmt.Errorf("Could not create reports directory. %v", err)
	}

	waitGroup := &sync.WaitGroup{}
	var merr error
	if err := filepath.Walk(project.ProjectPath(), getTestWalkerFunc(project, waitGroup, &merr)); err != nil {
		merr = multierror.Append(merr, err)
	}
	waitGroup.Wait()
	return merr
}

func getTestWalkerFunc(project *domain.Project, waitGroup *sync.WaitGroup, merr *error) filepath.WalkFunc {
	// TODO - Make this configurable
	var buildTool builder.BuildTool = &builder.GoBuildTool{}
	return func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			subDir := path[len(project.ProjectPath()):]
			subDirParts := strings.Split(subDir, string(filepath.Separator))
			matches, _ := filepath.Glob(filepath.Join(path, "*_test.go"))
			if len(matches) > 0 && !contains(subDirParts, map[string]struct{}{
				"vendor": struct{}{},
				"target": struct{}{},
				".git":   struct{}{}}) {
				fmt.Printf("Testing %v\n", subDir)
				waitGroup.Add(1)
				go func() {
					if err := buildTool.Test(subDir, project); err != nil {
						*merr = multierror.Append(*merr, err)
					}
					waitGroup.Done()
				}()
			}
		}
		return nil
	}
}

func contains(strings []string, exclusions map[string]struct{}) bool {
	for _, a := range strings {
		if _, ok := exclusions[a]; ok {
			return true
		}
	}
	return false
}
