package vcstool

import (
	"github.com/cbegin/graven/domain"
)

type VCSTool interface {
	VerifyRepoState(project *domain.Project, remote, branch string) error
	Tag(project *domain.Project, remote, tagname string) error
}
