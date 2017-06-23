package vendortool

import (
	"testing"

	"github.com/cbegin/graven/domain"
	"github.com/stretchr/testify/assert"
)

func TestLoadGlideFile(t *testing.T) {
	vendorTool := &GlideVendorTool{}
	p, err := domain.LoadProject("../hello/project.yaml")
	assert.NoError(t, err)
	err = vendorTool.LoadFile(p)
	assert.NoError(t, err)
	assert.Equal(t, 7, len(vendorTool.Dependencies()))
	assert.Equal(t, "github-com-fatih-color-570b54cabe6b8eb0bc2dfce68d964677d63b5260.zip", vendorTool.Dependencies()[0].ArchiveFileName())
}
