package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/cbegin/graven/hello/version"
)

func TestShouldFindProject(t *testing.T) {
	p, err := FindProject()
	assert.NoError(t, err)
	assert.Equal(t, p.Name, "graven")
}

func TestShouldLoadProject(t *testing.T) {
	p, err := LoadProject("../hello/project.yaml")
	assert.NoError(t, err)

	assert.Equal(t, p.Name, "hello")
	assert.Equal(t, p.Version, version.Version)
	assert.Equal(t, 3, len(p.Artifacts))
	assert.Equal(t, 2, len(p.Resources))
	artifactMap := map[string]Artifact{}
	for _, a := range p.Artifacts {
		artifactMap[a.Classifier] = a
	}
	linux, ok := artifactMap["linux"]
	assert.True(t, ok)
	assert.Equal(t, 1, len(linux.Targets))
	assert.Equal(t, 1, len(linux.Resources))
	assert.Equal(t, 2, len(linux.Environment))
	assert.Equal(t, "tar.gz", linux.Archive)

	for _, g := range linux.Targets {
		assert.Equal(t, "hello", g.Executable)
		assert.Equal(t, ".", g.Package)
		assert.Equal(t, "-p=4", g.Flags)
	}
}
