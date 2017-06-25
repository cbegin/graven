package vendortool

import (
	"testing"

	"github.com/cbegin/graven/domain"
	"github.com/stretchr/testify/assert"
)

func TestLoadDepFile(t *testing.T) {
	vendorTool := &DepVendorTool{}
	p, err := domain.LoadProject("../test_fixtures/hello/project.yaml")
	assert.NoError(t, err)
	err = vendorTool.LoadFile(p)
	assert.NoError(t, err)
	assert.Equal(t, 7, len(vendorTool.Dependencies()))
	assert.Equal(t, "github-com-davecgh-go-spew-6d212800a42e8ab5c146b8ace3490ee17e5225f9.zip", vendorTool.Dependencies()[0].ArchiveFileName())
}
