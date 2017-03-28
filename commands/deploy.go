package commands

import (
	"fmt"
	"context"
	"os"
	"io/ioutil"
	"os/user"
	"path"

	"github.com/urfave/cli"
	"github.com/cbegin/graven/domain"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
	"github.com/bgentry/speakeasy"
)

type ConfigMap map[string]map[string]string

var DeployCommand = cli.Command{
	Name: "deploy",
	Usage:       "Deploys artifacts to a repository",
	Action: deploy,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name: "login",
			Usage: "Prompts for repo login credentials.",
		},
	},
}

func deploy(c *cli.Context) error {
	if c.Bool("login") {
		return loginToGithub()
	}

	if err := pkg(c); err != nil {
		return err
	}

	return deployToGithub()
}

func loginToGithub() error {
	token, err := readSecret("Please type or paste a github token (will not echo): ")
	config, err := readConfig()
	if err != nil {
		config = ConfigMap{}
	}
	config["github"] = map[string]string{}
	config["github"]["token"] = token
	err = writeConfig(config)
	if err != nil {
		return fmt.Errorf("Error writing configuration file. %v", err)
	}
	return nil
}

func deployToGithub() error {
	project, err := domain.FindProject()
	if err != nil {
		return err
	}

	gh, ctx, err := authenticate()
	if err != nil {
		return err
	}

	repo, ok := project.Repositories["github"]
	if !ok {
		return fmt.Errorf("Sorry, could not find gihub repo configuration")
	}

	ownerName := repo["owner"]
	repoName := repo["repo"]

	tagName := fmt.Sprintf("v%s", project.Version)
	releaseName := tagName
	release := &github.RepositoryRelease{
		TagName: &tagName,
		Name: &releaseName,
	}

	release, _, err = gh.Repositories.CreateRelease(ctx, ownerName, repoName, release)
	if err != nil {
		return err
	}
	fmt.Printf("Created release %v/%v:%v\n", ownerName, repoName, *release.Name)

	for _, a := range project.Artifacts {
		filename := a.ArtifactFile(project)
		sourceFile, err := os.Open(project.TargetPath(filename))
		if err != nil {
			return err
		}
		opts := &github.UploadOptions{
			Name: filename,
		}
		_,_,err = gh.Repositories.UploadReleaseAsset(ctx, ownerName, repoName, *release.ID, opts, sourceFile)
		if err != nil {
			return err
		}
		fmt.Printf("Uploaded %v/%v/%v\n", ownerName, repoName, filename)
	}

	return err
}

func authenticate() (*github.Client, context.Context, error) {
	config, err := readConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("Error reading configuration (try: deploy --login): %v", err)
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config["github"]["token"]},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return client, ctx, nil
}

func readSecret(prompt string) (string, error) {
	password, err := speakeasy.Ask(prompt)
	if err != nil {
		return "", fmt.Errorf("Error reading secret from terminal: %v", err)
	}
	return password, nil
}

func readConfig() (ConfigMap, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}
	file, err := os.Open(path.Join(usr.HomeDir, ".graven.yaml"))
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	config := ConfigMap{}
	err = yaml.Unmarshal(bytes, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func writeConfig(config ConfigMap) (error) {
	usr, err := user.Current()
	if err != nil {
		return err
	}
	bytes, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path.Join(usr.HomeDir, ".graven.yaml"), bytes, 0600)
	if err != nil {
		return err
	}

	return nil
}