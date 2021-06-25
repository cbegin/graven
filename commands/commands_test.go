package commands

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/cbegin/graven/domain"
	"github.com/cbegin/graven/test_fixtures/hello/version"
	"github.com/cbegin/graven/util"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func init() {
	domain.FindProject = func() (*domain.Project, error) {
		relativePath := "../test_fixtures/hello/project.yaml"
		absolutePath, _ := filepath.Abs(relativePath)
		return domain.LoadProject(absolutePath)
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

	assert.True(t, util.PathExists("../test_fixtures/hello/target/darwin/1.txt"))
	assert.True(t, util.PathExists("../test_fixtures/hello/target/darwin/3.txt"))
	assert.True(t, util.PathExists("../test_fixtures/hello/target/darwin/hello"))
	assert.True(t, util.PathExists("../test_fixtures/hello/target/darwin/Readme.md"))

	assert.True(t, util.PathExists("../test_fixtures/hello/target/linux/1.txt"))
	assert.True(t, util.PathExists("../test_fixtures/hello/target/linux/3.txt"))
	assert.True(t, util.PathExists("../test_fixtures/hello/target/linux/hello"))
	assert.True(t, util.PathExists("../test_fixtures/hello/target/linux/Readme.md"))

	assert.True(t, util.PathExists("../test_fixtures/hello/target/win/2.txt"))
	assert.True(t, util.PathExists("../test_fixtures/hello/target/win/3.txt"))
	assert.True(t, util.PathExists("../test_fixtures/hello/target/win/hello.exe"))
	assert.True(t, util.PathExists("../test_fixtures/hello/target/win/Readme.md"))
}

func TestShouldCleanTargetDirectory(t *testing.T) {
	c := &cli.Context{}

	err := build(c)
	assert.NoError(t, err)

	err = clean(c)
	assert.NoError(t, err)

	assert.False(t, util.PathExists("../test_fixtures/hello/target/darwin/hello"))
}

func TestShouldPackageTargetDirectory(t *testing.T) {
	c := &cli.Context{}

	err := clean(c)
	assert.NoError(t, err)

	err = pkg(c)
	assert.NoError(t, err)

	darwinPath := fmt.Sprintf("../test_fixtures/hello/target/hello-%s-darwin.tgz", version.Version)
	linuxPath := fmt.Sprintf("../test_fixtures/hello/target/hello-%s-linux.tar.gz", version.Version)
	winPath := fmt.Sprintf("../test_fixtures/hello/target/hello-%s-win.zip", version.Version)

	assert.True(t, util.PathExists(darwinPath))
	assert.True(t, util.PathExists(linuxPath))
	assert.True(t, util.PathExists(winPath))
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

	assert.True(t, util.PathExists(path.Join(tempdir, "version", "version.go")))
	assert.True(t, util.PathExists(path.Join(tempdir, "main.go")))
	assert.True(t, util.PathExists(path.Join(tempdir, "project.yaml")))
}

func TestShouldFreezeResources(t *testing.T) {
	c := &cli.Context{}

	_ = os.RemoveAll("../test_fixtures/hello/.modules")

	err := freeze(c)
	assert.NoError(t, err)

	assert.True(t, util.PathExists("../test_fixtures/hello/.modules/github-com-davecgh-go-spew-spew-346938d642f2ec3594ed81d874461961cd0faa76.zip"))
	assert.True(t, util.PathExists("../test_fixtures/hello/.modules/github-com-fatih-color-62e9147c64a1ed519147b62a56a14e83e2be02c1.zip"))
	assert.True(t, util.PathExists("../test_fixtures/hello/.modules/github-com-mattn-go-colorable-941b50ebc6efddf4c41c8e4537a5f68a4e686b24.zip"))
	assert.True(t, util.PathExists("../test_fixtures/hello/.modules/github-com-mattn-go-isatty-fc9e8d8ef48496124e79ae0df75490096eccf6fe.zip"))
	assert.True(t, util.PathExists("../test_fixtures/hello/.modules/github-com-pmezard-go-difflib-difflib-792786c7400a136282c1664665ae0a8db921c6c2.zip"))
	assert.True(t, util.PathExists("../test_fixtures/hello/.modules/github-com-stretchr-testify-assert-f6abca593680b2315d2075e0f5e2a9751e3f431a.zip"))

}

func TestShouldUnfreezeResources(t *testing.T) {
	c := &cli.Context{}

	_ = os.RemoveAll("../test_fixtures/hello/vendor/github.com")
	_ = os.RemoveAll("../test_fixtures/hello/vendor/golang.org")

	err := unfreeze(c)
	assert.NoError(t, err)

	assert.True(t, util.PathExists("../test_fixtures/hello/vendor/github.com/fatih"))
	assert.True(t, util.PathExists("../test_fixtures/hello/vendor/github.com/mattn"))
	assert.True(t, util.PathExists("../test_fixtures/hello/vendor/golang.org/x"))

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

	p, err := domain.LoadProject("../test_fixtures/hello/project.yaml")
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

	p, err := domain.LoadProject("../test_fixtures/hello/project.yaml")
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

	p, err := domain.LoadProject("../test_fixtures/hello/project.yaml")
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

	p, err := domain.LoadProject("../test_fixtures/hello/project.yaml")
	assert.NoError(t, err)
	assert.Equal(t, "0.0.1-DEV", p.Version)
}

func resetVersion() {
	util.CopyFile("../test_fixtures/hello/project.fixture", "../test_fixtures/hello/project.yaml")
	util.CopyFile("../test_fixtures/hello/version/version.fixture", "../test_fixtures/hello/version/version.go")
}
