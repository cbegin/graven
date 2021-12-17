package builder

import (
	"github.com/cbegin/graven/internal/domain"
)

type BuildTool interface {
	Build(outputPath string, project *domain.Project, artifact *domain.Artifact, target *domain.Target) error
	Test(testPackage string, project *domain.Project) error
}
