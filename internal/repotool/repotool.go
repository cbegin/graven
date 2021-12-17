package repotool

import (
	"fmt"

	"github.com/cbegin/graven/internal/config"
	"github.com/cbegin/graven/internal/domain"
)

var RepoRegistry = map[string]RepoTool{}

func init() {
	RepoRegistry["github"] = &GithubRepoTool{}
	RepoRegistry["maven"] = &MavenRepoTool{}
	RepoRegistry["docker"] = &DockerRepotool{}
}

type RepoTool interface {
	Login(project *domain.Project, repo string) error
	LoginTest(project *domain.Project, repo string) error
	Release(project *domain.Project, repo string) error
}

func GenericLogin(project *domain.Project, repo string) error {
	c := config.NewConfig()
	if err := c.Read(); err != nil {
		// ignore
	}
	if err := c.PromptPlainText(project.Name, fmt.Sprintf("%v-username", repo), "Username: "); err != nil {

	}
	if err := c.PromptSecret(project.Name, fmt.Sprintf("%v-password", repo), "Password: "); err != nil {

	}
	if err := c.Write(); err != nil {
		return fmt.Errorf("Error writing configuration file. %v", err)
	}
	return nil
}
