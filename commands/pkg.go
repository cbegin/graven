package commands

import (
	"github.com/urfave/cli"
	"github.com/cbegin/graven/domain"
	"archive/zip"
	"path/filepath"
	"os"
	"strings"
	"io"
	"fmt"
	"archive/tar"
	"compress/gzip"
	"sync"
	"github.com/hashicorp/go-multierror"
)

var PackageCommand = cli.Command{
	Name: "package",
	Usage:       "project package",
	UsageText:   "package - project package",
	Description: "produces packaged archive",
	Action: pkg,
}

func pkg(c *cli.Context) error {
	if err := clean(c); err != nil {
		return err
	}

	if err := build(c); err != nil {
		return err
	}

	var merr error
	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}
	project := c.App.Metadata["project"].(*domain.Project)
	for _, artifact := range project.Artifacts {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := packageArtifact(project, &artifact)
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
	nameParts := strings.Split(project.Name, "/")
	shortName := nameParts[len(nameParts) - 1:][0]
	targetFile := fmt.Sprintf("%s-%s-%s.%s", shortName, project.Version, artifact.Classifier, artifact.Archive)
	source := project.TargetPath(artifact.Classifier)
	dest := project.TargetPath(targetFile)
	temp := fmt.Sprintf("%s.tmp", dest)
	switch artifact.Archive {
	case "zip":
		err := zipDir(source, dest)
		if err != nil {
			return err
		}
	case "tgz":
		fallthrough
	case "tar.gz":
		err := tarDir(source, dest)
		if err != nil {
			return err
		}
		err = gzipFile(dest, temp)
		if err != nil {
			return err
		}
		os.Remove(dest)
		os.Rename(temp, dest)
	}
	return nil
}

func zipDir(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path != source {

			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}

			header.Name = strings.TrimPrefix(strings.TrimPrefix(path, source), "/")

			if info.IsDir() {
				header.Name += "/"
			} else {
				header.Method = zip.Deflate
			}

			writer, err := archive.CreateHeader(header)
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
			return err
		}
		return err
	})

	return err
}

func tarDir(source, target string) error {
	tarfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer tarfile.Close()

	tarball := tar.NewWriter(tarfile)
	defer tarball.Close()

	return filepath.Walk(source,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if path != source {

				header, err := tar.FileInfoHeader(info, info.Name())
				if err != nil {
					return err
				}

				header.Name = strings.TrimPrefix(strings.TrimPrefix(path, source), "/")

				if err := tarball.WriteHeader(header); err != nil {
					return err
				}

				if info.IsDir() {
					return nil
				}

				file, err := os.Open(path)
				if err != nil {
					return err
				}
				defer file.Close()
				_, err = io.Copy(tarball, file)
				return err
			}
			return err
		})
}

func gzipFile(source, target string) error {
	reader, err := os.Open(source)
	if err != nil {
		return err
	}

	writer, err := os.Create(target)
	if err != nil {
		return err
	}
	defer writer.Close()

	archiver := gzip.NewWriter(writer)
	//archiver.Name = filename
	defer archiver.Close()

	_, err = io.Copy(archiver, reader)
	return err
}

