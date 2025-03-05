//go:build !go1.16
// +build !go1.16

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// GoVersion reads the version of `go` that is on the PATH. This is done
// instead of `runtime.Version()` because it is possible to run gox against
// another Go version.
func GoVersion() (string, error) {
	// NOTE: We use `go run` instead of `go version` because the output
	// of `go version` might change whereas the source is guaranteed to run
	// for some time thanks to Go's compatibility guarantee.

	td, err := ioutil.TempDir("", "gox")
	if err != nil {
		return "", err
	}

	defer os.RemoveAll(td)

	// Write the source code for the program that will generate the version
	sourcePath := filepath.Join(td, "version.go")
	if err := ioutil.WriteFile(sourcePath, []byte(versionSource), 0644); err != nil {
		return "", err
	}

	// Execute and read the version, which will be the only thing on stdout.
	version, err := execGo("go", nil, "", "run", sourcePath)

	fmt.Printf("Detected Go Version: %s\n", version)

	return version, err
}
