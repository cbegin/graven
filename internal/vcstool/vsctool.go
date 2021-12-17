package vcstool

import (
	"github.com/cbegin/graven/internal/domain"
)

type VCSTool interface {
	VerifyRepoState(project *domain.Project, remote, branch string) error
	Tag(project *domain.Project, remote, tagname string) error
}
