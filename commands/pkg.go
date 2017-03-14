package commands

import (
	"os"
	"fmt"
	"sync"

	"github.com/urfave/cli"
	"github.com/cbegin/graven/domain"
	"github.com/hashicorp/go-multierror"
	"github.com/cbegin/graven/util"
)

var PackageCommand = cli.Command{
	Name: "package",
	Usage:       "Produces packaged archive for deployment",
	Action: pkg,
}

func pkg(c *cli.Context) error {
	if err := build(c); err != nil {
		return err
	}

	project, err := domain.FindProject()
	if err != nil {
		return err
	}

	var merr error
	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}
	for _, artifact := range project.Artifacts {
		a := artifact
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := packageArtifact(project, &a)
			if err != nil {
				mutex.Lock()
				merr = multierror.Append(merr, err)
				mutex.Unlock()
			}
		}()
	}
	wg.Wait()
	return merr
}

func packageArtifact(project *domain.Project, artifact *domain.Artifact) error {
	targetFile := artifact.ArtifactFile(project)
	source := project.TargetPath(artifact.Classifier)
	dest := project.TargetPath(targetFile)
	temp := fmt.Sprintf("%s.tmp", dest)
	fmt.Printf("Packaging %v\n", targetFile)
	defer fmt.Printf("Done %v\n", targetFile)
	switch artifact.Archive {
	case "zip":
		err := util.ZipDir(source, dest)
		if err != nil {
			return err
		}
	case "tgz":
		fallthrough
	case "tar.gz":
		err := util.TarDir(source, dest)
		if err != nil {
			return err
		}
		err = util.GzipFile(dest, temp)
		if err != nil {
			return err
		}
		os.Remove(dest)
		os.Rename(temp, dest)
	}
	return nil
}
