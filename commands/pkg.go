package commands

import (
	"fmt"
	"os"
	"sync"

	"github.com/cbegin/graven/domain"
	"github.com/cbegin/graven/util"
	"github.com/hashicorp/go-multierror"
	"github.com/urfave/cli"
	"path"
	"path/filepath"
)

var PackageCommand = cli.Command{
	Name:   "package",
	Usage:  "Produces packaged archive for deployment",
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

	for _, repo := range project.Repositories {
		if repo.Type == "docker" {
			dockerPath := path.Join(repo.URL, repo.Group, repo.Artifact)
			dockerTag := fmt.Sprintf("%v:%v", dockerPath, project.Version)
			dockerDir := filepath.Dir(project.ProjectPath(repo.File))
			fmt.Printf("Building docker image %v\n", dockerTag)
			if sout, serr, err := util.RunCommand(project.ProjectPath(), nil, "docker", "build", "-f", repo.File, "-t", dockerTag, dockerDir); err != nil {
				fmt.Printf("%v\n%v\n", sout, serr)
				merr = multierror.Append(merr, err)
			}
		}
	}

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
		err := util.ZipDir(source, dest, true)
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
