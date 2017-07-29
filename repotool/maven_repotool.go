package repotool

import (
	"fmt"
	"os"
	"net/http"
	"bytes"
	"net/url"
	"path"
	"strings"

	"github.com/cbegin/graven/config"
	"github.com/cbegin/graven/domain"
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

		if err := uploadFile(repoURL.String(), username, password, filepath); err != nil {
			return err
		}

		fmt.Printf("Uploaded %v\n", repoURL.String())
	}
	return nil
}

func uploadFile(uri, username, password, filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	req, err := http.NewRequest("PUT", uri, file)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.SetBasicAuth(username, password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body := &bytes.Buffer{}

	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("Error [%v]: %v", resp.StatusCode, body)
	}

	return nil
}

