package repotool

import (
	"fmt"

	"github.com/cbegin/graven/config"
	"github.com/cbegin/graven/domain"
	"github.com/cbegin/graven/util"
	"path"
)

type DockerRepotool struct{}

func (r *DockerRepotool) Login(project *domain.Project, repo string) error {
	return GenericLogin(project, repo)
}

func (r *DockerRepotool) LoginTest(project *domain.Project, repo string) error {
	return nil
}

func (r *DockerRepotool) Release(project *domain.Project, repo string) error {
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

	if username == "" || password == "" {
		return fmt.Errorf("Could not find docker credentials. Please log in with: graven repo --login --name [reponame]")
	}
	if sout, serr, err := util.RunCommand(project.ProjectPath(), nil, "docker", "login", "-u", username, "-p", password, repository.URL); err != nil {
		fmt.Printf("Logging into Docker...  %v\n%v\n", sout, serr)
		return err
	}

	dockerPath := path.Join(repository.URL, repository.Group, repository.Artifact)
	dockerTag := fmt.Sprintf("%v:%v", dockerPath, project.Version)

	fmt.Printf("Pushing docker image %v\n", dockerTag)
	if sout, serr, err := util.RunCommand(project.ProjectPath(), nil, "docker", "push", dockerTag); err != nil {
		fmt.Printf("Running Docker build...  %v\n%v\n", sout, serr)
		return err
	}

	return nil
}
