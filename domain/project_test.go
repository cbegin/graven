package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/cbegin/graven/version"
	"strings"
)

func TestShouldLoadProject(t *testing.T) {
	p, err := LoadProject("../project.yaml")
	assert.NoError(t, err)

	assert.True(t, strings.HasSuffix(p.Name, "graven"))
	assert.Equal(t, p.Version, version.Version)
	assert.Equal(t, 3, len(p.Artifacts))
	assert.Equal(t, 0, len(p.Resources))
}
