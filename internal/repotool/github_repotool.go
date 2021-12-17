package repotool

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/cbegin/graven/internal/config"
	"github.com/cbegin/graven/internal/domain"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type GithubRepoTool struct{}

func (g *GithubRepoTool) Login(project *domain.Project, repo string, auth string) error {
	c := config.NewConfig()
	if err := c.Read(); err != nil {
		// ignore
	}
	if auth == "" {
		err := c.PromptSecret(project.Name, repo, "Please type or paste a github token (will not echo): ")
		if err != nil {
			return fmt.Errorf("Error reading token from stdin: %v", err)
		}
	} else {
		if err := c.SetSecret(project.Name, repo, auth); err != nil {
			return fmt.Errorf("Error writing secret to config: %v", err)
		}
	}
	if err := c.Write(); err != nil {
		return fmt.Errorf("Error writing configuration file: %v", err)
	}
	return nil
}

func (g *GithubRepoTool) LoginTest(project *domain.Project, repo string) error {
	gh, ctx, err := authenticate(project, repo)
	if err != nil {
		return fmt.Errorf("Authentication error for repo %v: %v", repo, err)
	}

	repository, found := project.Repositories[repo]
	if !found {
		return fmt.Errorf("No repo named %v is found in project.", repo)
	}

	_, _, err = gh.Repositories.ListReleases(ctx, repository.Group, repository.Artifact, &github.ListOptions{})
	if err != nil {
		return fmt.Errorf("Error listing releases for %v: %v", repo, err)
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

func authenticate(project *domain.Project, repo string) (*github.Client, context.Context, error) {
	c := config.NewConfig()

	if err := c.Read(); err != nil {
		return nil, nil, fmt.Errorf("Error reading configuration (try: graven repo login --name %v): %v", repo, err)
	}

	token, err := c.GetSecret(project.Name, repo)
	if err != nil {
		return nil, nil, err
	}
	if token == "" {
		return nil, nil, fmt.Errorf("Configuration missing token (try: graven repo login --name %v).", repo)
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
