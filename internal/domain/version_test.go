package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldParseMajorMinorPatchQualifier(t *testing.T) {
	v := Version{}
	v.Parse("1.2.3-DEV")
	assert.Equal(t, 1, v.Major)
	assert.Equal(t, 2, v.Minor)
	assert.Equal(t, 3, v.Patch)
	assert.Equal(t, "DEV", v.Qualifier)
	assert.Equal(t, "1.2.3-DEV", v.ToString())
}

func TestShouldParseMajorMinorPatch(t *testing.T) {
	v := Version{}
	v.Parse("1.2.3")
	assert.Equal(t, 1, v.Major)
	assert.Equal(t, 2, v.Minor)
	assert.Equal(t, 3, v.Patch)
	assert.Equal(t, "", v.Qualifier)
	assert.Equal(t, "1.2.3", v.ToString())
}

func TestShouldParseMajorMinor(t *testing.T) {
	v := Version{}
	v.Parse("1.2")
	assert.Equal(t, 1, v.Major)
	assert.Equal(t, 2, v.Minor)
	assert.Equal(t, 0, v.Patch)
	assert.Equal(t, "", v.Qualifier, "")
	assert.Equal(t, "1.2.0", v.ToString())
}

func TestShouldParseMajor(t *testing.T) {
	v := Version{}
	v.Parse("1")
	assert.Equal(t, 1, v.Major)
	assert.Equal(t, 0, v.Minor)
	assert.Equal(t, 0, v.Patch)
	assert.Equal(t, "", v.Qualifier)
	assert.Equal(t, "1.0.0", v.ToString())
}

func TestShouldParseBadQualifierFormat(t *testing.T) {
	v := Version{}
	err := v.Parse("1-1-1")
	assert.NotNil(t, err)
}

func TestShouldParseBadVersionFormat(t *testing.T) {
	v := Version{}
	err := v.Parse("1.1.1.1")
	assert.NotNil(t, err)
}
