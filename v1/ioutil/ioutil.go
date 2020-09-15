package ioutil

import (
	"os"
)

type LazyFile struct {
	path string
	file os.File
}

func OpenLazy(p string) *LazyFile {
	return &LazyFile{p, nil}
}

func (f *LazyFile) Read(p []byte) (int, error) {
	var err error
	if f.file == nil {
		f.file, err = os.Open(f.path)
		if err != nil {
			return 0, err
		}
	}
	return f.file.Read(p)
}

func (f *LazyFile) Close() error {
	if f.file != nil {
		return f.file.Close()
	} else {
		return nil
	}
}
