package vendortool

import (
	"github.com/cbegin/graven/internal/domain"
	"strings"
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

func Coordinates(p PackageDepencency) string {
	path := p.PackagePath()
	parts := strings.Split(path, "/")
	domain := strings.Split(parts[0], ".")
	reverse(domain)
	return strings.Join([]string{strings.Join(domain, "/"), strings.Join(parts[1:], "/"), p.Tag(), p.ArchiveFileName()}, "/")
}

func reverse(ss []string) {
	last := len(ss) - 1
	for i := 0; i < len(ss)/2; i++ {
		ss[i], ss[last-i] = ss[last-i], ss[i]
	}
}
