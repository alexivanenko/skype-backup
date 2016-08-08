package core

import (
	"fmt"
	"os"
	"path/filepath"
)

var rootDir string

func GetRootDir() string {
	if rootDir == "" {
		/*_, filename, _, ok := runtime.Caller(0)
		if !ok {
			panic("No caller information")
		}

		rootDir = path.Dir(filename)
		rootDir = strings.Replace(rootDir, "/config", "", 1)*/

		var err error
		if rootDir, err = filepath.Abs(fmt.Sprintf("%s/", filepath.Dir(os.Args[0]))); err != nil {
			panic(err)
		}
	}

	return rootDir
}

func OpenFile(path string) (*os.File, os.FileInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to open file: %s", err)
	}

	info, err := f.Stat()
	if err != nil {
		return nil, nil, fmt.Errorf("Failed getting file metadata: %s", err)
	}

	return f, info, nil
}

func CreateFile(path string) (*os.File, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to open file: %s", err)
	}

	return f, nil
}
