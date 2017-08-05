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
	config := config.NewConfig()
	if err := config.Read(); err != nil {
		// ignore
	}
	if err := config.PromptPlainText(project.Name, fmt.Sprintf("%v-username", repo), "Username: "); err != nil {

	}
	if err := config.PromptSecret(project.Name, fmt.Sprintf("%v-password", repo), "Password: "); err != nil {

	}
	if err := config.Write(); err != nil {
		return fmt.Errorf("Error writing configuration file. %v", err)
	}
	return nil
}

func (m *MavenRepoTool) Release(project *domain.Project, repo string) error {
	config := config.NewConfig()

	if err := config.Read(); err != nil {
		return fmt.Errorf("Error reading configuration (try: release --login): %v", err)
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
		groupPath := strings.Replace(repository.GroupID, ".", "/", -1)
		artifactId := repository.ArtifactID
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
