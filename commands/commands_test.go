package commands

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/cbegin/graven/domain"
	"github.com/cbegin/graven/hello/version"
	"github.com/cbegin/graven/util"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func init() {
	domain.FindProject = func() (*domain.Project, error) {
		return domain.LoadProject("../hello/project.yaml")
	}
	c := &cli.Context{}
	_ = unfreeze(c)
}

func TestShouldBuildTargetDirectory(t *testing.T) {
	c := &cli.Context{}

	err := clean(c)
	assert.NoError(t, err)

	err = build(c)
	assert.NoError(t, err)

	assert.True(t, pathExists("../hello/target/darwin/1.txt"))
	assert.True(t, pathExists("../hello/target/darwin/3.txt"))
	assert.True(t, pathExists("../hello/target/darwin/hello"))
	assert.True(t, pathExists("../hello/target/darwin/Readme.md"))

	assert.True(t, pathExists("../hello/target/linux/1.txt"))
	assert.True(t, pathExists("../hello/target/linux/3.txt"))
	assert.True(t, pathExists("../hello/target/linux/hello"))
	assert.True(t, pathExists("../hello/target/linux/Readme.md"))

	assert.True(t, pathExists("../hello/target/win/2.txt"))
	assert.True(t, pathExists("../hello/target/win/3.txt"))
	assert.True(t, pathExists("../hello/target/win/hello.exe"))
	assert.True(t, pathExists("../hello/target/win/Readme.md"))
}

func TestShouldCleanTargetDirectory(t *testing.T) {
	c := &cli.Context{}

	err := build(c)
	assert.NoError(t, err)

	err = clean(c)
	assert.NoError(t, err)

	assert.False(t, pathExists("../hello/target/darwin/hello"))
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

	assert.True(t, pathExists(darwinPath))
	assert.True(t, pathExists(linuxPath))
	assert.True(t, pathExists(winPath))
}

func TestShouldInitDirectory(t *testing.T) {
	tempdir := "../temp"
	_ = os.RemoveAll(tempdir)
	err := os.MkdirAll(tempdir, 0755)
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

	assert.True(t, pathExists(path.Join(tempdir, "version", "version.go")))
	assert.True(t, pathExists(path.Join(tempdir, "main.go")))
	assert.True(t, pathExists(path.Join(tempdir, "project.yaml")))
}

func TestShouldFreezeResources(t *testing.T) {
	c := &cli.Context{}

	_ = os.RemoveAll("../hello/.freezer")

	err := freeze(c)
	assert.NoError(t, err)

	assert.True(t, pathExists("../hello/.freezer/github-com-fatih-color-9131ab34cf20d2f6d83fdc67168a5430d1c7dc23.zip"))
	assert.True(t, pathExists("../hello/.freezer/github-com-mattn-go-colorable-acb9493f2794fd0f820de7a27a217dafbb1b65ea.zip"))
	assert.True(t, pathExists("../hello/.freezer/github-com-mattn-go-isatty-57fdcb988a5c543893cc61bce354a6e24ab70022.zip"))
	assert.True(t, pathExists("../hello/.freezer/golang-org-x-sys-unix-99f16d856c9836c42d24e7ab64ea72916925fa97.zip"))

}

func TestShouldUnfreezeResources(t *testing.T) {
	c := &cli.Context{}

	_ = os.RemoveAll("../hello/vendor/github.com")
	_ = os.RemoveAll("../hello/vendor/golang.org")

	err := unfreeze(c)
	assert.NoError(t, err)

	assert.True(t, pathExists("../hello/vendor/github.com/fatih"))
	assert.True(t, pathExists("../hello/vendor/github.com/mattn"))
	assert.True(t, pathExists("../hello/vendor/golang.org/x"))

}

func TestShouldRunTests(t *testing.T) {
	c := &cli.Context{}

	err := tester(c)
	assert.NoError(t, err)
}

func TestShouldBumpPatchVersion(t *testing.T) {
	// clean up before and after
	resetVersion()
	defer resetVersion()

	app := cli.NewApp()
	app.Commands = []cli.Command{
		BumpCommand,
	}
	err := app.Run([]string{"graven", "bump", "patch"})
	assert.NoError(t, err)

	p, err := domain.LoadProject("../hello/project.yaml")
	assert.NoError(t, err)
	assert.Equal(t, "0.0.2", p.Version)
}

func TestShouldBumpMinorVersion(t *testing.T) {
	// clean up before and after
	resetVersion()
	defer resetVersion()

	app := cli.NewApp()
	app.Commands = []cli.Command{
		BumpCommand,
	}
	err := app.Run([]string{"graven", "bump", "minor"})
	assert.NoError(t, err)

	p, err := domain.LoadProject("../hello/project.yaml")
	assert.NoError(t, err)
	assert.Equal(t, "0.1.0", p.Version)
}

func TestShouldBumpMajorVersion(t *testing.T) {
	// clean up before and after
	resetVersion()
	defer resetVersion()

	app := cli.NewApp()
	app.Commands = []cli.Command{
		BumpCommand,
	}
	err := app.Run([]string{"graven", "bump", "major"})
	assert.NoError(t, err)

	p, err := domain.LoadProject("../hello/project.yaml")
	assert.NoError(t, err)
	assert.Equal(t, "1.0.0", p.Version)
}

func TestShouldSetVersionQualifier(t *testing.T) {
	// clean up before and after
	resetVersion()
	defer resetVersion()

	app := cli.NewApp()
	app.Commands = []cli.Command{
		BumpCommand,
	}
	err := app.Run([]string{"graven", "bump", "DEV"})
	assert.NoError(t, err)

	p, err := domain.LoadProject("../hello/project.yaml")
	assert.NoError(t, err)
	assert.Equal(t, "0.0.1-DEV", p.Version)
}

func resetVersion() {
	util.CopyFile("../hello/project.fixture", "../hello/project.yaml")
	util.CopyFile("../hello/version/version.fixture", "../hello/version/version.go")
}

func pathExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
