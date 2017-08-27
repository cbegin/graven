package repotool

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/cbegin/graven/config"
	"github.com/cbegin/graven/domain"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type GithubRepoTool struct{}

func (g *GithubRepoTool) Login(project *domain.Project, repo string) error {
	config := config.NewConfig()
	if err := config.Read(); err != nil {
		// ignore
	}
	err := config.PromptSecret(project.Name, repo, "Please type or paste a github token (will not echo): ")
	err = config.Write()
	if err != nil {
		return fmt.Errorf("Error writing configuration file. %v", err)
	}
	return nil
}

func (g *GithubRepoTool) Release(project *domain.Project, repo string) error {
	gh, ctx, err := authenticate(project, repo)
	if err != nil {
		return err
	}

	repository, ok := project.Repositories[repo]
	if !ok {
		return fmt.Errorf("Sorry, could not find repo configuration named %v", repo)
	}

	ownerName := repository.Group
	repoName := repository.Artifact

	tagName := fmt.Sprintf("v%s", project.Version)
	releaseName := tagName
	release := &github.RepositoryRelease{
		TagName: &tagName,
		Name:    &releaseName,
	}

	release, _, err = gh.Repositories.CreateRelease(ctx, ownerName, repoName, release)
	if err != nil {
		return err
	}
	fmt.Printf("Created release %v/%v (%v):%v\n", ownerName, repoName, release.GetID(), release.GetName())

	for _, a := range project.Artifacts {
		filename := a.ArtifactFile(project)
		sourceFile, err := os.Open(project.TargetPath(filename))
		if err != nil {
			return err
		}
		opts := &github.UploadOptions{
			Name: filename,
		}

		_, _, err = gh.Repositories.UploadReleaseAsset(ctx, ownerName, repoName, release.GetID(), opts, sourceFile)
		if err != nil {
			return err
		}
		fmt.Printf("Uploaded %v/%v/%v\n", ownerName, repoName, filename)
	}

	return err
}

func (g *GithubRepoTool) UploadDependency(project *domain.Project, repo string, dependencyFile, dependencyPath string) error {
	return fmt.Errorf("Github repos don't support dependencies.")
}

func authenticate(project *domain.Project, repo string) (*github.Client, context.Context, error) {
	config := config.NewConfig()

	if err := config.Read(); err != nil {
		return nil, nil, fmt.Errorf("Error reading configuration (try: release --login): %v", err)
	}

	token, err := config.GetSecret(project.Name, repo)
	if err != nil {
		return nil, nil, err
	}
	if token == "" {
		return nil, nil, fmt.Errorf("Configuration missing token (try: release --login).")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	if repository, hasRepo := project.Repositories[repo]; hasRepo {
		if baseURL := repository.URL; baseURL != "" {
			apiUrl := fmt.Sprintf("%v/api/v3/", baseURL)
			uploadUrl := fmt.Sprintf("%v/api/uploads/", baseURL)
			if u, err := url.ParseRequestURI(apiUrl); err != nil {
				return nil, nil, fmt.Errorf("Error parsing repo URL : %v. Cause: %v", u, err)
			} else {
				client.BaseURL = u
			}
			if u, err := url.ParseRequestURI(uploadUrl); err != nil {
				return nil, nil, fmt.Errorf("Error parsing repo URL : %v. Cause: %v", u, err)
			} else {
				client.UploadURL = u
			}

		}
	}

	return client, ctx, nil
}
