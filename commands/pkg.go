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
)

var PackageCommand = cli.Command{
	Name: "package",
	Usage:       "project package",
	UsageText:   "package - project package",
	Description: "produces packaged archive",
	Action: pkg,
}

func pkg(c *cli.Context) error {
	project := c.App.Metadata["project"].(*domain.Project)

	clean(c)
	build(c)

	for _, artifact := range project.Artifacts {
		nameParts := strings.Split(project.Name, "/")
		shortName := nameParts[len(nameParts) - 1:][0]
		targetFile := fmt.Sprintf("%s-%s-%s.%s", shortName, project.Version, artifact.Classifier, artifact.Archive)
		switch artifact.Archive {
		case "zip":
			return zipDir(project.TargetPath(artifact.Classifier), project.TargetPath(targetFile))
		case "tar.gz":

		}

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