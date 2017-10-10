package vcstool

import (
	"github.com/cbegin/graven/domain"
)

type VCSTool interface {
	VerifyRepoState(project *domain.Project, branch string) error
	Tag(project *domain.Project, tagname string) error
}
