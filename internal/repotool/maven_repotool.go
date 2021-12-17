package repotool

import (
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/cbegin/graven/internal/config"
	"github.com/cbegin/graven/internal/domain"
	"github.com/cbegin/graven/internal/util"
)

type MavenRepoTool struct{}

func (m *MavenRepoTool) Login(project *domain.Project, repo string) error {
	return GenericLogin(project, repo)
}

func (m *MavenRepoTool) LoginTest(*domain.Project, string) error {
	return nil
}

func (m *MavenRepoTool) Release(project *domain.Project, repo string) error {
	c := config.NewConfig()

	if err := c.Read(); err != nil {
		return fmt.Errorf("Error reading configuration (try: graven repo login --name %v): %v", repo, err)
	}

	username := c.Get(project.Name, fmt.Sprintf("%v-username", repo))
	password, err := c.GetSecret(project.Name, fmt.Sprintf("%v-password", repo))
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
