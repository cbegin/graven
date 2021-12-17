package commands

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/cbegin/graven/internal/domain"
	"github.com/cbegin/graven/internal/util"
	"github.com/cbegin/graven/test/hello/version"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func init() {
	domain.FindProject = func() (*domain.Project, error) {
		relativePath := "../../test/hello/project.yaml"
		absolutePath, _ := filepath.Abs(relativePath)
		return domain.LoadProject(absolutePath)
	}
}

func TestShouldBuildTargetDirectory(t *testing.T) {
	c := &cli.Context{}

	err := clean(c)
	assert.NoError(t, err)

	err = build(c)
	assert.NoError(t, err)

	assert.True(t, util.PathExists("../../test/hello/target/darwin/1.txt"))
	assert.True(t, util.PathExists("../../test/hello/target/darwin/3.txt"))
	assert.True(t, util.PathExists("../../test/hello/target/darwin/hello"))
	assert.True(t, util.PathExists("../../test/hello/target/darwin/Readme.md"))

	assert.True(t, util.PathExists("../../test/hello/target/linux/1.txt"))
	assert.True(t, util.PathExists("../../test/hello/target/linux/3.txt"))
	assert.True(t, util.PathExists("../../test/hello/target/linux/hello"))
	assert.True(t, util.PathExists("../../test/hello/target/linux/Readme.md"))

	assert.True(t, util.PathExists("../../test/hello/target/win/2.txt"))
	assert.True(t, util.PathExists("../../test/hello/target/win/3.txt"))
	assert.True(t, util.PathExists("../../test/hello/target/win/hello.exe"))
	assert.True(t, util.PathExists("../../test/hello/target/win/Readme.md"))
}

func TestShouldCleanTargetDirectory(t *testing.T) {
	c := &cli.Context{}

	err := build(c)
	assert.NoError(t, err)

	err = clean(c)
	assert.NoError(t, err)

	assert.False(t, util.PathExists("../../test/hello/target/darwin/hello"))
}

func TestShouldPackageTargetDirectory(t *testing.T) {
	c := &cli.Context{}

	err := clean(c)
	assert.NoError(t, err)

	err = pkg(c)
	assert.NoError(t, err)

	darwinPath := fmt.Sprintf("../../test/hello/target/hello-%s-darwin.tgz", version.Version)
	linuxPath := fmt.Sprintf("../../test/hello/target/hello-%s-linux.tar.gz", version.Version)
	winPath := fmt.Sprintf("../../test/hello/target/hello-%s-win.zip", version.Version)

	assert.True(t, util.PathExists(darwinPath))
	assert.True(t, util.PathExists(linuxPath))
	assert.True(t, util.PathExists(winPath))
}

func TestShouldInitDirectory(t *testing.T) {
	tempdir := "../temp"
	_ = os.RemoveAll(tempdir)
	err := os.MkdirAll(tempdir, 0755)
	assert.NoError(t, err)
	defer func() {
		_ = os.RemoveAll(tempdir)
	}()

	c := &cli.Context{}

	wd, err := os.Getwd()
	assert.NoError(t, err)
	defer func() {
		_ = os.Chdir(wd) // making double sure we change the working directory back
	}()
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

	p, err := domain.LoadProject("../../test/hello/project.yaml")
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

	p, err := domain.LoadProject("../../test/hello/project.yaml")
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

	p, err := domain.LoadProject("../../test/hello/project.yaml")
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

	p, err := domain.LoadProject("../../test/hello/project.yaml")
	assert.NoError(t, err)
	assert.Equal(t, "0.0.1-DEV", p.Version)
}

func resetVersion() {
	_ = util.CopyFile("../../test/hello/project.fixture", "../../test/hello/project.yaml")
	_ = util.CopyFile("../../test/hello/version/version.fixture", "../../test/hello/version/version.go")
}
