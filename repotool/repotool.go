package repotool

import "github.com/cbegin/graven/domain"

type RepoTool interface {
	Release(project *domain.Project) error
	Login(project *domain.Project) error
}
