package domain

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestLoadVendorFile(t *testing.T) {
	p, err := LoadProject("../hello/project.yaml")
	assert.NoError(t, err)
	g, err := LoadGovendorFile(p)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(g.Packages))
	assert.Equal(t, "github-com-fatih-color-9131ab34cf20d2f6d83fdc67168a5430d1c7dc23.zip", g.Packages[0].ArchiveFileName())
}
