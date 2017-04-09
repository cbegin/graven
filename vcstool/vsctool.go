package vcs

import "github.com/cbegin/graven/domain"

type VCSTool interface {
	VerifyRepoState(project *domain.Project) error
	Tag(tagname string) error
}
