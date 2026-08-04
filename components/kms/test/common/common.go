package common

import (
	"log"
	"path"
	"runtime"
)

// GetTestRootDir gets the test root directory.
func GetTestRootDir() string {
	_, filename, _, _ := runtime.Caller(0)
	root := path.Clean(path.Join(path.Dir(filename), ".."))

	if path.Base(root) != "test" {
		log.Panic("Could not determine root directory of tests.")
	}

	return root
}
