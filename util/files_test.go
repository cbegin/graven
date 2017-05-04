package util

import (
	"testing"
	"os"

	"github.com/stretchr/testify/assert"
)

func TestCopyDir (t *testing.T) {
	CopyDir("../hello", "../temp/hello")
	if same, err := CompareDir("../hello", "../temp/hello"); err != nil {
		assert.FailNow(t, "Directory comparison failed: %v", err)
	} else {
		assert.True(t, same)
	}
	os.RemoveAll("../temp")
}

func TestCompareDirTrue(t *testing.T) {
	if same, err := CompareDir("../hello", "../hello"); err != nil {
		assert.FailNow(t, "Directory comparison failed: %v", err)
	} else {
		assert.True(t, same)
	}
}

func TestCompareDirFalse(t *testing.T) {
	if same, err := CompareDir("../resources", "../hello"); err != nil {
		assert.FailNow(t, "Directory comparison failed", "%v", err)
	} else {
		assert.False(t, same)
	}
}

func TestCompareDirFalseReverse(t *testing.T) {
	if same, err := CompareDir("../hello", "../resources"); err != nil {
		assert.FailNow(t, "Directory comparison failed", "%v", err)
	} else {
		assert.False(t, same)
	}
}

func TestCompareMissingDir(t *testing.T) {
	if same, err := CompareDir("../fakedir", "../hello"); err != nil {
		assert.Error(t, err)
	} else {
		assert.False(t, same)
	}
}
