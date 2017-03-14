package commands

import (
	"fmt"
	"context"
	"os"
	"io/ioutil"
	"os/user"
	"path"

	"github.com/urfave/cli"
	"github.com/cbegin/graven/domain"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
)

var DeployCommand = cli.Command{
	Name: "deploy",
	Usage:       "Deploys artifacts to a repository",
	Action: deploy,
}

func deploy(c *cli.Context) error {
	if err := pkg(c); err != nil {
		return err
	}

	project, err := domain.FindProject()
	if err != nil {
		return err
	}

	gh, ctx, err := authenticate()
	if err != nil {
		return err
	}

	repo, ok := project.Repositories["github"]
	if !ok {
		return fmt.Errorf("Sorry, could not find gihub repo configuration")
	}

	ownerName := repo["owner"]
	repoName := repo["repo"]

	tagName := fmt.Sprintf("v%s", project.Version)
	releaseName := tagName
	release := &github.RepositoryRelease{
		TagName: &tagName,
		Name: &releaseName,
	}

	release, _, err = gh.Repositories.CreateRelease(ctx, ownerName, repoName, release)
	if err != nil {
		return err
	}

	for _, a := range project.Artifacts {
		filename := a.ArtifactFile(project)
		sourceFile, err := os.Open(project.TargetPath(filename))
		if err != nil {
			return err
		}
		opts := &github.UploadOptions{
			Name: filename,
		}
		_,_,err = gh.Repositories.UploadReleaseAsset(ctx, ownerName, repoName, *release.ID, opts, sourceFile)
		if err != nil {
			return err
		}
	}

	return err
}

func authenticate() (*github.Client, context.Context, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, nil, err
	}
	file, err := os.Open(path.Join(usr.HomeDir, ".graven.yaml"))
	if err != nil {
		return nil, nil, err
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, nil, err
	}
	config := map[string]map[string]string{}
	err = yaml.Unmarshal(bytes, config)
	if err != nil {
		return nil, nil, err
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config["github"]["token"]},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return client, ctx, nil
}