package vendortool

import "github.com/cbegin/graven/domain"

type VendorTool interface {
	LoadFile(project *domain.Project) error
	Dependencies() []PackageDepencency
}

type PackageDepencency interface {
	ArchiveFileName() string
	PackagePath() string
}