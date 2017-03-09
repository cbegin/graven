package commands

import (
	"github.com/cbegin/graven/domain"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
	"os"
)

func init() {
	domain.FindProject = func () (*domain.Project, error) {
		return domain.LoadProject("../hello/project.yaml")
	}
}

func TestShouldBuildTargetDirectory(t *testing.T) {
	c := &cli.Context{}

	err := clean(c)
	assert.NoError(t, err)

	err = build(c)
	assert.NoError(t, err)

	assert.True(t, PathExists("../hello/target/darwin/1.txt"))
	assert.True(t, PathExists("../hello/target/darwin/3.txt"))
	assert.True(t, PathExists("../hello/target/darwin/hello"))
	assert.True(t, PathExists("../hello/target/darwin/Readme.md"))

	assert.True(t, PathExists("../hello/target/linux/1.txt"))
	assert.True(t, PathExists("../hello/target/linux/3.txt"))
	assert.True(t, PathExists("../hello/target/linux/hello"))
	assert.True(t, PathExists("../hello/target/linux/Readme.md"))

	assert.True(t, PathExists("../hello/target/win/2.txt"))
	assert.True(t, PathExists("../hello/target/win/3.txt"))
	assert.True(t, PathExists("../hello/target/win/hello.exe"))
	assert.True(t, PathExists("../hello/target/win/Readme.md"))
}

func TestShouldCleanTargetDirectory(t *testing.T) {
	c := &cli.Context{}

	err := build(c)
	assert.NoError(t, err)

	err = clean(c)
	assert.NoError(t, err)

	assert.False(t, PathExists("../hello/target/darwin/hello"))
}


func TestShouldPackageTargetDirectory(t *testing.T) {
	c := &cli.Context{}

	err := clean(c)
	assert.NoError(t, err)

	err = pkg(c)
	assert.NoError(t, err)

	//├── hello-0.0.1-darwin.tgz
	//├── hello-0.0.1-linux.tar.gz
	//├── hello-0.0.1-win.zip

	assert.True(t, PathExists("../hello/target/darwin"));
	assert.True(t, PathExists("../hello/target/linux"));
	assert.True(t, PathExists("../hello/target/win"));
}

func PathExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
