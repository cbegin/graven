package util


import (
	"testing"
	"fmt"
	"os"
	"io/ioutil"
	"github.com/stretchr/testify/assert"
)

func TestZipUnzip(t *testing.T) {
	os.Mkdir("../temp", 0755)
	if err := ZipDir("../hello", "../temp/hello.zip"); err != nil {
		fmt.Printf("Zip error: %v", err)
		t.FailNow()
	}
	if err := UnzipDir("../temp/hello.zip", "../temp/hello"); err != nil {
		fmt.Printf("Unzip error: %v", err)
		t.FailNow()
	}

	orig, err := ioutil.ReadFile("../hello/hello.go")
	if err != nil {
		fmt.Printf("Zip error: %v", err)
		t.FailNow()
	}
	unzipped, err := ioutil.ReadFile("../temp/hello/hello.go")
	if err != nil {
		fmt.Printf("Unzip error: %v", err)
		t.FailNow()
	}

	assert.Equal(t, orig, unzipped)

	os.RemoveAll("../temp")
}

func TestTarUntar(t *testing.T) {
	os.Mkdir("../temp", 0755)
	if err := TarDir("../hello", "../temp/hello.tar.gz"); err != nil {
		fmt.Printf("Tar error: %v", err)
		t.FailNow()
	}
	if err := UntarDir("../temp/hello.tar.gz", "../temp/hello"); err != nil {
		fmt.Printf("Untar error: %v", err)
		t.FailNow()
	}

	orig, err := ioutil.ReadFile("../hello/hello.go")
	if err != nil {
		fmt.Printf("Tar file error: %v", err)
		t.FailNow()
	}
	unzipped, err := ioutil.ReadFile("../temp/hello/hello.go")
	if err != nil {
		fmt.Printf("Untar file error: %v", err)
		t.FailNow()
	}

	assert.Equal(t, orig, unzipped)

	os.RemoveAll("../temp")
}