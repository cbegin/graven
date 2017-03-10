package commands

import (
	"os"
	"testing"

	"github.com/cbegin/graven/domain"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
	"path"
	"fmt"
	"github.com/cbegin/graven/hello/version"
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

	darwinPath := fmt.Sprintf("../hello/target/hello-%s-darwin.tgz", version.Version)
	linuxPath := fmt.Sprintf("../hello/target/hello-%s-linux.tar.gz", version.Version)
	winPath := fmt.Sprintf("../hello/target/hello-%s-win.zip", version.Version)

	assert.True(t, PathExists(darwinPath));
	assert.True(t, PathExists(linuxPath));
	assert.True(t, PathExists(winPath));
}

func TestShouldInitDirectory(t *testing.T) {
	tempdir := "../temp"
	_ = os.RemoveAll(tempdir)
	err := os.MkdirAll(tempdir,0755)
	assert.NoError(t, err)
	defer os.RemoveAll(tempdir)

	c := &cli.Context{}

	wd, err := os.Getwd()
	assert.NoError(t, err)
	defer os.Chdir(wd) // making double sure we change the working directory back

	err = os.Chdir(tempdir)
	assert.NoError(t, err)

	err = initialize(c)
	assert.NoError(t, err)

	err = os.Chdir(wd)
	assert.NoError(t, err)

	assert.True(t, PathExists(path.Join(tempdir,"version", "version.go")));
	assert.True(t, PathExists(path.Join(tempdir, "main.go")));
	assert.True(t, PathExists(path.Join(tempdir, "project.yaml")));
}

func TestShouldFreezeResources(t *testing.T) {
	c := &cli.Context{}

	_ = os.RemoveAll("../hello/.freezer")

	err := freeze(c)
	assert.NoError(t, err)

	assert.True(t, PathExists("../hello/.freezer/github-com-fatih-color-9131ab34cf20d2f6d83fdc67168a5430d1c7dc23.zip"));
	assert.True(t, PathExists("../hello/.freezer/github-com-mattn-go-colorable-acb9493f2794fd0f820de7a27a217dafbb1b65ea.zip"));
	assert.True(t, PathExists("../hello/.freezer/github-com-mattn-go-isatty-57fdcb988a5c543893cc61bce354a6e24ab70022.zip"));
	assert.True(t, PathExists("../hello/.freezer/golang-org-x-sys-unix-99f16d856c9836c42d24e7ab64ea72916925fa97.zip"));

}

func PathExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
