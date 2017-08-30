package repotool

import (
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/cbegin/graven/config"
	"github.com/cbegin/graven/domain"
	"github.com/cbegin/graven/util"
)

type MavenRepoTool struct{}

func (m *MavenRepoTool) Login(project *domain.Project, repo string) error {
	return GenericLogin(project, repo)
}

func (m *MavenRepoTool) Release(project *domain.Project, repo string) error {
	config := config.NewConfig()

	if err := config.Read(); err != nil {
		return fmt.Errorf("Error reading configuration (try: graven repo --login --name %v): %v", repo, err)
	}

	username := config.Get(project.Name, fmt.Sprintf("%v-username", repo))
	password, err := config.GetSecret(project.Name, fmt.Sprintf("%v-password", repo))
	if err != nil {
		return err
	}

	repository, ok := project.Repositories[repo]
	if !ok {
		return fmt.Errorf("Sorry, could not find repo configuration named %v", repo)
	}

	for _, a := range project.Artifacts {
		filename := a.ArtifactFile(project)
		filepath := project.TargetPath(filename)

		repoURL, err := url.Parse(repository.URL)
		groupPath := strings.Replace(repository.Group, ".", "/", -1)
		artifactId := repository.Artifact
		repoURL.Path = path.Join(repoURL.Path, groupPath, artifactId, project.Version, filename)
		if err != nil {
			return err
		}

		if err := util.UploadFile(repoURL.String(), username, password, filepath); err != nil {
			return err
		}

		fmt.Printf("Uploaded %v\n", repoURL.String())
	}
	return nil
}

func (m *MavenRepoTool) UploadDependency(project *domain.Project, repo string, dependencyFile, dependencyPath string) error {
	config := config.NewConfig()

	if err := config.Read(); err != nil {
		return fmt.Errorf("Error reading configuration (try: graven repo --login --name %v): %v", repo, err)
	}

	repository, ok := project.Repositories[repo]
	if !ok {
		return fmt.Errorf("Sorry, could not find repo configuration named %v", repo)
	}

	username := config.Get(project.Name, fmt.Sprintf("%v-username", repo))
	password, err := config.GetSecret(project.Name, fmt.Sprintf("%v-password", repo))
	if err != nil {
		return err
	}

	repoUrl, err := url.Parse(repository.URL)
	fullPath := path.Join(repoUrl.Path, dependencyPath)
	repoUrl.Path = fullPath

	if exists, err := util.HttpExists(repoUrl.String(), username, password); err != nil {
		return err
	} else if !exists {
		return util.UploadFile(repoUrl.String(), username, password, dependencyFile)
	}
	return nil
}

func (m *MavenRepoTool) DownloadDependency(project *domain.Project, repo string, dependencyFile, dependencyPath string) error {
	config := config.NewConfig()

	if err := config.Read(); err != nil {
		return fmt.Errorf("Error reading configuration (try: graven repo --login --name %v): %v", repo, err)
	}

	repository, ok := project.Repositories[repo]
	if !ok {
		return fmt.Errorf("Sorry, could not find repo configuration named %v", repo)
	}

	username := config.Get(project.Name, fmt.Sprintf("%v-username", repo))
	password, err := config.GetSecret(project.Name, fmt.Sprintf("%v-password", repo))
	if err != nil {
		return err
	}

	repoUrl, err := url.Parse(repository.URL)
	fullPath := path.Join(repoUrl.Path, dependencyPath)
	repoUrl.Path = fullPath

	if exists, err := util.HttpExists(repoUrl.String(), username, password); err != nil {
		return err
	} else if !exists {
		return fmt.Errorf("%v doesn't exist", repoUrl.String())
	}
	return util.DownloadFile(repoUrl.String(), username, password, dependencyFile)
}
