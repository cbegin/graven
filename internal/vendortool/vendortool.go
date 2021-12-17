package vendortool

import (
	"github.com/cbegin/graven/internal/domain"
)

type VendorTool interface {
	Name() string
	VendorFileExists(project *domain.Project) bool
	LoadFile(project *domain.Project) error
	Dependencies() []PackageDepencency
}

type PackageDepencency interface {
	ArchiveFileName() string
	PackagePath() string
	Tag() string
}
