package util

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopyDir(t *testing.T) {
	CopyDir("../../test/hello", "../temp/hello")
	if same, err := CompareDir("../../test/hello", "../temp/hello"); err != nil {
		assert.FailNow(t, "Directory comparison failed: %v", err)
	} else {
		assert.True(t, same)
	}
	os.RemoveAll("../temp")
}

func TestCompareDirTrue(t *testing.T) {
	if same, err := CompareDir("../../test/hello", "../../test/hello"); err != nil {
		assert.FailNow(t, "Directory comparison failed: %v", err)
	} else {
		assert.True(t, same)
	}
}

func TestCompareDirFalse(t *testing.T) {
	if same, err := CompareDir("../resources", "../../test/hello"); err != nil {
		assert.FailNow(t, "Directory comparison failed", "%v", err)
	} else {
		assert.False(t, same)
	}
}

func TestCompareDirFalseReverse(t *testing.T) {
	if same, err := CompareDir("../../test/hello", "../resources"); err != nil {
		assert.FailNow(t, "Directory comparison failed", "%v", err)
	} else {
		assert.False(t, same)
	}
}

func TestCompareMissingDir(t *testing.T) {
	if same, err := CompareDir("../fakedir", "../../test/hello"); err != nil {
		assert.Error(t, err)
	} else {
		assert.False(t, same)
	}
}
