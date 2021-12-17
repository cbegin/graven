package util

import (
	"io/fs"
	"time"
)

type CustomFileInfo struct {
	name    string
	size    int64
	mode    fs.FileMode
	modTime time.Time
	isDir   bool
	sys     interface{}
}

func CopyFileInfo(fi fs.FileInfo) *CustomFileInfo {
	return &CustomFileInfo{
		name:    fi.Name(),
		size:    fi.Size(),
		mode:    fi.Mode(),
		modTime: fi.ModTime(),
		isDir:   fi.IsDir(),
		sys:     fi.Sys(),
	}
}

func (c CustomFileInfo) Name() string {
	return c.name
}

func (c CustomFileInfo) Size() int64 {
	return c.size
}

func (c CustomFileInfo) Mode() fs.FileMode {
	return c.mode
}

func (c CustomFileInfo) ModTime() time.Time {
	return c.modTime
}

func (c CustomFileInfo) IsDir() bool {
	return c.isDir
}

func (c CustomFileInfo) Sys() interface{} {
	return c.sys
}
