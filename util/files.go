package util

import (
	"crypto/md5"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

// CopyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file. The file mode will be copied from the source and
// the copied data is synced/flushed to stable storage.
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	err = out.Sync()
	if err != nil {
		return err
	}

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return err
	}

	return err
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist. Symlinks are ignored and skipped.
func CopyDir(src string, dst string) error {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}

	mode := si.Mode()
	if !si.IsDir() {
		dirsi, err := os.Stat(path.Dir(src))
		if err != nil {
			return err
		}
		mode = dirsi.Mode()
	}

	err = os.MkdirAll(dst, mode)
	if err != nil {
		return err
	}

	if !si.IsDir() {
		return CopyFile(src, path.Join(dst, path.Base(src)))
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}

	return err
}

func CompareFileContents(a, b string) (bool, error) {
	adata, err := ioutil.ReadFile(a)
	if err != nil {
		return false, fmt.Errorf("File compare error (a): %v", err)

	}
	bdata, err := ioutil.ReadFile(b)
	if err != nil {
		return false, fmt.Errorf("File compare error (b): %v", err)
	}
	return string(adata) == string(bdata), nil
}

func MD5File(filePath string, h hash.Hash) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := io.Copy(h, file); err != nil {
		return err
	}
	return nil
}

func GetMD5Walker(basePath string, h hash.Hash) func(fp string, fi os.FileInfo, err error) error {
	return func(fp string, fi os.FileInfo, err error) error {
		if err != nil {
			return nil // skip on error
		}
		if fi.IsDir() {
			return nil // ignore directories
		}
		shortPath := fp[len(basePath):]
		_, err = h.Write([]byte(shortPath))
		if err != nil {
			return err
		}
		MD5File(fp, h)
		return nil
	}
}

func MD5Dir(basePath string) ([]byte, error) {
	var result []byte
	h := md5.New()
	err := filepath.Walk(basePath, GetMD5Walker(basePath, h))
	if err != nil {
		return result, err
	}
	result = h.Sum(result)
	return result, nil
}

func CompareDir(a, b string) (bool, error) {
	ahash, err := MD5Dir(a)
	if err != nil {
		return false, err
	}
	bhash, err := MD5Dir(b)
	if err != nil {
		return false, err
	}
	return string(ahash) == string(bhash), nil
}
