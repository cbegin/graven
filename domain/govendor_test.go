package domain

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestLoadVendorFile(t *testing.T) {
	p, err := LoadProject("../hello/project.yaml")
	assert.NoError(t, err)
	g, err := ReadGovendorFile(p)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(g.Packages))
}
