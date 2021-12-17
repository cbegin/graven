package util

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZipUnzip(t *testing.T) {
	os.Mkdir("../temp", 0755)
	if err := ZipDir("../../test/hello", "../temp/hello.zip", true); err != nil {
		fmt.Printf("Zip error: %v", err)
		t.FailNow()
	}

	if err := UnzipDir("../temp/hello.zip", "../temp/hello"); err != nil {
		fmt.Printf("Unzip error: %v", err)
		t.FailNow()
	}

	if same, err := CompareFileContents("../../test/hello/hello.go", "../temp/hello/hello.go"); err != nil {
		fmt.Printf("Error comparing files: %v", err)
		t.FailNow()
	} else {
		assert.True(t, same)
	}

	os.RemoveAll("../temp")
}

func TestTarUntar(t *testing.T) {
	os.Mkdir("../temp", 0755)
	if err := TarDir("../../test/hello", "../temp/hello.tar.gz"); err != nil {
		fmt.Printf("Tar error: %v", err)
		t.FailNow()
	}

	if err := UntarDir("../temp/hello.tar.gz", "../temp/hello"); err != nil {
		fmt.Printf("Untar error: %v", err)
		t.FailNow()
	}

	if same, err := CompareFileContents("../../test/hello/hello.go", "../temp/hello/hello.go"); err != nil {
		fmt.Printf("Error comparing files: %v", err)
		t.FailNow()
	} else {
		assert.True(t, same)
	}

	os.RemoveAll("../temp")
}
