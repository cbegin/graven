package repotool

import (
	"fmt"

	"github.com/cbegin/graven/config"
	"github.com/cbegin/graven/domain"
)

type MavenRepoTool struct{}

func (m *MavenRepoTool) Login(project *domain.Project, repo string) error {
	config := config.NewConfig()
	if err := config.Read(); err != nil {
		// ignore
	}
	if err := config.SetPlainText(project.Name, fmt.Sprintf("%v-username", repo), "Username: "); err != nil {

	}
	if err := config.SetSecret(project.Name, fmt.Sprintf("%v-password", repo), "Password: "); err != nil {

	}
	if err := config.Write(); err != nil {
		return fmt.Errorf("Error writing configuration file. %v", err)
	}
	return nil
}

func (m *MavenRepoTool) Release(project *domain.Project, repo string) error {

	return nil
}

