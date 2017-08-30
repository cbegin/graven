package repotool

import "github.com/cbegin/graven/domain"

var RepoRegistry = map[string]RepoTool{}

func init() {
	RepoRegistry["github"] = &GithubRepoTool{}
	RepoRegistry["maven"] = &MavenRepoTool{}
}

type RepoTool interface {
	Release(project *domain.Project, repo string) error
	Login(project *domain.Project, repo string) error
	UploadDependency(project *domain.Project, repo string, dependencyFile, dependencyPath string) error
	DownloadDependency(project *domain.Project, repo string, dependencyFile, dependencyPath string) error
}
