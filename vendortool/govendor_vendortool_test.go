package vendortool

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/cbegin/graven/domain"
)

func TestLoadVendorFile(t *testing.T) {
	f := &GovendorFile{}
	p, err := domain.LoadProject("../hello/project.yaml")
	assert.NoError(t, err)
	err = f.LoadFile(p)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(f.Dependencies()))
	assert.Equal(t, "github-com-fatih-color-9131ab34cf20d2f6d83fdc67168a5430d1c7dc23.zip", f.Dependencies()[0].ArchiveFileName())
}
