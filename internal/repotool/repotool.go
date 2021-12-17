package repotool

import (
	"fmt"
	"strings"

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
	Login(project *domain.Project, repo string, auth string) error
	LoginTest(project *domain.Project, repo string) error
	Release(project *domain.Project, repo string) error
}

func GenericLogin(project *domain.Project, repo, auth string) error {
	c := config.NewConfig()
	if err := c.Read(); err != nil {
		// ignore
	}

	usernameField := fmt.Sprintf("%v-username", repo)
	passwordField := fmt.Sprintf("%v-password", repo)

	if auth == "" {
		if err := c.PromptPlainText(project.Name, usernameField, "Username: "); err != nil {
			return fmt.Errorf("Error reading username from stdin. %v", err)
		}
		if err := c.PromptSecret(project.Name, passwordField, "Password: "); err != nil {
			return fmt.Errorf("Error reading password from stdin. %v", err)
		}
	} else {
		userPass := strings.SplitN(auth, ":", 2)
		if len(userPass) != 2 {
			return fmt.Errorf("Format of provided auth string must be username:password")
		}
		c.Set(project.Name, usernameField, userPass[0])
		if err := c.SetSecret(project.Name, passwordField, userPass[1]); err != nil {
			return fmt.Errorf("Error writing secret to config: %v", err)
		}
	}

	if err := c.Write(); err != nil {
		return fmt.Errorf("Error writing configuration file. %v", err)
	}
	return nil
}
