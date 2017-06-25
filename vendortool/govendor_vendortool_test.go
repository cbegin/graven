package vendortool

import (
	"testing"

	"github.com/cbegin/graven/domain"
	"github.com/stretchr/testify/assert"
)

func TestLoadGovendorFile(t *testing.T) {
	vendorTool := &GovendorVendorTool{}
	p, err := domain.LoadProject("../test_fixtures/hello/project.yaml")
	assert.NoError(t, err)
	err = vendorTool.LoadFile(p)
	assert.NoError(t, err)
	assert.Equal(t, 7, len(vendorTool.Dependencies()))
	assert.Equal(t, "github-com-davecgh-go-spew-spew-346938d642f2ec3594ed81d874461961cd0faa76.zip", vendorTool.Dependencies()[0].ArchiveFileName())
}
