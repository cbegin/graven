package vendortool

import (
	"testing"

	"github.com/cbegin/graven/domain"
	"github.com/stretchr/testify/assert"
)

func TestLoadVendorFile(t *testing.T) {
	vendorTool := &GovendorVendorTool{}
	p, err := domain.LoadProject("../hello/project.yaml")
	assert.NoError(t, err)
	err = vendorTool.LoadFile(p)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(vendorTool.Dependencies()))
	assert.Equal(t, "github-com-fatih-color-9131ab34cf20d2f6d83fdc67168a5430d1c7dc23.zip", vendorTool.Dependencies()[0].ArchiveFileName())
}
